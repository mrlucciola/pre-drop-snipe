package main

import (
	"fmt"
)

// https://opensea.io/collection/azuki1
// var collectionSlug = "we-asuki"
var tokenSlug = "azuki1"

const COLOR_YELLOW = "\033[33m"
const COLOR_GREEN = "\033[32m"
const COLOR_RED = "\033[31m"
const COLOR_RESET = "\033[0m"

func main() {
	// client := httpClient()
	// res := getOsCollection(client, collectionSlug)

	// preallocate an array - as long as we know up front how many tokens we need to call
	// we can store using their id as this array's index
	tokens := make([]Token, 10000)

	// retrieve tokens from server
	// tokens := getTokens(tokenSlug, int(10))
	getTokensConcurrently(tokenSlug, tokens)

	// create probability map
	// TODO: move in to the concurrent logic
	probMap := buildTraitProbabilityMap(tokens, len(tokens))

	tokenRarityArr := calculateTokensRarity(tokens, probMap)

	// sort
	sortedArr := sortRarityArr(tokenRarityArr)

	// Display the top five
	fmt.Println(sortedArr[:5])
}
