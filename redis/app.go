package redis

import (
	"github.com/go-redis/redis/v7"
)

func ConnectClient(redisHost string) (pingResult string, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = client.Ping().Result()

	pingResult = "Redis client connected"

	return
}
