package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Query OpenSea's api and get information about the collection and sales
func getOsCollection(client *http.Client, collectionSlug string) CollectionPayload {
	// base url for the OpenSea API
	baseUrlOs := "https://api.opensea.io/api/v1/collection"

	// OpenSea API endpoint
	endpointOs := fmt.Sprintf("%v/%v", baseUrlOs, collectionSlug)

	// build and send the request
	res := sendGetRequest(client, endpointOs)

	// deserialize response and handle errors
	var collectionRes collectionResponse
	if err := json.Unmarshal(res, &collectionRes); err != nil {
		log.Fatal("Error deserializing data")
	}

	return CollectionPayload{
		Traits:      collectionRes.Collection.Traits,
		TotalSupply: collectionRes.Collection.Stats.TotalSupply,
		Count:       collectionRes.Collection.Stats.Count,
	}
}

// Reformatted Collection response body
type CollectionPayload struct {
	Traits      map[string]map[string]int
	TotalSupply float64
	Count       float64
}

// Raw collection response format
type collectionResponse struct {
	Collection struct {
		Traits map[string]map[string]int
		Stats  struct {
			TotalSupply float64
			Count       float64
		}
	}
}
