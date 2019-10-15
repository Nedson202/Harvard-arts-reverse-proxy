package reverse_proxy

import "net/http"

type App struct {
	baseURL       string
	harvardAPIKey string
}

// Route defines a structure for routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes defines the list of routes of our API
type Routes []Route

// RootPayload structure for error responses
type RootPayload struct {
	Error   bool        `json:"error"`
	Payload interface{} `json:"payload"`
}

// DataPayload structure for error responses
type DataPayload struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CollectionsResponse struct {
	Records interface{} `json:"records"`
}
