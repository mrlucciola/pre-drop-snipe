package main

import "fmt"

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

	// retrieve tokens from server
	// tokens := getTokens(tokenSlug, int(10))
	tokens := getTokensConcurrently(tokenSlug, int(10))

	// create probability map
	probMap := buildTraitProbabilityMap(tokens, len(tokens))

	tokenRarityArr := calculateTokensRarity(tokens, probMap)

	for _, token := range tokens {
		token.lookupRarityRank(tokenRarityArr)
	}

	// Display the top five
	fmt.Println(sortRarityArr(tokenRarityArr)[:5])
}
