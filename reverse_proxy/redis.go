package reverse_proxy

import (
	"log"

	"github.com/go-redis/redis/v7"
)

func ConnectRedisClient(redisHost string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Unable to establish connection to redis")

		return
	}

	log.Println("Redis client connected")

	return
}

func (app App) GetDataFromRedis(redisHash string) (data string) {
	data, err := app.redisClient.Get(redisHash).Result()

	if err != nil {
		log.Println(err)

		return
	}

	return data
}

func (app App) AddDataToRedis(redisHash string, data interface{}) {
	err := app.redisClient.Set(redisHash, data, 0).Err()

	if err != nil {
		log.Println(err)
	}

	return
}
