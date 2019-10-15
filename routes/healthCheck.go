package routes

import (
	"net/http"

	"github.com/nedson202/harvard-arts-reverse-proxy/config"
)

var healthCheckRoute []Route

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	config.RespondWithJSON(w, 200,
		config.RootPayload{
			Error:   false,
			Payload: "Harvard Art Museum Reverse Proxy API running",
		},
	)
}

// GetHealthChecker .
func GetHealthChecker() []Route {
	healthCheckRoute = append(healthCheckRoute,
		Route{
			"Index",
			"GET",
			"/health",
			healthCheckHandler,
		},
	)

	return healthCheckRoute
}
