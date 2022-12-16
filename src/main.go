package main

import "fmt"

// https://opensea.io/collection/we-asuki
// var collectionSlug = "we-asuki"
var tokenSlug = "azuki1"

const ColorYellow = "\033[33m"
const ColorGreen = "\033[32m"
const ColorRed = "\033[31m"
const ColorReset = "\033[0m"

func main() {
	// preallocate an array - as long as we know up front how many tokens we need to call
	// we can store using their id as this array's index
	tokens := make([]Token, 10000)
	freqMap := TraitFrequencyMap{groups: make(map[string]*TraitValueFreqMap)}

	// retrieve tokens from server
	useConcurrency := true
	if useConcurrency {
		getTokensConcurrently(tokenSlug, tokens, &freqMap)
	} else {
		freqMap = getTokens(tokenSlug, tokens)
	}

	var tokenRarityArr []TokenRarity
	useRarityScore := true
	if useRarityScore {
		// TODO: move in to the concurrent logic
		rarityScoreMap := buildTraitScoreMap(tokens, &freqMap)
		tokenRarityArr = calculateTokensRarityScores(tokens, rarityScoreMap)
	} else {
		probMap := buildTraitProbabilityMap(tokens, &freqMap)
		tokenRarityArr = calculateTokensRarity(tokens, probMap)
	}

	// sort
	sortedArr := sortRarityArr(tokenRarityArr)

	// Display the top five
	for _, v := range sortedArr[:5] {
		fmt.Println(v.id, v.rarity.StringFixed(20))
	}
}
