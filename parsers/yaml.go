package parsers

import (
	"github.com/ghodss/yaml"
)

type YAMLParser struct{}

func (yamlParser *YAMLParser) Parse(content []byte, c map[string]interface{}) error {
	return yaml.Unmarshal(content, &c)
}
