package imaging

import (
	"context"
	"fmt"
	"io"

	"github.com/davidbyttow/govips/v2/vips"
)

type Service interface {
	Resize(ctx context.Context, imageData io.Reader) ([]byte, error)
}

func New() (Service, func()) {
	vips.Startup(nil)
	
	return service{}, func() {
		vips.Shutdown()
	}
}

type service struct {
}

func (svc service) Resize(ctx context.Context, imgData io.Reader) ([]byte, error) {
	inputImage, err := vips.NewImageFromReader(imgData)
	if err != nil {
		return nil, fmt.Errorf("load image error %w", err)
	}

	ep := vips.NewJpegExportParams()
	ep.StripMetadata = true
	ep.Quality = 75
	ep.Interlace = true
	ep.OptimizeCoding = true
	ep.SubsampleMode = vips.VipsForeignSubsampleAuto
	ep.TrellisQuant = true
	ep.OvershootDeringing = true
	ep.OptimizeScans = true
	ep.QuantTable = 3

	imageBytes, _, err := inputImage.ExportJpeg(ep)
	if err != nil {
		return nil, fmt.Errorf("export jpeg error %w", err)
	}

	return imageBytes, nil
}
