package png

import (
	"image"
	"image/jpeg"
	_ "image/png" // import for side-effect of registering png decoder
	"os"
	"path/filepath"

	"github.com/renja-g/convert/internal/converter"
	"github.com/spf13/pflag"
)

type pngToJpegConverter struct{}

func (c *pngToJpegConverter) From() string {
	return ".png"
}

func (c *pngToJpegConverter) To() string {
	return ".jpeg"
}

func (c *pngToJpegConverter) Convert(inputPath string, outputPath string, options converter.Options) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return err
	}

	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = inputPath[0:len(inputPath)-len(ext)] + c.To()
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// jpeg.Encode options can be provided here, for now using defaults.
	return jpeg.Encode(outputFile, img, nil)
}

func (c *pngToJpegConverter) GetFlags() *pflag.FlagSet {
	// We could add flags for JPEG quality here in the future
	return pflag.NewFlagSet("png-to-jpeg", pflag.ExitOnError)
}

func init() {
	converter.Register(&pngToJpegConverter{})
}
