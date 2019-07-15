package infrastructure

import (
	"fmt"

	"github.com/go-redis/redis"
)

// NewRedisClient - return a new instance of redis client,
// using library "github.com/go-redis/redis"
func NewRedisClient(url string) (*redis.Client, error) {
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
