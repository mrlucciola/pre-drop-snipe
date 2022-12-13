package main

import (
	"fmt"
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

// Query OpenSea's api and get information about the collection and sales
func getOsCollection(client *http.Client, collectionSlug string) []byte {
	// base url for the OpenSea API
	baseUrlOs := "https://api.opensea.io/api/v1/collection"

	// OpenSea API endpoint
	endpointOs := fmt.Sprintf("%v/%v", baseUrlOs, collectionSlug)

	// build and send the request
	resJson := sendGetRequest(client, endpointOs)
	return resJson
}
