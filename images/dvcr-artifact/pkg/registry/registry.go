package registry

import (
	"archive/tar"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/djherbis/buffer"
	"github.com/djherbis/nio/v3"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/stream"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"k8s.io/klog/v2"
	"kubevirt.io/containerized-data-importer/pkg/importer"
	"kubevirt.io/containerized-data-importer/pkg/util"

	"github.com/deckhouse/virtualization-controller/dvcr-importers/pkg/datasource"
	"github.com/deckhouse/virtualization-controller/dvcr-importers/pkg/monitoring"
)

// FIXME(ilya-lesikov): certdir

const (
	imageLabelSourceImageSize        = "source-image-size"
	imageLabelSourceImageVirtualSize = "source-image-virtual-size"
	imageLabelSourceImageFormat      = "source-image-format"
)

type ImageInfo struct {
	VirtualSize int    `json:"virtual-size"`
	Format      string `json:"format"`
}

const (
	imageInfoSize        = 64 * 1024 * 1024
	pipeBufSize          = 64 * 1024 * 1024
	tempImageInfoPattern = "tempfile"
	isoImageType         = "iso"
)

type DataProcessor struct {
	ds            datasource.DataSourceInterface
	destUsername  string
	destPassword  string
	destImageName string
	sha256Sum     string
	md5Sum        string
	destInsecure  bool
}

type DestinationRegistry struct {
	ImageName string
	Username  string
	Password  string
	Insecure  bool
}

func NewDataProcessor(ds datasource.DataSourceInterface, dest DestinationRegistry, sha256Sum, md5Sum string) (*DataProcessor, error) {
	return &DataProcessor{
		ds,
		dest.Username,
		dest.Password,
		dest.ImageName,
		sha256Sum,
		md5Sum,
		dest.Insecure,
	}, nil
}

func (p DataProcessor) Process(ctx context.Context) error {
	sourceImageFilename, err := p.ds.Filename()
	if err != nil {
		return fmt.Errorf("error getting source filename: %w", err)
	}

	sourceImageSize, err := p.ds.Length()
	if err != nil {
		return fmt.Errorf("error getting source image size: %w", err)
	}

	if sourceImageSize == 0 {
		return fmt.Errorf("zero data source image size")
	}

	sourceImageReader, err := p.ds.ReadCloser()
	if err != nil {
		return fmt.Errorf("error getting source image reader: %w", err)
	}

	// Wrap data source reader with progress and speed metrics.
	progressMeter := monitoring.NewProgressMeter(sourceImageReader, uint64(sourceImageSize))
	progressMeter.Start()
	defer progressMeter.Stop()

	pipeReader, pipeWriter := nio.Pipe(buffer.New(pipeBufSize))
	imageInfoCh := make(chan ImageInfo)
	errsGroup, ctx := errgroup.WithContext(ctx)
	errsGroup.Go(func() error {
		return p.inspectAndStreamSourceImage(ctx, sourceImageFilename, sourceImageSize, progressMeter, pipeWriter, imageInfoCh)
	})
	errsGroup.Go(func() error {
		defer pipeReader.Close()
		return p.uploadLayersAndImage(ctx, pipeReader, sourceImageSize, progressMeter, imageInfoCh)
	})

	return errsGroup.Wait()
}

func (p DataProcessor) inspectAndStreamSourceImage(
	ctx context.Context,
	sourceImageFilename string,
	sourceImageSize int,
	sourceImageReader io.ReadCloser,
	pipeWriter *nio.PipeWriter,
	imageInfoCh chan ImageInfo,
) error {
	var tarWriter *tar.Writer
	{
		tarWriter = tar.NewWriter(pipeWriter)
		header := &tar.Header{
			Name:     path.Join("disk", sourceImageFilename),
			Size:     int64(sourceImageSize),
			Mode:     0o644,
			Typeflag: tar.TypeReg,
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("error writing tar header: %w", err)
		}
	}

	var checksumWriters []io.Writer
	var checksumCheckFuncList []func() error
	{
		if p.sha256Sum != "" {
			hash := sha256.New()
			checksumWriters = append(checksumWriters, hash)
			checksumCheckFuncList = append(checksumCheckFuncList, func() error {
				sum := hex.EncodeToString(hash.Sum(nil))
				if sum != p.sha256Sum {
					return fmt.Errorf("sha256 sum mismatch: %s != %s", sum, p.sha256Sum)
				}

				return nil
			})
		}

		if p.md5Sum != "" {
			hash := md5.New()
			checksumWriters = append(checksumWriters, hash)
			checksumCheckFuncList = append(checksumCheckFuncList, func() error {
				sum := hex.EncodeToString(hash.Sum(nil))
				if sum != p.md5Sum {
					return fmt.Errorf("md5 sum mismatch: %s != %s", sum, p.md5Sum)
				}

				return nil
			})
		}
	}

	var streamWriter io.Writer
	{
		writers := []io.Writer{tarWriter}
		writers = append(writers, checksumWriters...)
		streamWriter = io.MultiWriter(writers...)
	}

	errsGroup, ctx := errgroup.WithContext(ctx)

	imageInfoReader, imageInfoWriter := nio.Pipe(buffer.New(imageInfoSize))

	errsGroup.Go(func() error {
		defer tarWriter.Close()
		defer pipeWriter.Close()
		defer sourceImageReader.Close()
		defer imageInfoWriter.Close()

		klog.Infoln("Streaming from the source")
		doneSize, err := io.Copy(streamWriter, io.TeeReader(sourceImageReader, imageInfoWriter))
		if err != nil {
			return fmt.Errorf("error copying from the source: %w", err)
		}

		if doneSize != int64(sourceImageSize) {
			return fmt.Errorf("source image size mismatch: %d != %d", doneSize, sourceImageSize)
		}

		for _, checksumCheckFunc := range checksumCheckFuncList {
			if err = checksumCheckFunc(); err != nil {
				return err
			}
		}

		klog.Infoln("Source streaming completed")

		return nil
	})

	errsGroup.Go(func() error {
		defer imageInfoReader.Close()

		info, err := getImageInfo(ctx, imageInfoReader)
		if err != nil {
			return err
		}

		imageInfoCh <- info

		return nil
	})

	return errsGroup.Wait()
}

func (p DataProcessor) uploadLayersAndImage(
	ctx context.Context,
	pipeReader *nio.PipeReader,
	sourceImageSize int,
	progressMeter *monitoring.ProgressMeter,
	imageInfoCh chan ImageInfo,
) error {
	nameOpts := destNameOptions(p.destInsecure)
	remoteOpts := destRemoteOptions(ctx, p.destUsername, p.destPassword, p.destInsecure)
	image := empty.Image

	ref, err := name.ParseReference(p.destImageName, nameOpts...)
	if err != nil {
		return fmt.Errorf("error parsing image name: %w", err)
	}

	repo, err := name.NewRepository(ref.Context().Name(), nameOpts...)
	if err != nil {
		return fmt.Errorf("error constructing new repository: %w", err)
	}

	layer := stream.NewLayer(pipeReader)

	klog.Infoln("Uploading layer to registry")
	if err := remote.WriteLayer(repo, layer, remoteOpts...); err != nil {
		return fmt.Errorf("error uploading layer: %w", err)
	}
	klog.Infoln("Layer uploaded")

	cnf, err := image.ConfigFile()
	if err != nil {
		return fmt.Errorf("error getting image config: %w", err)
	}

	imageInfo := <-imageInfoCh

	klog.Infof("Got image info: %+v", imageInfo)

	cnf.Config.Labels = map[string]string{}
	cnf.Config.Labels[imageLabelSourceImageVirtualSize] = fmt.Sprintf("%d", imageInfo.VirtualSize)
	cnf.Config.Labels[imageLabelSourceImageSize] = fmt.Sprintf("%d", sourceImageSize)
	cnf.Config.Labels[imageLabelSourceImageFormat] = imageInfo.Format

	image, err = mutate.ConfigFile(image, cnf)
	if err != nil {
		return fmt.Errorf("error mutating image config: %w", err)
	}

	image, err = mutate.AppendLayers(image, layer)
	if err != nil {
		return fmt.Errorf("error appending layer to image: %w", err)
	}

	klog.Infof("Uploading image %q to registry", p.destImageName)
	if err = remote.Write(ref, image, remoteOpts...); err != nil {
		return fmt.Errorf("error uploading image: %w", err)
	}

	if err = WriteImportCompleteMessage(uint64(sourceImageSize), uint64(imageInfo.VirtualSize), progressMeter.GetAvgSpeed(), imageInfo.Format); err != nil {
		return fmt.Errorf("error writing import complete message: %w", err)
	}

	return nil
}

func getImageInfo(ctx context.Context, sourceReader io.ReadCloser) (ImageInfo, error) {
	formatSourceReaders, err := importer.NewFormatReaders(sourceReader, 0)
	if err != nil {
		return ImageInfo{}, fmt.Errorf("error creating format readers: %w", err)
	}

	var uncompressedN int64
	var tempImageInfoFile *os.File

	klog.Infoln("Write image info to temp file")
	{
		tempImageInfoFile, err = os.CreateTemp("", tempImageInfoPattern)
		if err != nil {
			return ImageInfo{}, fmt.Errorf("error creating temp file: %w", err)
		}

		uncompressedN, err = io.CopyN(tempImageInfoFile, formatSourceReaders.TopReader(), imageInfoSize)
		if err != nil && !errors.Is(err, io.EOF) {
			return ImageInfo{}, fmt.Errorf("error writing to temp file: %w", err)
		}

		if err = tempImageInfoFile.Close(); err != nil {
			return ImageInfo{}, fmt.Errorf("error closing temp file: %w", err)
		}
	}

	klog.Infoln("Get image info from temp file")
	var imageInfo ImageInfo
	{
		cmd := exec.CommandContext(ctx, "qemu-img", "info", "--output=json", tempImageInfoFile.Name())
		rawOut, err := cmd.Output()
		if err != nil {
			return ImageInfo{}, fmt.Errorf("error running qemu-img info: %w", err)
		}

		klog.Infoln("Qemu-img command output:", string(rawOut))

		if err = json.Unmarshal(rawOut, &imageInfo); err != nil {
			return ImageInfo{}, fmt.Errorf("error parsing qemu-img info output: %w", err)
		}

		if imageInfo.Format != "raw" {
			// It's necessary to read everything from the original image to avoid blocking.
			_, err = io.Copy(&util.EmptyWriter{}, sourceReader)
			if err != nil {
				return ImageInfo{}, fmt.Errorf("error copying to nowhere: %w", err)
			}

			return imageInfo, nil
		}
	}

	// `qemu-img` command does not support getting information about iso files.
	// It is necessary to obtain this information in another way (using the `file` command).
	klog.Infoln("Check the image as it may be an iso")
	{
		cmd := exec.CommandContext(ctx, "file", "-b", tempImageInfoFile.Name())
		rawOut, err := cmd.Output()
		if err != nil {
			return ImageInfo{}, fmt.Errorf("error running file info: %w", err)
		}

		out := string(rawOut)

		klog.Infoln("File command output:", out)

		if strings.HasPrefix(strings.ToLower(out), isoImageType) {
			imageInfo.Format = isoImageType
		}

		// Count uncompressed size of source image.
		n, err := io.Copy(&util.EmptyWriter{}, formatSourceReaders.TopReader())
		if err != nil {
			return ImageInfo{}, fmt.Errorf("error copying to nowhere: %w", err)
		}

		imageInfo.VirtualSize = int(uncompressedN + n)

		return imageInfo, nil
	}
}

type ImportInfo struct {
	SourceImageSize        uint64 `json:"source-image-size"`
	SourceImageVirtualSize uint64 `json:"source-image-virtual-size"`
	SourceImageFormat      string `json:"source-image-format"`
	AverageSpeed           uint64 `json:"average-speed"`
}

func WriteImportCompleteMessage(sourceImageSize, sourceImageVirtualSize, avgSpeed uint64, sourceImageFormat string) error {
	rawMsg, err := json.Marshal(ImportInfo{
		SourceImageSize:        sourceImageSize,
		SourceImageVirtualSize: sourceImageVirtualSize,
		SourceImageFormat:      sourceImageFormat,
		AverageSpeed:           avgSpeed,
	})
	if err != nil {
		return err
	}

	message := string(rawMsg)

	err = util.WriteTerminationMessage(message)
	if err != nil {
		return err
	}

	klog.Infoln("Image uploaded: " + message)

	return nil
}

func destNameOptions(destInsecure bool) []name.Option {
	nameOpts := []name.Option{}

	if destInsecure {
		nameOpts = append(nameOpts, name.Insecure)
	}

	return nameOpts
}

func destRemoteOptions(ctx context.Context, destUsername, destPassword string, destInsecure bool) []remote.Option {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: destInsecure,
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = tlsConfig

	remoteOpts := []remote.Option{
		remote.WithContext(ctx),
		remote.WithTransport(transport),
		remote.WithAuth(&authn.Basic{Username: destUsername, Password: destPassword}),
	}

	return remoteOpts
}