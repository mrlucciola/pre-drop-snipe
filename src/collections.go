package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// DEPRECATED
// Query OpenSea's api and get information about the collection and sales
func getOsCollection(collectionSlug string) CollectionPayload {
	client := httpClient()
	// base url for the OpenSea API
	baseUrlOs := "https://api.opensea.io/api/v1/collection"

	// OpenSea API endpoint
	endpointOs := fmt.Sprintf("%v/%v", baseUrlOs, collectionSlug)

	// Build and send the request
	logger.Println(string(ColorYellow), fmt.Sprintf("Getting collection `%s`", collectionSlug), string(ColorReset))
	res := sendGetRequest(client, endpointOs)
	logger.Println(string(ColorGreen), fmt.Sprintf("Success: retrieved collection `%s`", collectionSlug), string(ColorReset))

	// Deserialize response and handle errors
	var collectionRes CollectionResponse
	if err := json.Unmarshal(res, &collectionRes); err != nil {
		log.Fatal("Error deserializing data")
	}

	return CollectionPayload{
		// "copy lock value" warning is inconsequential here, since all ops are sequential.
		Traits: collectionRes.Collection.Traits,
		Count:  collectionRes.Collection.Stats.Count,
	}
}

// Reformatted Collection response body for readability
type CollectionPayload struct {
	Traits TraitFrequencyMap
	Count  float64
}

// Body of the OpenSea API response for `collection`.
//
// Only includes relevant properties.
type CollectionResponse struct {
	Collection struct {
		Traits TraitFrequencyMap
		Stats  struct {
			Count float64
		}
	}
}
