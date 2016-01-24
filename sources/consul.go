package sources

import (
	consul "github.com/hashicorp/consul/api"
	"strings"
)

type ConsulSource struct {
	Address string `json:"address"`
	Prefix  string `json:"prefix"`
	Scheme  string `json:"scheme"`
}

func (consulSource *ConsulSource) Get() (map[string]interface{}, error) {
	config := consul.DefaultConfig()

	config.Address = consulSource.Address

	if consulSource.Scheme != "" {
		config.Scheme = consulSource.Scheme
	}

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	pairs, _, err := client.KV().List(consulSource.Prefix, nil)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, pair := range pairs {
		parts := strings.Split(pair.Key, "/")
		result[parts[len(parts)-1]] = string(pair.Value)
	}

	return result, nil
}
