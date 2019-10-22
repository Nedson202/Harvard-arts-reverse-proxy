package reverse_proxy

import (
	"net/http"

	"github.com/gorilla/mux"
)

var routes []Route
var healthCheckRoute Route
var collectionsRoute []Route
var publicationsRoute []Route
var placesRoute []Route

func (app App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.RespondWithJSON(w, http.StatusOK,
		RootPayload{
			Error:   false,
			Payload: "Harvard Art Museum Reverse Proxy API running",
		},
	)
}

func (app App) getHealthChecker() Route {
	healthCheckRoute = Route{
		"Index",
		"GET",
		"/health",
		app.healthCheckHandler,
	}

	return healthCheckRoute
}

func (app App) getCollectionRoutes() []Route {
	collectionsRoute = append(collectionsRoute,
		Route{
			"Index",
			"GET",
			"/harvard-arts/object",
			app.GetCollections,
		},
		Route{
			"Index",
			"GET",
			"/harvard-arts/object/{objectId}",
			app.GetCollection,
		},
		Route{
			"Index",
			"GET",
			"/harvard-arts/search",
			app.SearchCollections,
		},
	)

	return collectionsRoute
}

func (app App) getPublicationRoutes() []Route {
	publicationsRoute = append(publicationsRoute,
		Route{
			"Index",
			"GET",
			"/harvard-arts/publications",
			app.GetPublications,
		},
	)

	return publicationsRoute
}

func (app App) getPlacesRoutes() []Route {
	placesRoute = append(placesRoute,
		Route{
			"Index",
			"GET",
			"/harvard-arts/places/id",
			app.GetPlaceIds,
		},
		Route{
			"Index",
			"GET",
			"/harvard-arts/places",
			app.GetPlaces,
		},
	)

	return placesRoute
}

//NewRouter configures a new router to the API
func (app App) NewRouter(router *mux.Router) *mux.Router {
	healthChecker := app.getHealthChecker()
	getCollectionRoutes := app.getCollectionRoutes()
	getPublicationRoutes := app.getPublicationRoutes()
	getPlacesRoutes := app.getPlacesRoutes()

	routes = append(routes, healthChecker)
	routes = append(routes, getCollectionRoutes...)
	routes = append(routes, getPublicationRoutes...)
	routes = append(routes, getPlacesRoutes...)
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
