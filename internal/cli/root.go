package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/renja-g/convert/internal/alias"
	"github.com/renja-g/convert/internal/converter"

	// Import the image package to register the converters
	_ "github.com/renja-g/convert/internal/converter/image"
	"github.com/renja-g/convert/internal/detect"
	"github.com/spf13/cobra"
)

var output string
var to string

var rootCmd = &cobra.Command{
	Use:   "convert [input file]",
	Short: "A universal file converter",
	Long:  `A universal file converter that supports various file formats.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]

		mimeType, err := detect.MimeType(inputFile)
		if err != nil {
			return fmt.Errorf("could not detect mime type: %w", err)
		}

		from, ok := detect.ExtensionFromMimeType(mimeType)
		if !ok {
			// fallback to extension
			from = strings.ToLower(filepath.Ext(inputFile))
		}

		if to != "" {
			// User has specified a target format.
			// The `to` flag should not include a dot, but our registry uses it.
			if !strings.HasPrefix(to, ".") {
				to = "." + to
			}

			to = alias.Resolve(to)

			c, found := converter.GetConverter(from, to)
			if !found {
				return fmt.Errorf("no converter found from %s to %s", from, to)
			}

			fmt.Printf("Converting %s to %s...\n", inputFile, to)
			return c.Convert(inputFile, output, nil)
		} else {
			// User has not specified a target format.
			// List available conversions.
			converters := converter.GetConvertersFor(from)
			if len(converters) == 0 {
				return fmt.Errorf("no converters found for %s", from)
			}

			fmt.Printf("Available conversions for %s:\n", from)
			for to_format := range converters {
				fmt.Printf("- %s\n", to_format)
			}
			fmt.Println("\nPlease specify a target format with the --to flag.")
			return nil
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output file path")
	rootCmd.PersistentFlags().StringVarP(&to, "to", "t", "", "Target format (e.g., png, jpg)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
