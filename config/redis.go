package config

import (
	"github.com/redis/go-redis/v9"
)

func ProvideRedis(env *EnvironmentVariables) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.REDIS_ADDR,
		Password: env.REDIS_PASSWORD, // no password set
		DB:       0,                  // use default DB
	})

	return rdb
}
