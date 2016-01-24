package main

import (
	"encoding/json"
	"fmt"
	"github.com/zeroturnaround/configo/sources"
	"reflect"
)

// Source defines methods a source should implement, to resolve a config
type Source interface {

	// Get resolves a config using the Source's config
	Get() (map[string]interface{}, error)
}

type configDTO struct {
	Type string `json:"type"`
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
	configObject := configDTO{}

	sourceBytes := []byte(source)
	if err := json.Unmarshal(sourceBytes, &configObject); err != nil {
		return nil, err
	}

	sourceType, found := configMappings[configObject.Type]
	if !found {
		return nil, fmt.Errorf("Failed to find source type: %s", configObject.Type)
	}

	loader := reflect.New(sourceType).Interface().(Source)
	if err := json.Unmarshal(sourceBytes, loader); err != nil {
		return nil, err
	}

	return loader.Get()
}
