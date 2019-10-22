package reverse_proxy

import (
	"github.com/gorilla/mux"
)

func New(redisHost string, router *mux.Router, baseURL string, harvardAPIKey string, elasticSearchURL string) (app App, err error) {
	app.baseURL = baseURL
	app.harvardAPIKey = harvardAPIKey

	// Connect redis client
	client := ConnectRedisClient(redisHost)
	app.redisClient = client

	elasticSearchClient, err := ConnectElasticClient(elasticSearchURL)
	if err == nil {
		app.elasticClient = elasticSearchClient
	}

	app.NewRouter(router)

	return
}
