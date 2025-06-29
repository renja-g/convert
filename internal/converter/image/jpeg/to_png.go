package jpeg

import (
	"image"
	_ "image/jpeg" // import for side-effect of registering jpeg decoder
	"image/png"
	"os"
	"path/filepath"

	"github.com/renja-g/convert/internal/converter"
	"github.com/spf13/pflag"
)

type jpegToPngConverter struct{}

func (c *jpegToPngConverter) From() string {
	return ".jpeg"
}

func (c *jpegToPngConverter) To() string {
	return ".png"
}

func (c *jpegToPngConverter) Convert(inputPath string, outputPath string, options converter.Options) error {
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

	return png.Encode(outputFile, img)
}

func (c *jpegToPngConverter) GetFlags() *pflag.FlagSet {
	// No specific flags for this conversion
	return pflag.NewFlagSet("jpeg-to-png", pflag.ExitOnError)
}

func init() {
	converter.Register(&jpegToPngConverter{})
	// Also handle .jpg
	converter.Register(&jpegToPngConverterAliasJpg{})
}

// Alias for .jpg
type jpegToPngConverterAliasJpg struct {
	jpegToPngConverter
}

func (c *jpegToPngConverterAliasJpg) From() string {
	return ".jpg"
}
