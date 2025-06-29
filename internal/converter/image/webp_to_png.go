package image

import (
	"image/png"
	"os"
	"path/filepath"

	"github.com/renja-g/convert/internal/converter"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/spf13/pflag"
)

// webpToPngConverter converts WebP images to PNG.
// It uses the go-webp library to decode WebP and the standard library to encode PNG.
// Currently no specific options are exposed; the FlagSet is empty for future extension.

type webpToPngConverter struct{}

func (c *webpToPngConverter) From() string { return ".webp" }

func (c *webpToPngConverter) To() string { return ".png" }

func (c *webpToPngConverter) Convert(inputPath, outputPath string, options converter.Options) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	img, err := webp.Decode(inputFile, &decoder.Options{})
	if err != nil {
		return err
	}

	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = inputPath[:len(inputPath)-len(ext)] + c.To()
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return png.Encode(outputFile, img)
}

func (c *webpToPngConverter) GetFlags() *pflag.FlagSet {
	return pflag.NewFlagSet("webp-to-png", pflag.ExitOnError)
}

func init() {
	converter.Register(&webpToPngConverter{})
}
