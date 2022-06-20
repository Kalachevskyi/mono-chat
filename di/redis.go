package di

import (
	"fmt"

	"github.com/go-redis/redis"
)

// RedisClient - initialize the instance of redis client.
func RedisClient(url string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("can't initialize Redis client: err=%s", err.Error())
	}

	return client, nil
}
