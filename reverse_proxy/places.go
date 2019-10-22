package reverse_proxy

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (app App) GetPlaceIds(w http.ResponseWriter, r *http.Request) {
	var allPlacesIdsRedis []PlaceID

	from := r.URL.Query().Get("from")
	size := r.URL.Query().Get("size")

	parseFrom, err := strconv.Atoi(from)
	if err != nil {
		log.Println(err, "Error parsing from to number", from, size)
	}

	parseSize, err := strconv.Atoi(size)
	if err != nil {
		log.Println(err, "Error parsing size to number")
	}

	redisHash := "all_parent_IDS_for_places"
	placesDataFromRedis := app.GetDataFromRedis(redisHash)

	redisDataToByte := []byte(placesDataFromRedis)
	if err = json.Unmarshal(redisDataToByte, &allPlacesIdsRedis); err != nil {
		log.Println(err, "PlaceIds retrieval from redis")
	}

	if allPlacesIdsRedis == nil {
		// app.AddDataToRedis(redisHash, "")
	}

	partitionData := allPlacesIdsRedis[parseFrom : parseSize+parseFrom]

	app.RespondWithJSON(w, http.StatusOK,
		PlaceIdPayload{
			Error:   false,
			Message: "Harvard art place IDs retrieved successfully",
			Places:  partitionData,
		},
	)

	return
}

func (app App) GetPlaces(w http.ResponseWriter, r *http.Request) {
	var allPlaces []Place
	var filteredPlaces []Place
	var err error

	placeID := r.URL.Query().Get("placeId")
	parsePlaceID, err := strconv.Atoi(placeID)
	if err != nil {
		log.Println(err, "Error parsing placeId to number")
	}

	redisHash := "all_places_2356"
	placesDataFromRedis := app.GetDataFromRedis(redisHash)

	redisDataToByte := []byte(placesDataFromRedis)
	if err = json.Unmarshal(redisDataToByte, &allPlaces); err != nil {
		log.Println(err, "Places data retrieval from redis")
	}

	if allPlaces == nil {
		// app.AddDataToRedis(redisHash, "")
	}

	if placeID != "" {
		for _, place := range allPlaces {
			if place.ParentPlaceID == parsePlaceID {
				filteredPlaces = append(filteredPlaces, place)
			}
		}

		allPlaces = filteredPlaces
	}

	if len(allPlaces) > 100 {
		allPlaces = allPlaces[0:100]
	}

	app.RespondWithJSON(w, http.StatusOK,
		PlacesPayload{
			Error:   false,
			Message: "Harvard art place IDs retrieved successfully",
			Places:  allPlaces,
		},
	)

	return
}
