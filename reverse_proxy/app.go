package reverse_proxy

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/nedson202/harvard-arts-reverse-proxy/redis"
)

func New(redisHost string, router *mux.Router, baseURL string, harvardAPIKey string) (app App, err error) {
	app.baseURL = baseURL
	app.harvardAPIKey = harvardAPIKey

	app.NewRouter(router)
	// Connect redis client
	result, err := redis.ConnectClient(redisHost)
	if err != nil {
		log.Println(err)

		return
	}
	log.Println(result)

	return
}
