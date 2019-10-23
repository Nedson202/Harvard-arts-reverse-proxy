package reverse_proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	stringAction "strings"
)

func (p *PlaceID) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Place) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

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
		places, _ := app.ReadPlacesData()
		allPlacesIdsRedis = places
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
	parsePlaceID, err := strconv.ParseInt(placeID, 10, 64)
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
		_, parsedData := app.ReadPlacesData()
		allPlaces = parsedData
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

func (app App) ReadPlacesData() (places []PlaceID, parsedData []Place) {
	const allPlacesHash = "all_places_2356"
	const placeIdsHash = "all_parent_IDS_for_places"

	data, err := ioutil.ReadFile("places.json")
	if err != nil {
		log.Println("File reading error", err)

		return
	}

	dataToByte := []byte(data)
	if err := json.Unmarshal(dataToByte, &parsedData); err != nil {
		log.Println("Unable to parse json file data", err)
	}

	hashIDs := make(map[string]interface{})

	for _, place := range parsedData {
		var pathForward string

		parsePath := stringAction.Split(place.PathForward, "\\")

		trimPath := parsePath[0 : len(parsePath)-1]
		if len(trimPath) > 0 {
			lastIndex := len(trimPath) - 1
			pathForward = trimPath[lastIndex]
		}

		parent := PlaceID{place.ParentPlaceID, pathForward}

		hasGeolocationProperty := false
		hasParentIDProperty := false

		if place.Geo.Lat != 0 {
			hasGeolocationProperty = true
		}

		if place.Geo.Lon != 0 {
			hasGeolocationProperty = true
		}

		if parent.ParentPlaceId != 0 {
			hasParentIDProperty = true
		}

		if !hasGeolocationProperty || !hasParentIDProperty {
			continue
		}

		_, ok := hashIDs[pathForward]

		if !ok {
			places = append(places, parent)
		}

		hashIDs[pathForward] = place
	}

	app.AddDataToRedis(placeIdsHash, places)
	app.AddDataToRedis(allPlacesHash, parsedData)

	return
}
