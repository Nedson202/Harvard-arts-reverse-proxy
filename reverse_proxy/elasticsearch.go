package reverse_proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
		log.Println("Error creating the client: %s", err)

		return
	}

	err = verifyOrCreateIndex(ctx, indices, client)
	if err == nil {
		log.Println("Elastic client connected")
	}

	return
}

func verifyOrCreateIndex(ctx context.Context, indices []string, client *elasticsearch.Client) (err error) {
	req := esapi.IndicesExistsRequest{
		Index: indices,
	}

	indexExistsRes, err := req.Do(ctx, client)
	if err != nil {
		log.Println("Error getting response: %s", err)

		return
	}
	defer indexExistsRes.Body.Close()

	if indexExistsRes.IsError() {
		req := esapi.IndicesCreateRequest{
			Index: "arts",
		}

		indexCreateRes, err := req.Do(ctx, client)
		if err != nil {
			log.Println("Error getting response: %s", err)
		}

		if !indexCreateRes.IsError() {
			log.Println("Index specified created successfully")
			createMapping(ctx, indices, client)
		}
	}

	return
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
		log.Println("Error getting response: %s", err)

		return
	}

	if !putMappingRes.IsError() {
		log.Println("Mapping added successfully")

		return
	}
}

func (app App) SearchCollections(w http.ResponseWriter, r *http.Request) {
	var resultHits map[string]interface{}
	ctx := context.Background()

	query := r.URL.Query().Get("query")

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

func (app App) ElasticBulkWrite(collections []CollectionsObject) {
	var (
		_ = fmt.Print
	)

	flag.Parse()

	log.SetFlags(0)

	type bulkResponse struct {
		Errors bool `json:"errors"`
		Items  []struct {
			Index struct {
				ID     string `json:"_id"`
				Result string `json:"result"`
				Status int    `json:"status"`
				Error  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
					Cause  struct {
						Type   string `json:"type"`
						Reason string `json:"reason"`
					} `json:"caused_by"`
				} `json:"error"`
			} `json:"index"`
		} `json:"items"`
	}

	var (
		buf bytes.Buffer
		res *esapi.Response
		raw map[string]interface{}
		blk *bulkResponse

		indexName = "arts"

		numItems   int
		numErrors  int
		numIndexed int
	)

	start := time.Now().UTC()

	// Loop over the collection
	for _, c := range collections {
		numItems++

		// Prepare the metadata payload
		meta := []byte(
			fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, c.ID, "\n"),
		)
		// Prepare the data payload: encode article to JSON
		data, err := json.Marshal(c)
		if err != nil {
			log.Println("Cannot encode article %d: %s", c.ID, err)

			return
		}

		// Append newline to the data payload
		data = append(data, "\n"...)

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	res, err := app.elasticClient.Bulk(
		bytes.NewReader(buf.Bytes()),
		app.elasticClient.Bulk.WithIndex(indexName),
		app.elasticClient.Bulk.WithDocumentType("art"),
	)

	if err != nil {
		log.Printf("Failure indexing batch %d: %s", err)

		return
	}

	if res.IsError() {
		numErrors += numItems
		if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
			log.Printf("Failure to to parse response body: %s", err)
		} else {
			log.Printf("  Error: [%d] %s: %s",
				res.StatusCode,
				raw["error"].(map[string]interface{})["type"],
				raw["error"].(map[string]interface{})["reason"],
			)
		}

		// A successful response might still contain errors for particular documents...
	} else {
		if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
			log.Fatalf("Failure to to parse response body: %s", err)
		} else {
			for _, d := range blk.Items {
				// ... so for any HTTP status above 201 ...
				if d.Index.Status > 201 {
					// ... increment the error counter ...
					numErrors++

					// ... and print the response status and error information ...
					log.Printf("  Error: [%d]: %s: %s: %s: %s",
						d.Index.Status,
						d.Index.Error.Type,
						d.Index.Error.Reason,
						d.Index.Error.Cause.Type,
						d.Index.Error.Cause.Reason,
					)
				} else {
					// ... otherwise increase the success counter.
					numIndexed++
				}
			}
		}
	}

	// Close the response body, to prevent reaching the limit for goroutines or file handles
	res.Body.Close()

	// Reset the buffer and items counter
	// buf.Reset()

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	log.Println(strings.Repeat("=", 80))

	dur := time.Since(start)

	if numErrors > 0 {
		log.Fatalf(
			"Indexed [%d] documents with [%d] errors in %s (%.0f docs/sec)",
			numIndexed,
			numErrors,
			dur.Truncate(time.Millisecond),
			1000.0/float64(dur/time.Millisecond)*float64(numIndexed),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%d] documents in %s (%.0f docs/sec)",
			numIndexed,
			dur.Truncate(time.Millisecond),
			1000.0/float64(dur/time.Millisecond)*float64(numIndexed),
		)
	}

	return
}
