package main

import (
	"fmt"
	"sort"
)

type TokenRarity struct {
	rarity float64
	id     int
}

// Calculate the probability for an individual token.
//
// Call after any update to token datastore, or token datastore init
func calculateTokenRarity(tokenToEval Token, probabilityMap TraitProbabilityMap, tokenRarityArr []TokenRarity) float64 {
	// initialize at 1
	prob := 1.

	// iterate thru each trait on the token
	// TODO: parallelize if possible
	for traitGroup, traitValue := range tokenToEval.traits {
		// look up value in the map
		prob *= probabilityMap[traitGroup][traitValue]
	}

	// add to datastore
	tokenRarityArr[tokenToEval.id] = TokenRarity{id: tokenToEval.id, rarity: roundFloat(prob, 15)}

	return prob
}

// Build list of rarity/probability values for an array of tokens.
//
// Loop over calculate-token-rarity fxn.
//
// # Returns array of token rarities, in the order indexible by token.id
//
// Array is not sorted
func calculateTokensRarity(tokensToEval []Token, probabilityMap TraitProbabilityMap) []TokenRarity {
	// temporary: filter empty tokens out of array
	filteredMap := make(map[string]Token)
	for idx, token := range tokensToEval {
		fmt.Println("token pre", idx, token)
		// if token's trait struct is empty, ignore
		if len(token.traits) > 0 {
			idStr := fmt.Sprintf("%d", token.id)
			filteredMap[idStr] = token
		}
	}
	filteredTokenArr := make([]Token, len(filteredMap))
	for idx, token := range filteredMap {
		fmt.Println("token post", idx, token)
		filteredTokenArr[token.id] = token
	}

	// init arr
	tokenRarityArr := make([]TokenRarity, len(filteredMap))
	// parallelize?
	for _, token := range tokensToEval {
		calculateTokenRarity(token, probabilityMap, tokenRarityArr)
	}

	return tokenRarityArr
}

// Standard sorting function
func sortRarityArr(tokenRarityArr []TokenRarity) []TokenRarity {
	sort.Slice(tokenRarityArr[:], func(i, j int) bool {
		return tokenRarityArr[i].rarity < tokenRarityArr[j].rarity
	})

	return tokenRarityArr
}
