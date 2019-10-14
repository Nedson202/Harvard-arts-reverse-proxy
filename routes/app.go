package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nedson202/harvard-arts-reverse-proxy/config"
)

var routes []Route

//NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	// userRoutes := GetUserRoutes()
	healthChecker := GetHealthChecker()

	// routes = append(routes, userRoutes...)
	routes = append(routes, healthChecker...)
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = config.Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
