package reverse_proxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

var routes []Route
var collectionsRoute []Route
var healthCheckRoute Route

func (app App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.RespondWithJSON(w, http.StatusOK,
		RootPayload{
			Error:   false,
			Payload: "Harvard Art Museum Reverse Proxy API running",
		},
	)
}

// GetHealthChecker .
func (app App) GetHealthChecker() Route {
	healthCheckRoute = Route{
		"Index",
		"GET",
		"/health",
		app.healthCheckHandler,
	}

	return healthCheckRoute
}

// GetHealthChecker .
func (app App) GetCollectionRoutes() []Route {
	collectionsRoute = append(collectionsRoute, Route{
		"Index",
		"GET",
		"/objects",
		app.GetCollections,
	})

	return collectionsRoute
}

//NewRouter configures a new router to the API
func (app App) NewRouter(router *mux.Router) *mux.Router {
	healthChecker := app.GetHealthChecker()
	GetCollectionRoutes := app.GetCollectionRoutes()

	routes = append(routes, healthChecker)
	routes = append(routes, GetCollectionRoutes...)
	combineRoutes := router.StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = app.Logger(handler, route.Name)
		combineRoutes.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
