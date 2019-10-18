package reverse_proxy

import (
	"net/http"

	"github.com/go-redis/redis/v7"
)

type App struct {
	baseURL       string
	harvardAPIKey string
	redisClient   *redis.Client
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
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type RecordsPayload struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Records interface{} `json:"records"`
}

type RecordPayload struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Record  interface{} `json:"record"`
}

type CollectionsResponse struct {
	Records []interface{} `json:"records"`
}

type PublicationsPayload struct {
	Error        bool          `json:"error"`
	Message      string        `json:"message"`
	Publications []interface{} `json:"publications"`
}

type PlaceIdPayload struct {
	Error   bool      `json:"error"`
	Message string    `json:"message"`
	Places  []PlaceID `json:"places"`
}

type PlacesPayload struct {
	Error   bool    `json:"error"`
	Message string  `json:"message"`
	Places  []Place `json:"places"`
}

type PlaceID struct {
	ParentPlaceId int64  `json:"parentPlaceID"`
	PathForward   string `json:"pathForward"`
}

type Place struct {
	Objectcount   int    `json:"objectcount"`
	ID            int    `json:"id"`
	LastUpdate    string `json:"lastupdate"`
	HasChildren   int    `json:"haschildren"`
	Level         int    `json:"level"`
	PlaceID       int    `json:"placeid"`
	PathForward   string `json:"pathforward"`
	ParentPlaceID int    `json:"parentplaceid"`
	Name          string `json:"name"`
	TgnID         int    `json:"tgn_id,omitempty"`
	Geo           struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"geo,omitempty"`
}
