package parsers

import (
	"strings"
)

// Parser defines methods a source should implement, to parse a config from
// array of bytes
type Parser interface {

	// Parse parses a config using the Parser's algorithm
	Parse(content []byte, c map[string]interface{}) error
}

// MustGetParser gets parser by format. It will call panic() if format is unknown
func MustGetParser(format string) Parser {
	parser := GetParser(format)

	if parser == nil {
		panic("Unknown format: " + format)
	}

	return parser
}

// GetParser gets parser by format or nil if it's unknown
func GetParser(format string) Parser {
	switch strings.ToLower(format) {
	case "yaml", "yml":
		return &YAMLParser{}
	case "json":
		return &JSONParser{}
	case "hcl":
		return &HCLParser{}
	case "toml":
		return &TOMLParser{}
	case "properties", "props", "prop":
		return &PropertiesParser{}
	}

	return nil
}
