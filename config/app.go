package config

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

// LogFatal to handle logging errors
func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Logger function
func Logger(inner http.Handler, name string) http.Handler {
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

// RespondWithError handler for sending errors over http
func RespondWithError(w http.ResponseWriter, code int, errorData interface{}) {
	RespondWithJSON(w, code, RootPayload{Error: true, Payload: errorData})
}

// RespondWithJSON handler for sending responses over http
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
