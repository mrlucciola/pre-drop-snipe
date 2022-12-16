package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

const baseUrlToken = "https://go-challenge.skip.money"

// Map of trait group to value
//
// Data structure containing all of the traits for a single token
//
//	ex. traitValueMap: {
//	 "Background": "Off White A",
//	 "Clothing":   "Pink Oversized Kimono",
//	 "Eyes":       "Striking",
//	 "Offhand":    "Monkey King Staff",
//	 "Type":       "Human",
//	 "Hair":       "Water",
//	 "Mouth":      "Frown",
//	}
type TraitValueMap map[string]string

// We do not put the rarity on this class because it would require
// updating all token structs everytime a new token is minted.
//
// Use lookup table to find rarity.
type Token struct {
	id     int
	traits TraitValueMap
}

// ## Look up the current token's rarity in the rarity array
func (thisToken Token) lookupRarity(tokenRarityArr []TokenRarity) decimal.Decimal {
	return tokenRarityArr[thisToken.id].rarity
}

// # Calculate token rank.
//
// Lower probability = lower index order = higher rank.
//
// Optimized using 2 pointers.
//
// > TODO: parallelize
func calcRankOpt(searchValue decimal.Decimal, tokenRarityArr []TokenRarity) int {
	rank := 1
	arrLen := len(tokenRarityArr)
	isOdd := arrLen%2 != 0
	arrEndIdx := int(math.Floor(float64(arrLen) / 2.))

	p2 := arrLen
	for p1 := 0; p1 < arrEndIdx; p1++ {
		p2--

		if tokenRarityArr[p1].rarity != searchValue && tokenRarityArr[p1].rarity.LessThan(searchValue) {
			rank++
		}
		if tokenRarityArr[p2].rarity != searchValue && tokenRarityArr[p2].rarity.LessThan(searchValue) {
			rank++
		}
	}

	// if arr length is not even, eval the middle arr item
	if isOdd {
		if tokenRarityArr[arrEndIdx+1].rarity != searchValue && tokenRarityArr[arrEndIdx+1].rarity.LessThan(searchValue) {
			rank++
		}
	}

	return rank
}

// ## Sort the array and search for the index (i.e. rank) using the rarity value.
//
// All algos in the go::sort package are `O(n log n)`
//
// With how the data stores are designed:
//   - Regardless of whether the array is presorted or not, there will always be a O(n) lookup.
func (thisToken Token) lookupRarityRank(tokenRarityArr []TokenRarity) int {
	// get the the rarity value
	searchValue := thisToken.lookupRarity(tokenRarityArr)
	rank := calcRankOpt(searchValue, tokenRarityArr)

	return rank
}

// ## Makes GET request to Skip's servers, retrieves asset
//
// # From the stub script
//
// After receiving value from request, update frequency map
func getToken(client *http.Client, collectionSlug string, tokenId int, freqMap TraitFrequencyMap) Token {
	// build the endpoint string
	url := fmt.Sprintf("%s/%s/%d.json", baseUrlToken, collectionSlug, tokenId)

	// send request
	res, err := client.Get(url)

	// for now we are handling errors by returning an empty struct
	if err != nil {
		return Token{}
	}
	defer res.Body.Close()

	// read buffer
	body, err := io.ReadAll(res.Body)

	// for now we are handling errors by returning an empty struct
	if err != nil {
		return Token{}
	}

	// init the trait map
	traits := make(TraitValueMap)

	// deserialize token's traits from the response body's byte arr into our map
	json.Unmarshal(body, &traits)
	token := Token{
		id:     tokenId,
		traits: traits,
	}

	// NEW: update the trait value map
	// iterate thru the traits, add occurences to the map
	for traitGroup, traitValue := range token.traits {
		if _, found := freqMap[traitGroup]; !found {
			freqMap[traitGroup] = TraitValueFreqMap{}
		}

		freqMap[traitGroup][traitValue] += 1
	}

	return token
}

// ## Fetch all tokens for a given collection, without concurrency
//
// 1. Get the amount of total available tokens (normally would be from OpenSea collection stats)
//
// 2. Iterate through this range to get the collection's tokens
func getTokens(collectionSlug string, tokenArr []Token) TraitFrequencyMap {
	tokenCt := len(tokenArr)
	client := &http.Client{}
	freqMap := make(TraitFrequencyMap)

	for tokenId := 0; tokenId < tokenCt; tokenId++ {
		// log the token
		logger.Println(string(COLOR_GREEN), fmt.Sprintf("Getting token %d", tokenId), string(COLOR_RESET))
		token := getToken(client, collectionSlug, tokenId, freqMap)
		// add to the array
		tokenArr[tokenId] = token
	}

	// TODO: move this to the http request logic
	traitOccurences := buildTraitFrequencyMap(tokenArr)

	return traitOccurences
}

// ## Fetch all tokens for a given collection, using concurrency
//
// 1. Get the amount of total available tokens (normally would be from OpenSea collection stats)
//
// 2. Iterate through this range to get the collection's tokens
func getTokensConcurrently(collectionSlug string, tokenArr []Token) TraitFrequencyMap {
	jobCt := len(tokenArr)
	// init frequency map
	freqMap := make(TraitFrequencyMap)
	// set up workers and job pool
	workerCt := 2000
	transport := &http.Transport{
		ResponseHeaderTimeout: time.Hour,
		MaxConnsPerHost:       99999,
		DisableKeepAlives:     true,
	}
	client := &http.Client{Transport: transport}

	// init job channel and wait group
	jobChannel := make(chan int)
	var waitGroup sync.WaitGroup

	for workerId := 0; workerId < workerCt; workerId++ {
		waitGroup.Add(1)

		go func(wid int, wg *sync.WaitGroup) {
			defer wg.Done()
			for jobId := range jobChannel {
				// assign incoming token data to the pre-allocated array
				tokenArr[jobId] = getToken(client, collectionSlug, jobId, freqMap)
			}

		}(workerId, &waitGroup)
	}

	// Assign jobs to channel - add jobs to pool
	for jobId := 0; jobId < jobCt; jobId++ {
		jobChannel <- jobId
	}

	close(jobChannel)
	waitGroup.Wait()
	fmt.Println("(done) \n1  ", tokenArr[:5], "\n2  ", tokenArr[len(tokenArr)/2:5+len(tokenArr)/2], "\n3  ", tokenArr[len(tokenArr)-5:])

	// TODO: move this to the http request logic
	traitOccurences := buildTraitFrequencyMap(tokenArr)

	return traitOccurences
}

// DEPRECATED
// get the rank - lower probability = higher rank
// unoptomized
// func calcRankBF(searchValue float64, tokenRarityArr []TokenRarity) int {
// 	rank := 1

// 	for idx := 0; idx < len(tokenRarityArr); idx++ {
// 		lookupValue := tokenRarityArr[idx].rarity
// 		// handle duplicates
// 		if lookupValue != searchValue && lookupValue < searchValue {
// 			rank++
// 		}

// 	}
// 	return rank
// }
