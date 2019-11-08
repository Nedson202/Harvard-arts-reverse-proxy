package reverse_proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (app App) getCollections(w http.ResponseWriter, r *http.Request) {
	var response CollectionsResponse
	var err error

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")

	redisHash := fmt.Sprintf("objects - size %s - page %s", size, page)
	collectionsFromRedis := app.getDataFromRedis(redisHash)

	if err = json.Unmarshal([]byte(collectionsFromRedis), &response); err != nil {
		log.Println(err)
	}

	if response.Records == nil {
		requestURL := fmt.Sprintf("object?apikey=%s&hasimage=1&size=%s&page=%s", app.harvardAPIKey, size, page)

		result := app.retrieveDataFromHarvardAPI(requestURL)
		if err = json.Unmarshal(result, &response); err != nil {
			log.Println(err)

			return
		}

		app.addDataToRedis(redisHash, result)
		app.elasticBulkWrite(response.Records)
	}

	app.respondWithJSON(w, http.StatusOK,
		RecordsPayload{
			Error:   false,
			Message: "Harvard art objects retrieved successfully",
			Records: response.Records,
		},
	)

	return
}

func (app App) getCollection(w http.ResponseWriter, r *http.Request) {
	var response interface{}
	var err error

	params := mux.Vars(r)
	objectID := params["objectId"]

	redisHash := fmt.Sprintf("objectId-%s", objectID)
	objectFromRedis := app.getDataFromRedis(redisHash)

	if err := json.Unmarshal([]byte(objectFromRedis), &response); err != nil {
		log.Println(err)
	}

	if response != nil {
		app.respondWithJSON(w, http.StatusOK,
			RecordPayload{
				Error:   false,
				Message: "Harvard art objects retrieved successfully",
				Record:  response,
			},
		)

		return
	}

	requestURL := fmt.Sprintf("object/%s?apikey=%s", objectID, app.harvardAPIKey)

	result := app.retrieveDataFromHarvardAPI(requestURL)
	if err = json.Unmarshal(result, &response); err != nil {
		log.Println(err)

		return
	}

	app.addDataToRedis(redisHash, result)

	app.respondWithJSON(w, http.StatusOK,
		RecordPayload{
			Error:   false,
			Message: "Harvard art objects retrieved successfully",
			Record:  response,
		},
	)

	return
}
