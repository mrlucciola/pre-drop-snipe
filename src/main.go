package main

// https://opensea.io/collection/azuki1
var collectionSlug = "we-asuki"
var tokenSlug = "azuki1"

const COLOR_YELLOW = "\033[33m"
const COLOR_GREEN = "\033[32m"
const COLOR_RED = "\033[31m"
const COLOR_RESET = "\033[0m"

func main() {
	// client := httpClient()
	// res := getOsCollection(client, collectionSlug)

	tokens := getTokens(tokenSlug, int(10))
	tokenRarityArr := make([]float64, len(tokens))

	probMap := getAllTraitStatsSkip(tokens, len(tokens))

	for _, token := range tokens {
		calculateRarity(*token, probMap, tokenRarityArr)
	}

	for _, token := range tokens {
		token.lookupRarityRank(tokenRarityArr)
	}
}
