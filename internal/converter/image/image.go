// Package image provides image format converters.
// It imports all subpackages to ensure converter registration.
package image

import (
	// Import all converter subpackages for side-effect of registration
	_ "github.com/renja-g/convert/internal/converter/image/jpeg"
	_ "github.com/renja-g/convert/internal/converter/image/png"
	_ "github.com/renja-g/convert/internal/converter/image/webp"
)
