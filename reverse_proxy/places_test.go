package reverse_proxy

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

const (
	okResponse = `{
		"users": [
			{"id": 1, "name": "Roman"},
			{"id": 2, "name": "Dmitry"}
		]	
	}`
)

func TestGetPlaces(t *testing.T) {
	router := mux.NewRouter()
	app, _ := New("test_redis_host", router, "http://localhost", "test_api_key", "http://127.0.0.1:9200")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(okResponse))
	})

	_, teardown := app.testingHTTPClient(h)
	defer teardown()
}
