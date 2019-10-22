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

// const allPlacesHash = 'all_places_2356';
// const placeIdsHash = 'all_parent_IDS_for_places';

// const places: any = [];

// const getUnique = (data: any, keyToCheck: string) => {
//   return data.filter((obj: any, pos: any, arr: any) => {
//     return arr.map((mapObj: any) => mapObj[keyToCheck])
//     .indexOf(obj[keyToCheck]) === pos;
//   });
// };

// const placesData = (): Promise<object>  => {
//   return new Promise((resolve, reject) => {
//     return fs.readFile('places.json', 'utf8', (err: any, readResponse: any) => {
//       if (err) {
//         reject(err);
//       }
//       const readData = JSON.parse(readResponse);
//       readData.forEach((element: any) => {
//         const parsePath = element.pathforward.split('\\').filter((elem: string) => String(elem));
//         const pathForward = parsePath[parsePath.length - 1];
//         let parent: any = {};
//         parent.parentPlaceID = element.parentplaceid;
//         parent.pathForward = pathForward;
//         if (element.geo && Object.keys(element.geo).length && parent.parentPlaceID) {
//           places.push(parent);
//           parent = {};
//         }
//       });
//       const removedDuplicates = getUnique(places, 'pathForward');
//       const sortData = removedDuplicates.sort((indexA: PlaceIDObject, indexB: PlaceIDObject) => {
//         if (indexA.pathForward > indexB.pathForward) {
//           return 1;
//         } else {
//           return -1;
//         }
//       });
//       addDataToRedis(placeIdsHash, sortData);
//       addDataToRedis(allPlacesHash, readData);
//       resolve({
//         allPlacesData: readData,
//         allPlacesIds: places,
//       });
//     });
//   });
// };
