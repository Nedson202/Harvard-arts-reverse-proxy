package reverse_proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// LogFatal to handle logging errors
func (app App) LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Logger function
func (app App) Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

// RespondWithJSON handler for sending responses over http
func (app App) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError handler for sending errors over http
func (app App) RespondWithError(w http.ResponseWriter, code int, errorData interface{}) {
	app.RespondWithJSON(w, code, RootPayload{Error: true, Payload: errorData})
}

func (app App) RetrieveDataFromHarvardAPI(url string) (body []byte) {
	client := app.NewClient()

	combinedURL := fmt.Sprintf("%s%s", app.baseURL, url)

	req, err := http.NewRequest("GET", combinedURL, nil)
	if err != nil {
		log.Println(err)

		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return
	}

	return body
}

func (app App) RandomizeData(data []interface{}) (shuffled []interface{}) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for n := len(data); n > 0; n-- {
		randIndex := r.Intn(n)

		data[n-1], data[randIndex] = data[randIndex], data[n-1]
	}

	return data
}
