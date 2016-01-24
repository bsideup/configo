package parsers

import (
	"github.com/magiconair/properties"
)

type PropertiesParser struct{}

func (propertiesParser *PropertiesParser) Parse(content []byte, c map[string]interface{}) error {

	var p *properties.Properties
	p, err := properties.Load(content, properties.UTF8)
	if err != nil {
		return err
	}
	for _, key := range p.Keys() {
		value, _ := p.Get(key)
		c[key] = value
	}
	return err
}
