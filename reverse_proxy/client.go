package reverse_proxy

import (
	"net/http"
	"time"
)

// NewClient create client with timeout
func (app App) NewClient() (client *http.Client) {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}

	return
}
