package reverse_proxy

import (
	"net/http"
	"time"
)

func (app App) NewClient() (client *http.Client) {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}

	return
}
