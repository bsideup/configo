package sources

import (
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"
	"strings"
)

type EtcdSource struct {
	Endpoints  []string `json:"endpoints"`
	Prefix     string   `json:"prefix"`
	KeepPrefix bool     `json:"keepPrefix"`
}

func (etcdSource *EtcdSource) Get() (map[string]interface{}, error) {
	cfg := etcd.Config{
		Endpoints: etcdSource.Endpoints,
		Transport: etcd.DefaultTransport,
	}

	client, err := etcd.New(cfg)
	if err != nil {
		return nil, err
	}

	keysAPI := etcd.NewKeysAPI(client)

	response, err := keysAPI.Get(context.Background(), etcdSource.Prefix, &etcd.GetOptions{
		Recursive: true,
	})

	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	etcdSource.nodeToMap(response.Node, result)

	return result, nil
}

func (etcdSource *EtcdSource) nodeToMap(node *etcd.Node, result map[string]interface{}) {
	key := node.Key
	if !node.Dir {
		key = strings.TrimPrefix(key, "/")
		if !etcdSource.KeepPrefix {
			key = strings.TrimPrefix(key, etcdSource.Prefix)
		}
		key = strings.Replace(key, "/", "_", -1)
		result[key] = node.Value
	} else {
		for _, subNode := range node.Nodes {
			etcdSource.nodeToMap(subNode, result)
		}
	}
}
