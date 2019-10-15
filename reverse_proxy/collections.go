package reverse_proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (app App) GetCollections(w http.ResponseWriter, r *http.Request) {
	var response CollectionsResponse
	requestURL := fmt.Sprintf("object?apikey=%s", app.harvardAPIKey)
	// requestURL := `object?apikey=${HARVARD_API_KEY}&hasimage=1&size=${size}&page=${page}`

	result := app.RetrieveDataFromHarvardAPI(requestURL)
	err := json.Unmarshal(result, &response)
	if err != nil {
		log.Println(err)

		return
	}

	app.RespondWithJSON(w, http.StatusOK,
		DataPayload{
			Error:   false,
			Message: "Collection objects retrieved successfully",
			Data:    response.Records,
		},
	)

	return
}

// public static async getObjects(req: Request, res: Response, next: NextFunction) {
//   try {

//     let objects: any;
//     const { size, page } = req.query;
//     const redisHash = `objects - size ${size} - page ${page}`;
//     objects = await getDataFromRedis(redisHash);

//     if (!objects) {
//       const url = `object?apikey=${HARVARD_API_KEY}&hasimage=1&size=${size}&page=${page}`;
//       objects = await ArtsController.retrieveDataFromHarvardAPI(url);
//       addDataToRedis(redisHash, objects);
//       elasticBulkCreate(objects.records);
//     }
//     const randomizedData = ArtsController.dataRandomizer(objects.records);

//     if (objects) {
//       return res.status(200).json({
//         error: false,
//         message: 'Harvard art objects retrieved successfully',
//         records: randomizedData,
//       });
//     }
//   } catch (error) {
//     error.httpStatusCode = 500;
//     return next(error);
//   }
// }

// public static async getObject(req: Request, res: Response, next: NextFunction) {
//   try {
//     const { objectId } = req.params;
//     let queryResponse: any = '';
//     const redisHash = `objectId-${objectId}`;
//     queryResponse = await getDataFromRedis(redisHash);

//     if (!queryResponse) {
//       const url = `object/${objectId}?apikey=${HARVARD_API_KEY}`;
//       queryResponse = await ArtsController.retrieveDataFromHarvardAPI(url);
//       addDocument(queryResponse, 'art');
//       addDataToRedis(redisHash, queryResponse);
//     }

//     if (queryResponse) {
//       return res.status(200).json({
//         error: false,
//         message: 'Harvard art object retrieved successfully',
//         record: queryResponse,
//       });
//     }
//   } catch (error) {
//     error.httpStatusCode = 500;
//     return next(error);
//   }
// }
