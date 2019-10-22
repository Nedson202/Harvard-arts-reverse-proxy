package reverse_proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (app App) GetCollections(w http.ResponseWriter, r *http.Request) {
	var response CollectionsResponse
	var err error

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")

	redisHash := fmt.Sprintf("objects - size %s - page %s", size, page)
	collectionsFromRedis := app.GetDataFromRedis(redisHash)

	redisDataToByte := []byte(collectionsFromRedis)
	if err = json.Unmarshal(redisDataToByte, &response); err != nil {
		log.Println(err)
	}

	if response.Records == nil {
		requestURL := fmt.Sprintf("object?apikey=%s&hasimage=1&size=%s&page=%s", app.harvardAPIKey, size, page)

		result := app.RetrieveDataFromHarvardAPI(requestURL)
		if err = json.Unmarshal(result, &response); err != nil {
			log.Println(err)

			return
		}

		app.AddDataToRedis(redisHash, result)
	}

	randomizedData := app.RandomizeData(response.Records)

	app.RespondWithJSON(w, http.StatusOK,
		RecordsPayload{
			Error:   false,
			Message: "Harvard art objects retrieved successfully",
			Records: randomizedData,
		},
	)

	return
}

func (app App) GetCollection(w http.ResponseWriter, r *http.Request) {
	var response interface{}
	var err error

	params := mux.Vars(r)
	objectID := params["objectId"]

	redisHash := fmt.Sprintf("objectId-%s", objectID)
	objectFromRedis := app.GetDataFromRedis(redisHash)

	redisDataToByte := []byte(objectFromRedis)
	if err := json.Unmarshal(redisDataToByte, &response); err != nil {
		log.Println(err)
	}

	if response != nil {
		app.RespondWithJSON(w, http.StatusOK,
			RecordPayload{
				Error:   false,
				Message: "Harvard art objects retrieved successfully",
				Record:  response,
			},
		)

		return
	}

	requestURL := fmt.Sprintf("object/%s?apikey=%s", objectID, app.harvardAPIKey)

	result := app.RetrieveDataFromHarvardAPI(requestURL)
	if err = json.Unmarshal(result, &response); err != nil {
		log.Println(err)

		return
	}

	app.AddDataToRedis(redisHash, result)

	app.RespondWithJSON(w, http.StatusOK,
		RecordPayload{
			Error:   false,
			Message: "Harvard art objects retrieved successfully",
			Record:  response,
		},
	)

	return
}
