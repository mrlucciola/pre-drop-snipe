package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
)

const baseUrlToken = "https://go-challenge.skip.money"

// We do not put the rarity on this class because it updates everytime a new token is minted.
//
// Use lookup table to find rarity.
type Token struct {
	id     int
	traits map[string]string
}

func (thisToken Token) lookupRarity(tokenRarityArr []TokenRarity) float64 {
	return tokenRarityArr[thisToken.id].rarity
}

// get the rank - lower probability = higher rank
// unoptomized
func calcRankBF(searchValue float64, tokenRarityArr []TokenRarity) int {
	rank := 1

	for idx := 0; idx < len(tokenRarityArr); idx++ {
		lookupValue := tokenRarityArr[idx].rarity
		// handle duplicates
		if lookupValue != searchValue && lookupValue < searchValue {
			rank++
		}

	}
	return rank
}

// get the rank - lower probability = lower index order (i.e. higher rank)
// optimized using 2 pointers
// TODO: parallelize
func calcRankOpt(searchValue float64, tokenRarityArr []TokenRarity) int {
	rank := 1
	arrLen := len(tokenRarityArr)
	isOdd := arrLen%2 != 0
	arrEndIdx := int(math.Floor(float64(arrLen) / 2.))

	p2 := arrLen
	for p1 := 0; p1 < arrEndIdx; p1++ {
		p2--

		if tokenRarityArr[p1].rarity != searchValue && tokenRarityArr[p1].rarity < searchValue {
			rank++
		}
		if tokenRarityArr[p2].rarity != searchValue && tokenRarityArr[p2].rarity < searchValue {
			rank++
		}
	}

	// if arr length is not even, eval the middle arr item
	if isOdd {
		if tokenRarityArr[arrEndIdx+1].rarity != searchValue && tokenRarityArr[arrEndIdx+1].rarity < searchValue {
			rank++
		}
	}

	return rank
}

// Sort the array and search for the index (i.e. rank) using the rarity value.
//
// All algos in the go::sort package are `O(n log n)`
//
// With how the data stores are designed:
//   - Regardless of whether the array is presorted or not, there will always be a O(n) lookup.
func (thisToken Token) lookupRarityRank(tokenRarityArr []TokenRarity) int {
	// get the the rarity value
	searchValue := roundFloat(thisToken.lookupRarity(tokenRarityArr), 15)

	rank := calcRankBF(searchValue, tokenRarityArr)
	rank2 := calcRankOpt(searchValue, tokenRarityArr)

	fmt.Println("value:", searchValue, "    rank 1:", rank, "rank2:", rank2)
	return rank
}

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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Println(string(COLOR_RED), fmt.Sprintf("Error reading response for token %d :", tokenId), err, string(COLOR_RESET))
		return &Token{}
	}

	traits := make(map[string]string)
	json.Unmarshal(body, &traits)

	return &Token{
		id:     tokenId,
		traits: traits,
	}
}

// Retrieve all tokens for a given collection
//
// 1. Get the amount of total available tokens (normally would be from OpenSea collection stats)
//
// 2. Iterate through this range to get the collection's tokens
func getTokens(collectionSlug string, tokenCt int) []*Token {
	tokenArr := make([]*Token, tokenCt)

	for tokenId := 0; tokenId < tokenCt; tokenId++ {
		// log the token
		logger.Println(string(COLOR_GREEN), fmt.Sprintf("Getting token %d", tokenId), string(COLOR_RESET))
		token := getToken(collectionSlug, tokenId)
		// add to the array
		tokenArr[tokenId] = token
	}

	return tokenArr
}
func getTokensConcurrently(collectionSlug string, tokenCt int) []*Token {
	wg := sync.WaitGroup{}

	tokenArr := make([]*Token, tokenCt)

	for tokenId := 0; tokenId < tokenCt; tokenId++ {
		wg.Add(1)

		go func(tid int) {
			// log the token
			logger.Println(string(COLOR_GREEN), fmt.Sprintf("Getting token %d", tid), string(COLOR_RESET))
			token := getToken(collectionSlug, tid)
			// add to the array
			tokenArr[tid] = token

		}(tokenId)
	}

	return tokenArr
}
