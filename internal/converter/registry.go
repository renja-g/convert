package converter

import "fmt"

var registry = make(map[string]map[string]Converter)

// Register adds a converter to the registry.
func Register(c Converter) {
	from := c.From()
	to := c.To()

	if _, ok := registry[from]; !ok {
		registry[from] = make(map[string]Converter)
	}
	if _, ok := registry[from][to]; ok {
		// This should be a panic because it's a developer error.
		// It means two converters are trying to handle the same conversion.
		panic(fmt.Sprintf("converter from %s to %s already registered", from, to))
	}
	registry[from][to] = c
}

// GetConvertersFor returns all available converters for a given source extension.
func GetConvertersFor(from string) map[string]Converter {
	return registry[from]
}

// GetConverter returns a specific converter for a source and destination extension.
func GetConverter(from, to string) (Converter, bool) {
	if converters, ok := registry[from]; ok {
		converter, ok := converters[to]
		return converter, ok
	}
	return nil, false
}
