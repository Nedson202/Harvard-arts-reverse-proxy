package reverse_proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/estransport"
)

func ConnectElasticClient(elasticSearchURL string) (client *elasticsearch.Client, err error) {
	log.SetFlags(0)

	ctx := context.Background()

	indices := []string{"arts"}
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticSearchURL,
		},
		Logger: &estransport.ColorLogger{
			Output:            os.Stdout,
			EnableRequestBody: true,
		},
	}

	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	verifyOrCreateIndex(ctx, indices, client)

	log.Println("Elastic client connected")

	return
}

func verifyOrCreateIndex(ctx context.Context, indices []string, client *elasticsearch.Client) {
	req := esapi.IndicesExistsRequest{
		Index: indices,
	}

	indexExistsRes, err := req.Do(ctx, client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer indexExistsRes.Body.Close()

	if indexExistsRes.IsError() {
		req := esapi.IndicesCreateRequest{
			Index: "arts",
		}

		indexCreateRes, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}

		if !indexCreateRes.IsError() {
			log.Println("Index specified created successfully")
			createMapping(ctx, indices, client)
		}
	}
}

func createMapping(ctx context.Context, indices []string, client *elasticsearch.Client) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"mappings": map[string]interface{}{
			"art": map[string]interface{}{
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type": "text",
					},
					"id": map[string]interface{}{
						"type": "integer",
					},
					"accessionyear": map[string]interface{}{
						"type": "integer",
					},
					"suggest": map[string]interface{}{
						"type":            "completion",
						"analyzer":        "simple",
						"search_analyzer": "standard",
					},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	req := esapi.IndicesPutMappingRequest{
		Index: indices,
		Body:  &buf,
	}

	putMappingRes, err := req.Do(ctx, client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if !putMappingRes.IsError() {
		log.Println("Mapping added successfully")
	}
	log.Println(putMappingRes)
}

func (app App) SearchCollections(w http.ResponseWriter, r *http.Request) {
	var resultHits map[string]interface{}
	ctx := context.Background()

	query := r.URL.Query().Get("query")
	from := r.URL.Query().Get("from")
	size := r.URL.Query().Get("size")

	parseFrom, err := strconv.Atoi(from)
	if err != nil {
		log.Println(err, "Error parsing from to number")
	}

	parseSize, err := strconv.Atoi(size)
	if err != nil {
		log.Println(err, "Error parsing size to number")
	}

	fields := []string{
		"title", "century", "accessionmethod", "period", "technique",
		"classification", "department", "culture", "medium", "verificationleveldescription",
	}

	var buf bytes.Buffer
	multiMatchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"type":   "phrase_prefix",
				"fields": fields,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(multiMatchQuery); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	log.Println(parseFrom, "parsing size to number", parseSize, ctx)

	// Perform the search request.
	searchRes, err := app.elasticClient.Search(
		app.elasticClient.Search.WithContext(ctx),
		app.elasticClient.Search.WithIndex("arts"),
		app.elasticClient.Search.WithBody(&buf),
		app.elasticClient.Search.WithTrackTotalHits(true),
		app.elasticClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer searchRes.Body.Close()

	if searchRes.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(searchRes.Body).Decode(&e); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Printf("[%s] %s: %s",
				searchRes.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(searchRes.Body).Decode(&resultHits); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	var searchResults []interface{}
	for _, hit := range resultHits["hits"].(map[string]interface{})["hits"].([]interface{}) {
		searchResults = append(searchResults, hit.(map[string]interface{})["_source"])
	}

	app.RespondWithJSON(w, http.StatusOK,
		SearchResultsPayload{
			Error:   false,
			Message: "Harvard art objects retrieved successfully",
			Results: searchResults,
		},
	)

	return
}
