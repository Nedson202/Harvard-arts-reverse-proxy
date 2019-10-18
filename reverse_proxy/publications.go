package reverse_proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (app App) GetPublications(w http.ResponseWriter, r *http.Request) {
	var response CollectionsResponse

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	year := r.URL.Query().Get("year")

	redisHash := fmt.Sprintf("size %s - page %s - publications", size, page)
	publicationsFromRedis := app.GetDataFromRedis(redisHash)

	redisDataToByte := []byte(publicationsFromRedis)
	err := json.Unmarshal(redisDataToByte, &response)
	if err != nil {
		log.Println(err, "Data retrieval from redis")
	}

	if response.Records == nil {
		requestURL := fmt.Sprintf("publication?q=publicationyear:%s&apikey=%s&page=%s&size=%s", year, app.harvardAPIKey, page, size)

		result := app.RetrieveDataFromHarvardAPI(requestURL)
		err := json.Unmarshal(result, &response)
		if err != nil {
			log.Println(err)

			return
		}

		app.AddDataToRedis(redisHash, result)
	}

	randomizedData := app.RandomizeData(response.Records)

	app.RespondWithJSON(w, http.StatusOK,
		PublicationsPayload{
			Error:        false,
			Message:      "Harvard art publications retrieved successfully",
			Publications: randomizedData,
		},
	)

	return
}
