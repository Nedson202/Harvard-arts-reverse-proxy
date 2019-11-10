package reverse_proxy

import (
	"net/http"
)

func (app App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.respondWithJSON(w, http.StatusOK,
		RootPayload{
			Error:   false,
			Payload: "Harvard Art Museum Reverse Proxy API running",
		},
	)
}
