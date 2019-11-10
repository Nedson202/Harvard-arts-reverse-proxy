package reverse_proxy

import (
	"log"
	s "strings"

	"github.com/go-redis/redis/v7"
)

func ConnectRedisClient(redisHost string) (client *redis.Client) {
	redisProdHost := s.HasPrefix(redisHost, "redis")

	if redisProdHost {
		redisHostOptions, err := redis.ParseURL(redisHost)
		if err != nil {
			log.Println(err, "Unable to parse redis host")
			return
		}

		client = redis.NewClient(redisHostOptions)
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: "",
			DB:       0,
		})
	}

	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Unable to establish connection to redis")

		return
	}

	log.Println("Redis client connected")

	return
}

func (app App) getDataFromRedis(redisHash string) (data string) {
	data, err := app.redisClient.Get(redisHash).Result()

	if err != nil {
		log.Println(err)

		return
	}

	return data
}

func (app App) addDataToRedis(redisHash string, data interface{}) {
	err := app.redisClient.Set(redisHash, data, 0).Err()

	if err != nil {
		log.Println(err)
	}

	return
}
