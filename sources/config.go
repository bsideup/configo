package sources

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// Source defines methods a source should implement, to resolve a config
type Source interface {

	// Get resolves a config using the Source's config
	Get() (map[string]interface{}, error)
}

var configMappings = map[string]reflect.Type{
	"composite": reflect.TypeOf(CompositeSource{}),
	"consul":    reflect.TypeOf(ConsulSource{}),
	"dynamodb":  reflect.TypeOf(DynamoDBSource{}),
	"etcd":      reflect.TypeOf(EtcdSource{}),
	"file":      reflect.TypeOf(FileSource{}),
	"http":      reflect.TypeOf(HTTPSource{}),
	"redis":     reflect.TypeOf(RedisSource{}),
	"shell":     reflect.TypeOf(ShellSource{}),
	"vault":     reflect.TypeOf(VaultSource{}),
}

// GetSource resolves source from a map.
// Source must contain at least one property with name "type", which will be
// used to select proper source implementation.
func GetSource(rawSource map[string]interface{}) (Source, error) {
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

	return loader, nil
}
