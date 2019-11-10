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

func (app App) getHealthChecker() Route {
	healthCheckRoute = Route{
		"HealthCheck",
		"GET",
		"/health",
		app.healthCheckHandler,
	}

	return healthCheckRoute
}

func (app App) getCollectionRoutes() []Route {
	collectionsRoute = append(collectionsRoute,
		Route{
			"Objects",
			"GET",
			"/api/v1/objects",
			app.getCollections,
		},
		Route{
			"Object",
			"GET",
			"/api/v1/object/{objectId}",
			app.getCollection,
		},
		Route{
			"Search",
			"GET",
			"/api/v1/search",
			app.searchCollections,
		},
	)

	return collectionsRoute
}

func (app App) getPublicationRoutes() []Route {
	publicationsRoute = append(publicationsRoute,
		Route{
			"Publications",
			"GET",
			"/api/v1/publications",
			app.getPublications,
		},
	)

	return publicationsRoute
}

func (app App) getPlacesRoutes() []Route {
	placesRoute = append(placesRoute,
		Route{
			"PlaceIDs",
			"GET",
			"/api/v1/places/id",
			app.getPlaceIds,
		},
		Route{
			"Places",
			"GET",
			"/api/v1/places",
			app.getPlaces,
		},
	)

	return placesRoute
}

//NewRouter configures a new router to the API
func (app App) newRouter(router *mux.Router) *mux.Router {
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
		combineRoutes.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.Use(app.Logger)

	return router
}
