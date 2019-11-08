package reverse_proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (app App) getPublications(w http.ResponseWriter, r *http.Request) {
	var response PublicationsResponse

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	year := r.URL.Query().Get("year")

	redisHash := fmt.Sprintf("size %s - page %s - publications - year %s", size, page, year)
	publicationsFromRedis := app.getDataFromRedis(redisHash)

	if err := json.Unmarshal([]byte(publicationsFromRedis), &response); err != nil {
		log.Println(err, "Data retrieval from redis")
	}

	if response.Records == nil {
		requestURL := fmt.Sprintf("publication?q=publicationyear:%s&apikey=%s&page=%s&size=%s", year, app.harvardAPIKey, page, size)

		result := app.retrieveDataFromHarvardAPI(requestURL)
		if err := json.Unmarshal(result, &response); err != nil {
			log.Println(err)

			return
		}

		app.addDataToRedis(redisHash, result)
	}

	app.respondWithJSON(w, http.StatusOK,
		PublicationsPayload{
			Error:        false,
			Message:      "Harvard art publications retrieved successfully",
			Publications: response.Records,
		},
	)

	return
}
