package parsers

import (
	"encoding/json"
)

type JSONParser struct{}

func (jsonParser *JSONParser) Parse(content []byte, c map[string]interface{}) error {

	return json.Unmarshal(content, &c)
}
