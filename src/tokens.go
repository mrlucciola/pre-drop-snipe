package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseUrlToken = "https://go-challenge.skip.money"

// Makes GET request to Skip's servers, retrieves asset
// From the stub script
func getToken(collectionSlug string, tokenId int) *Token {
	// build the endpoint string
	url := fmt.Sprintf("%s/%s/%d.json", baseUrlToken, collectionSlug, tokenId)

	// send request
	res, err := http.Get(url)

	// handle res
	if err != nil {
		logger.Println(string(COLOR_RED), fmt.Sprintf("Error getting token %d :", tokenId), err, string(COLOR_RESET))
		return &Token{}
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Println(string(COLOR_RED), fmt.Sprintf("Error reading response for token %d :", tokenId), err, string(COLOR_RESET))
		return &Token{}
	}

	traits := make(map[string]string)
	json.Unmarshal(body, &traits)

	for traitType, traitTypeOccurrence := range traits {
		fmt.Println("traitType:", traitType, "occurrence:", traitTypeOccurrence)
	}

	return &Token{
		id:     tokenId,
		traits: traits,
	}
}

type Token struct {
	id     int
	traits map[string]string
}

type Tokens []*Token
