package sources

import (
	"github.com/hashicorp/vault/api"
)

type VaultSource struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Path    string `json:"path"`
}

func (vaultSource *VaultSource) Get() (map[string]interface{}, error) {
	config := api.DefaultConfig()
	config.Address = vaultSource.Address

	client, err := api.NewClient(config)

	if err != nil {
		return nil, err
	}

	client.SetToken(vaultSource.Token)

	secret, err := client.Logical().Read(vaultSource.Path)

	if err != nil {
		return nil, err
	}

	return secret.Data, nil

}
