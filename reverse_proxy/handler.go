package reverse_proxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (app App) retrieveDataFromHarvardAPI(url string) (body []byte) {
	client := app.newClient()

	combinedURL := fmt.Sprintf("%s%s", app.baseURL, url)

	req, err := http.NewRequest("GET", combinedURL, nil)
	if err != nil {
		log.Println(err)

		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)

		return
	}

	return body
}
