package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	virtv1 "kubevirt.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	virtv2alpha1 "github.com/deckhouse/virtualization-controller/api/v2alpha1"
	"github.com/deckhouse/virtualization-controller/pkg/common"
	appconfig "github.com/deckhouse/virtualization-controller/pkg/config"
	"github.com/deckhouse/virtualization-controller/pkg/controller"
	"github.com/deckhouse/virtualization-controller/pkg/controller/ipam"
)

var (
	log                  = logf.Log.WithName("cmd")
	resourcesSchemeFuncs = []func(*apiruntime.Scheme) error{
		clientgoscheme.AddToScheme,
		extv1.AddToScheme,
		virtv2alpha1.AddToScheme,
		cdiv1.AddToScheme,
		virtv1.AddToScheme,
	}
	importerImage       string
	uploaderImage       string
	controllerNamespace string
)

const defaultVerbosity = "1"

func init() {
	importerImage = getRequiredEnvVar(common.ImporterPodImageNameVar)
	uploaderImage = getRequiredEnvVar(common.UploaderPodImageNameVar)
	controllerNamespace = getRequiredEnvVar(common.PodNamespaceVar)
}

func setupLogger() {
	verbose := defaultVerbosity
	if verboseEnvVarVal := os.Getenv("VERBOSITY"); verboseEnvVarVal != "" {
		verbose = verboseEnvVarVal
	}
	// visit actual flags passed in and if passed check -v and set verbose
	if fv := flag.Lookup("v"); fv != nil {
		verbose = fv.Value.String()
	}
	if verbose == defaultVerbosity {
		log.V(1).Info(fmt.Sprintf("Note: increase the -v level in the controller deployment for more detailed logging, eg. -v=%d or -v=%d\n", 2, 3))
	}
	verbosityLevel, err := strconv.Atoi(verbose)
	debug := false
	if err == nil && verbosityLevel > 1 {
		debug = true
	}

	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	logf.SetLogger(zap.New(zap.Level(zapcore.Level(-1*verbosityLevel)), zap.UseDevMode(debug)))
}

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

func getRequiredEnvVar(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Error(fmt.Errorf("environment variable %q undefined", name), "")
	}
	return val
}

func main() {
	flag.Parse()

	setupLogger()
	printVersion()

	dvcrSettings, err := appconfig.LoadDVCRSettingsFromEnvs(controllerNamespace)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	leaderElectionNS := os.Getenv(common.PodNamespaceVar)
	if leaderElectionNS == "" {
		leaderElectionNS = "default"
	}

	// Setup scheme for all resources
	scheme := apiruntime.NewScheme()
	for _, f := range resourcesSchemeFuncs {
		err := f(scheme)
		if err != nil {
			log.Error(err, "Failed to add to scheme")
			os.Exit(1)
		}
	}

	managerOpts := manager.Options{
		// This controller watches resources in all namespaces.
		LeaderElection:             true,
		LeaderElectionNamespace:    leaderElectionNS,
		LeaderElectionID:           "d8-virt-operator-leader-election-helper",
		LeaderElectionResourceLock: "leases",
		Scheme:                     scheme,
	}

	vmCIDRsRaw := os.Getenv(common.VirtualMachineCIDRs)
	if vmCIDRsRaw == "" {
		log.Error(errors.New("vmCIDRs not found, but required"), "Failed to get vmCIDRs")
		os.Exit(1)
	}
	vmCIDRs := strings.Split(vmCIDRsRaw, ",")

	// Create a new Manager to provide shared dependencies and start components
	mgr, err := manager.New(cfg, managerOpts)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup context to gracefully handle termination.
	ctx := signals.SetupSignalHandler()

	if _, err := controller.NewVMDController(ctx, mgr, log, importerImage, uploaderImage, dvcrSettings); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if _, err := controller.NewCVMIController(ctx, mgr, log, importerImage, uploaderImage, controllerNamespace, dvcrSettings); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if _, err := controller.NewVMIController(ctx, mgr, log, importerImage, uploaderImage, dvcrSettings); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if _, err := controller.NewVMController(ctx, mgr, log, dvcrSettings); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if _, err := controller.NewVMBDAController(ctx, mgr, log); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if _, err := ipam.NewClaimController(ctx, mgr, log, vmCIDRs); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if _, err := ipam.NewLeaseController(ctx, mgr, log); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Starting the Manager.")

	// Start the Manager
	if err := mgr.Start(ctx); err != nil {
		log.Error(err, "manager exited non-zero")
		os.Exit(1)
	}
}