package persistence

import (
	"ddd-demo/common/factory"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(factory *factory.ObjectFactory, uri string) (*redis.Client, error) {
	client, err := factory.Get(
		uri,
		func(uri string) (interface{}, error) {
			opts, err := redis.ParseURL(uri)
			if err != nil {
				return nil, err
			}

			return redis.NewClient(opts), nil
		},
		func(obj interface{}) error {
			client, ok := obj.(*redis.Client)
			if !ok {
				return nil
			}

			return client.Close()
		},
	)
	if err != nil {
		return nil, err
	}

	return client.(*redis.Client), nil
}
