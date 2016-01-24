package parsers

import (
	"github.com/hashicorp/hcl"
)

type HCLParser struct{}

func (hclParser *HCLParser) Parse(content []byte, c map[string]interface{}) error {
	return hcl.Decode(&c, string(content))
}
