package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func httpClient() *http.Client {
	// set a 5 second timeout
	client := &http.Client{Timeout: 5 * time.Second}

	return client
}

func sendGetRequest(client *http.Client, endpoint string) []byte {

	// build request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatalf("Error building request. %+v", err)
	}

	// send request
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request. %+v", err)
	}

	// close connection
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error parsing response body. %+v", err)
	}

	return body
}
