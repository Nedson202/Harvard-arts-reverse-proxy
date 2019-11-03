package reverse_proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (app App) GetPublications(w http.ResponseWriter, r *http.Request) {
	var response PublicationsResponse

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	year := r.URL.Query().Get("year")

	redisHash := fmt.Sprintf("size %s - page %s - publications - year %s", size, page, year)
	publicationsFromRedis := app.GetDataFromRedis(redisHash)

	redisDataToByte := []byte(publicationsFromRedis)
	if err := json.Unmarshal(redisDataToByte, &response); err != nil {
		log.Println(err, "Data retrieval from redis")
	}

	if response.Records == nil {
		requestURL := fmt.Sprintf("publication?q=publicationyear:%s&apikey=%s&page=%s&size=%s", year, app.harvardAPIKey, page, size)

		result := app.RetrieveDataFromHarvardAPI(requestURL)
		if err := json.Unmarshal(result, &response); err != nil {
			log.Println(err)

			return
		}

		app.AddDataToRedis(redisHash, result)
	}

	app.RespondWithJSON(w, http.StatusOK,
		PublicationsPayload{
			Error:        false,
			Message:      "Harvard art publications retrieved successfully",
			Publications: response.Records,
		},
	)

	return
}
