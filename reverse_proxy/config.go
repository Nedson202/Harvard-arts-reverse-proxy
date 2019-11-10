package reverse_proxy

import (
	"encoding/json"
	"log"
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
func (app App) Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			r.UserAgent(),
			time.Since(start),
		)
		inner.ServeHTTP(w, r)
	})
}

// RespondWithJSON handler for sending responses over http
func (app App) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError handler for sending errors over http
func (app App) RespondWithError(w http.ResponseWriter, code int, errorData interface{}) {
	app.respondWithJSON(w, code, RootPayload{Error: true, Payload: errorData})
}
