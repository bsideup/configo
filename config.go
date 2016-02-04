package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/zeroturnaround/configo/sources"
	"reflect"
)

// Source defines methods a source should implement, to resolve a config
type Source interface {

	// Get resolves a config using the Source's config
	Get() (map[string]interface{}, error)
}

var configMappings = map[string]reflect.Type{
	"consul":   reflect.TypeOf(sources.ConsulSource{}),
	"dynamodb": reflect.TypeOf(sources.DynamoDBSource{}),
	"etcd":     reflect.TypeOf(sources.EtcdSource{}),
	"file":     reflect.TypeOf(sources.FileSource{}),
	"http":     reflect.TypeOf(sources.HTTPSource{}),
	"redis":    reflect.TypeOf(sources.RedisSource{}),
}

// GetConfig resolves config by source (in JSON format).
// Source must contain at least one property with name "type", which will be
// used to select proper source implementation.
func GetConfig(source string) (map[string]interface{}, error) {
	rawSource := make(map[string]interface{})

	sourceBytes := []byte(source)
	if err := json.Unmarshal(sourceBytes, &rawSource); err != nil {
		return nil, err
	}

	sourceType, found := configMappings[rawSource["type"].(string)]
	if !found {
		return nil, fmt.Errorf("Failed to find source type: %s", rawSource["type"])
	}
	delete(rawSource, "type")

	var metadata mapstructure.Metadata
	loader := reflect.New(sourceType).Interface().(Source)

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: &metadata,
		Result:   loader,
	})

	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(rawSource); err != nil {
		return nil, err
	}

	if len(metadata.Unused) > 0 {
		return nil, fmt.Errorf("unknown configuration keys: %v", metadata.Unused)
	}

	return loader.Get()
}
