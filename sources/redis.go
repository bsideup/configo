package sources

import (
	"github.com/garyburd/redigo/redis"
)

type RedisSource struct {
	URI string `json:"uri"`
	Key string `json:"key"`
}

func (redisSource *RedisSource) Get() (map[string]interface{}, error) {
	connection, err := redis.DialURL(redisSource.URI)

	if err != nil {
		return nil, err
	}

	reply, err := redis.StringMap(connection.Do("HGETALL", redisSource.Key))

	result := make(map[string]interface{})

	for key, value := range reply {
		result[key] = value
	}

	return result, nil
}
