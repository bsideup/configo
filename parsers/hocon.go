package parsers

import (
	"fmt"
	hocon "github.com/go-akka/configuration"
	hoconTokens "github.com/go-akka/configuration/hocon"
)

type HOCONParser struct{}

func (hoconParser *HOCONParser) Parse(content []byte, c map[string]interface{}) error {
	root := hocon.ParseString(string(content)).Root()

	if !root.IsObject() {
		return fmt.Errorf("Root must be an object")
	}

	for key, value := range toObj(root).(map[string]interface{}) {
		c[key] = value
	}

	return nil
}

// TODO https://github.com/go-akka/configuration/issues/2
func toObj(p *hoconTokens.HoconValue) interface{} {
	if p.IsString() {
		result := p.GetString()

		switch result {
		case "true":
			return true
		case "false":
			return false
		default:
			return result
		}
	}

	if p.IsObject() {
		obj := p.GetObject()
		result := make(map[string]interface{})

		for _, key := range obj.GetKeys() {
			result[key] = toObj(obj.GetKey(key))
		}

		return result
	}

	if p.IsArray() {
		var result []interface{}
		for _, item := range p.GetArray() {
			result = append(result, toObj(item))
		}
		return result
	}

	panic("token is not a String, Object or Array")
}
