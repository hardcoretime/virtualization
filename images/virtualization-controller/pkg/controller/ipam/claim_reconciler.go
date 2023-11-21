package ipam

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	virtv2 "github.com/deckhouse/virtualization-controller/api/v2alpha1"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/two_phase_reconciler"
)

type ClaimReconciler struct {
	ParsedCIDRs []*net.IPNet
}

func NewClaimReconciler(vmCIDRs []string) (*ClaimReconciler, error) {
	parsedCIDRs := make([]*net.IPNet, len(vmCIDRs))

	for i, cidr := range vmCIDRs {
		_, parsedCIDR, err := net.ParseCIDR(cidr)
		if err != nil || parsedCIDR == nil {
			return nil, fmt.Errorf("failed to parse virtual cide %s: %w", cidr, err)
		}

		parsedCIDRs[i] = parsedCIDR
	}

	return &ClaimReconciler{
		ParsedCIDRs: parsedCIDRs,
	}, nil
}

func (r *ClaimReconciler) SetupController(_ context.Context, mgr manager.Manager, ctr controller.Controller) error {
	if err := ctr.Watch(
		source.Kind(mgr.GetCache(), &virtv2.VirtualMachineIPAddressLease{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueRequestsFromLeases),
	); err != nil {
		return fmt.Errorf("error setting watch on leases: %w", err)
	}

	return ctr.Watch(source.Kind(mgr.GetCache(), &virtv2.VirtualMachineIPAddressClaim{}), &handler.EnqueueRequestForObject{},
		predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool { return true },
			DeleteFunc: func(e event.DeleteEvent) bool { return true },
			UpdateFunc: func(e event.UpdateEvent) bool { return true },
		},
	)
}

func (r *ClaimReconciler) enqueueRequestsFromLeases(_ context.Context, obj client.Object) []reconcile.Request {
	lease, ok := obj.(*virtv2.VirtualMachineIPAddressLease)
	if !ok {
		return nil
	}

	if lease.Spec.ClaimRef == nil {
		return nil
	}

	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Namespace: lease.Spec.ClaimRef.Namespace,
				Name:      lease.Spec.ClaimRef.Name,
			},
		},
	}
}

func (r *ClaimReconciler) Sync(ctx context.Context, _ reconcile.Request, state *ClaimReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	switch {
	case state.Lease == nil && state.Claim.Current().Spec.LeaseName != "":
		opts.Log.Info("Lease by name not found: waiting for the lease to be available")
		return nil

	case state.Lease == nil:
		// Lease not found by spec.LeaseName or spec.Address: it doesn't exist.
		opts.Log.Info("No Lease for Claim: create the new one", "address", state.Claim.Current().Spec.Address, "leaseName", state.Claim.Current().Spec.LeaseName)

		leaseName := state.Claim.Current().Spec.LeaseName

		if state.Claim.Current().Spec.Address == "" {
			if leaseName != "" {
				opts.Log.Info("Claim address omitted in the spec: extract from the lease name")
				state.Claim.Changed().Spec.Address = leaseNameToIP(leaseName)
			} else {
				opts.Log.Info("Claim address omitted in the spec: allocate the new one")
				var err error
				state.Claim.Changed().Spec.Address, err = r.allocateNewIP(state.AllocatedIPs)
				if err != nil {
					return err
				}
			}
		}

		if !r.isAvailableAddress(state.Claim.Changed().Spec.Address, state.AllocatedIPs) {
			opts.Log.Info("Claim cannot be created: the address has already been allocated for another claim", "address", state.Claim.Current().Spec.Address)
			return nil
		}

		if leaseName == "" {
			leaseName = ipToLeaseName(state.Claim.Changed().Spec.Address)
		}

		err := opts.Client.Create(ctx, &virtv2.VirtualMachineIPAddressLease{
			ObjectMeta: metav1.ObjectMeta{
				Name: leaseName,
			},
			Spec: virtv2.VirtualMachineIPAddressLeaseSpec{
				ClaimRef: &virtv2.VirtualMachineIPAddressLeaseClaimRef{
					Name:      state.Claim.Name().Name,
					Namespace: state.Claim.Name().Namespace,
				},
			},
		})
		if err != nil {
			return err
		}

		state.Claim.Changed().Spec.LeaseName = leaseName
		return opts.Client.Update(ctx, state.Claim.Changed())

	case isBoundedLease(state):
		opts.Log.Info("Lease already exists, claim ref is valid")
		return nil

	case state.Lease.Status.Phase == "":
		opts.Log.Info("Lease is not ready: waiting for the lease")
		state.SetReconcilerResult(&reconcile.Result{Requeue: true, RequeueAfter: 2 * time.Second})
		return nil

	case state.Lease.Status.Phase == virtv2.VirtualMachineIPAddressLeasePhaseBound:
		opts.Log.Info("Lease is bounded to another claim: recreate claim when the lease is released")
		return nil

	default:
		opts.Log.Info("Lease is released: set binding")

		state.Lease.Spec.ClaimRef = &virtv2.VirtualMachineIPAddressLeaseClaimRef{
			Name:      state.Claim.Name().Name,
			Namespace: state.Claim.Name().Namespace,
		}

		err := opts.Client.Update(ctx, state.Lease)
		if err != nil {
			return err
		}

		state.Claim.Changed().Spec.LeaseName = state.Lease.Name
		state.Claim.Changed().Spec.Address = leaseNameToIP(state.Lease.Name)
		return opts.Client.Update(ctx, state.Claim.Changed())
	}
}

func (r *ClaimReconciler) UpdateStatus(_ context.Context, _ reconcile.Request, state *ClaimReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	// Do nothing if object is being deleted as any update will lead to en error.
	if state.isDeletion() {
		return nil
	}

	claimStatus := state.Claim.Current().Status.DeepCopy()

	claimStatus.Address = ""
	claimStatus.ConflictMessage = ""

	switch {
	case state.Lease == nil && state.Claim.Current().Spec.LeaseName != "":
		claimStatus.Phase = virtv2.VirtualMachineIPAddressClaimPhaseLost

	case state.Lease == nil:
		claimStatus.Phase = virtv2.VirtualMachineIPAddressClaimPhasePending

	case isBoundedLease(state):
		claimStatus.Phase = virtv2.VirtualMachineIPAddressClaimPhaseBound
		claimStatus.Address = state.Claim.Current().Spec.Address

	case state.Lease.Status.Phase == virtv2.VirtualMachineIPAddressLeasePhaseBound:
		claimStatus.Phase = virtv2.VirtualMachineIPAddressClaimPhaseConflict

		// There is only one way to automatically link Claim in phase Conflict with recently released Lease: only with cyclic reconciliation (with an interval of N seconds).
		// At the moment this looks redundant, so Claim in the phase Conflict will not be able to bind the recently released Lease.
		// It is necessary to recreate Claim manually in order to link it to released Lease.
		claimStatus.ConflictMessage = "Lease is bounded to another claim: please recreate claim when the lease is released"

	default:
		claimStatus.Phase = virtv2.VirtualMachineIPAddressClaimPhasePending
	}

	opts.Log.Info("Set claim phase Pending", "phase", claimStatus.Phase)

	state.Claim.Changed().Status = *claimStatus

	return nil
}

func isBoundedLease(state *ClaimReconcilerState) bool {
	if state.Lease.Spec.ClaimRef == nil {
		return false
	}

	if state.Lease.Spec.ClaimRef.Namespace != state.Claim.Name().Namespace || state.Lease.Spec.ClaimRef.Name != state.Claim.Name().Name {
		return false
	}

	return true
}

func (r *ClaimReconciler) isAvailableAddress(address string, allocatedIPs AllocatedIPs) bool {
	ip := net.ParseIP(address)

	if _, ok := allocatedIPs[ip.String()]; !ok {
		for _, cidr := range r.ParsedCIDRs {
			if cidr.Contains(ip) {
				// available
				return true
			}
		}
		// out of range
		return false
	}
	// already exists
	return false
}

func (r *ClaimReconciler) allocateNewIP(allocatedIPs AllocatedIPs) (string, error) {
	for _, cidr := range r.ParsedCIDRs {
		for ip := cidr.IP.Mask(cidr.Mask); cidr.Contains(ip); inc(ip) {
			_, ok := allocatedIPs[ip.String()]
			if !ok {
				return ip.String(), nil
			}
		}
	}
	return "", errors.New("no remaining ips")
}

const ipPrefix = "ip-"

func ipToLeaseName(ip string) string {
	addr := net.ParseIP(ip)
	if addr.To4() != nil {
		// IPv4 address
		return ipPrefix + strings.ReplaceAll(addr.String(), ".", "-")
	}

	return ""
}

func leaseNameToIP(leaseName string) string {
	if strings.HasPrefix(leaseName, ipPrefix) && len(leaseName) > len(ipPrefix) {
		return strings.ReplaceAll(leaseName[len(ipPrefix):], "-", ".")
	}

	return ""
}

// http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
