package di

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client //nolint:gochecknoglobals
	redisOnce   sync.Once     //nolint:gochecknoglobals
)

// RedisClient - initialize the instance of redis client
func RedisClient(url string) (*redis.Client, error) {
	var err error

	redisOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     url,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		if _, err = client.Ping().Result(); err != nil {
			err = fmt.Errorf("can't initialize Redis client: err=%s", err.Error())
			return
		}
		redisClient = client
	})

	return redisClient, err
}
