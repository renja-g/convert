package image

import (
	"image"
	_ "image/png" // register PNG decoder
	"os"
	"path/filepath"

	"github.com/renja-g/convert/internal/converter"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/spf13/pflag"
)

type pngToWebpConverter struct{}

func (c *pngToWebpConverter) From() string { return ".png" }
func (c *pngToWebpConverter) To() string   { return ".webp" }

func (c *pngToWebpConverter) Convert(inputPath, outputPath string, options converter.Options) error {
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
		outputPath = inputPath[:len(inputPath)-len(ext)] + c.To()
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	encOptions, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return err
	}

	return webp.Encode(outputFile, img, encOptions)
}

func (c *pngToWebpConverter) GetFlags() *pflag.FlagSet {
	return pflag.NewFlagSet("png-to-webp", pflag.ExitOnError)
}

func init() {
	converter.Register(&pngToWebpConverter{})
}
