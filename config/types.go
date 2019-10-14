package config

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
