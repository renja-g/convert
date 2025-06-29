package jpeg

import (
	"image"
	_ "image/jpeg" // register JPEG decoder
	"os"
	"path/filepath"

	"github.com/renja-g/convert/internal/converter"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/spf13/pflag"
)

type jpegToWebpConverter struct{}

func (c *jpegToWebpConverter) From() string { return ".jpeg" }
func (c *jpegToWebpConverter) To() string   { return ".webp" }

func (c *jpegToWebpConverter) Convert(inputPath, outputPath string, options converter.Options) error {
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

func (c *jpegToWebpConverter) GetFlags() *pflag.FlagSet {
	return pflag.NewFlagSet("jpeg-to-webp", pflag.ExitOnError)
}

func init() {
	converter.Register(&jpegToWebpConverter{})
	// also register alias .jpg
	converter.Register(&jpegToWebpConverterAliasJpg{})
}

type jpegToWebpConverterAliasJpg struct {
	jpegToWebpConverter
}

func (c *jpegToWebpConverterAliasJpg) From() string { return ".jpg" }
