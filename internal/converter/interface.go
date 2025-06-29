package converter

import "github.com/spf13/pflag"

// Options contains common conversion options.
type Options map[string]interface{}

// Converter defines the contract for any file converter.
type Converter interface {
	// From returns the source file extension (e.g., ".jpeg").
	From() string
	// To returns the destination file extension (e.g., ".png").
	To() string
	// Convert performs the file conversion.
	Convert(inputPath string, outputPath string, options Options) error
	// GetFlags returns a set of `pflag.FlagSet` for this converter's specific options.
	GetFlags() *pflag.FlagSet
}
