package parsers

import (
	"github.com/BurntSushi/toml"
)

type TOMLParser struct{}

func (tomlParser *TOMLParser) Parse(content []byte, c map[string]interface{}) error {
	_, err := toml.Decode(string(content), &c)

	return err
}
