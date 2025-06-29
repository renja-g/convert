package image

import (
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/renja-g/convert/internal/converter"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/spf13/pflag"
)

type webpToJpegConverter struct{}

func (c *webpToJpegConverter) From() string { return ".webp" }
func (c *webpToJpegConverter) To() string   { return ".jpeg" }

func (c *webpToJpegConverter) Convert(inputPath, outputPath string, options converter.Options) error {
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

	return jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 75})
}

func (c *webpToJpegConverter) GetFlags() *pflag.FlagSet {
	return pflag.NewFlagSet("webp-to-jpeg", pflag.ExitOnError)
}

func init() {
	converter.Register(&webpToJpegConverter{})
}
