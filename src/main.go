package main

import (
	"fmt"
)

// https://opensea.io/collection/we-asuki
// var collectionSlug = "we-asuki"
var tokenSlug = "azuki1"

const ColorYellow = "\033[33m"
const ColorGreen = "\033[32m"
const ColorRed = "\033[31m"
const ColorReset = "\033[0m"

func main() {
	// Check the collection size
	collectionSize := 10000
	freqMap := TraitFrequencyMap{
		categories: make(map[string]*TraitValueFreqMap),
	}
	tokenMap := TokensMap{v: make(map[int]Token)}

	// retrieve tokens from server
	useConcurrency := true
	if useConcurrency {
		getTokensConcurrently(tokenSlug, &freqMap, &tokenMap, collectionSize)
	} else {
		freqMap = getTokens(tokenSlug, &tokenMap)
	}
	for key, token := range tokenMap.v {
		fmt.Println(key, token)
	}

	var tokenRarityArr []TokenRarity
	// tokenMap := freqMap.tokens.v
	useRarityScore := true
	if useRarityScore {
		// TODO: move in to the concurrent logic
		rarityScoreMap := buildTraitScoreMap(&freqMap)
		tokenRarityArr = calculateTokensRarityScores(rarityScoreMap, &tokenMap)
	} else {
		probMap := buildTraitProbabilityMap(&freqMap)
		tokenRarityArr = calculateTokensRarity(probMap, &tokenMap)
	}

	// sort
	sortedArr := sortRarityArr(tokenRarityArr)

	// Display the top five
	for _, v := range sortedArr[:5] {
		fmt.Println(v.id, v.rarity.StringFixed(20))
	}
}
