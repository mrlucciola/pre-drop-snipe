package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

const baseUrlToken = "https://go-challenge.skip.money"

// # Data store for all tokens from API
// We are more concerned with data consistency so we are trading off speed
// for a map that only contains valid entries (no blanks) and improved readability via reduced code
type TokensMap struct {
	mu sync.RWMutex
	v  map[int]Token
}

// Map of trait category to value
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
type TokenTraitMap map[string]string

// # Holds information about a single token of a single collection relevant to ranking
//
// - Token ID
// - List of traits assigned to the token
//
// All tokens of a collection have the same categories of traits,
// but may have a different "flavor" (or "type") for each category of trait.
//
// We do not put the rarity on this class because it would require
// updating all token structs everytime a new token is minted.
//
// Use lookup table to find rarity.
type Token struct {
	id     int
	traits TokenTraitMap
}

// # Look up the current token's rarity in the rarity array
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

// # Sort the array and search for the index (i.e. rank) using the rarity value.
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

// # Makes GET request to Skip's servers, retrieves asset
//
// After receiving value from request, update frequency map
// - From the stub script.
func getToken(client *http.Client, collectionSlug string, tokenId int, freqMap *TraitFrequencyMap, tokensMap *TokensMap) {
	// build the endpoint string
	url := fmt.Sprintf("%s/%s/%d.json", baseUrlToken, collectionSlug, tokenId)

	// send request
	res, err := client.Get(url)

	// for now we are handling errors by returning an empty struct
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching token: %v", tokenId)
	}
	defer res.Body.Close()

	// read buffer
	body, err := io.ReadAll(res.Body)

	// for now we are handling errors by returning an empty struct
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error handling token: %v", tokenId)
	}

	// init the trait map
	traits := make(TokenTraitMap)

	// deserialize token's traits from the response body's byte arr into our map
	json.Unmarshal(body, &traits)
	token := Token{
		id:     tokenId,
		traits: traits,
	}

	// NEW: update the trait value map
	// iterate thru the traits, add occurences to the map
	for traitCategory, traitValue := range token.traits {
		// add traits to the map
		if _, found := freqMap.categories[traitCategory]; !found {
			// data race: if two threads wanted to create a category and didnt know about each other, one would overwrite the other's data
			// handle concurrent access
			freqMap.mu.Lock()
			freqMap.categories[traitCategory] = &TraitValueFreqMap{values: make(map[string]int)}
			freqMap.mu.Unlock()
		}
		// handle concurrent access
		freqMap.categories[traitCategory].mu.Lock()
		freqMap.categories[traitCategory].values[traitValue] += 1
		freqMap.categories[traitCategory].mu.Unlock()
	}

	// NEW: add to token map
	tokensMap.mu.Lock()
	tokensMap.v[tokenId] = token
	tokensMap.mu.Unlock()
}

// # Fetch all tokens for a given collection, using concurrency
//
// 1. Get the amount of total available tokens (normally would be from OpenSea collection stats)
//
// 2. Iterate through this range to get the collection's tokens
func getTokensConcurrently(collectionSlug string, freqMap *TraitFrequencyMap, tokensMap *TokensMap, collectionSize int) {
	// job count = token count
	tokenCt := collectionSize

	// set up workers and job pool
	workerCt := 2000
	transport := &http.Transport{
		ResponseHeaderTimeout: time.Hour,
		MaxConnsPerHost:       99999,
		DisableKeepAlives:     true,
	}
	client := &http.Client{Transport: transport}

	// init job channel and wait category
	jobChannel := make(chan int)
	var waitGroup sync.WaitGroup

	for workerId := 0; workerId < workerCt; workerId++ {
		waitGroup.Add(1)

		go func(wid int, wg *sync.WaitGroup) {
			defer wg.Done()
			for tokenId := range jobChannel {
				getToken(client, collectionSlug, tokenId, freqMap, tokensMap)
			}

		}(workerId, &waitGroup)
	}

	// Assign jobs to channel - add jobs to pool
	for jobId := 0; jobId < tokenCt; jobId++ {
		jobChannel <- jobId
	}

	close(jobChannel)
	waitGroup.Wait()
}

// DEPRECATED
// # Fetch all tokens for a given collection, without concurrency
//
// 1. Get the amount of total available tokens (normally would be from OpenSea collection stats)
//
// 2. Iterate through this range to get the collection's tokens
func getTokens(collectionSlug string, tokenMap *TokensMap) TraitFrequencyMap {
	tokenCt := len(tokenMap.v)
	client := &http.Client{}
	// var rwMu sync.RWMutex
	freqMap := TraitFrequencyMap{categories: make(map[string]*TraitValueFreqMap)}

	for tokenId := 0; tokenId < tokenCt; tokenId++ {
		// log the token
		logger.Println(string(ColorGreen), fmt.Sprintf("Getting token %d", tokenId), string(ColorReset))
		getToken(client, collectionSlug, tokenId, &freqMap, tokenMap)
	}

	// TODO: move this to the http request logic
	traitOccurences := buildTraitFrequencyMap(tokenMap)

	// "copy lock value" warning is inconsequential here, since all ops are sequential.
	return traitOccurences
}
