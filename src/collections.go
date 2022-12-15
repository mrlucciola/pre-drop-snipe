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

	// Build and send the request
	logger.Println(string(COLOR_YELLOW), fmt.Sprintf("Getting collection `%s`", collectionSlug), string(COLOR_RESET))
	res := sendGetRequest(client, endpointOs)
	logger.Println(string(COLOR_GREEN), fmt.Sprintf("Success: retrieved collection `%s`", collectionSlug), string(COLOR_RESET))

	// Deserialize response and handle errors
	var collectionRes collectionResponse
	if err := json.Unmarshal(res, &collectionRes); err != nil {
		log.Fatal("Error deserializing data")
	}

	return CollectionPayload{
		Traits: collectionRes.Collection.Traits,
		Count:  collectionRes.Collection.Stats.Count,
	}
}

type CollectionTraits map[string]TraitValueMap

// Reformatted Collection response body for readability
type CollectionPayload struct {
	Traits CollectionTraits
	Count  float64
}

// Body of the OpenSea API response for `collection`.
//
// Only includes relevant properties.
type collectionResponse struct {
	Collection struct {
		Traits CollectionTraits
		Stats  struct {
			Count float64
		}
	}
}
