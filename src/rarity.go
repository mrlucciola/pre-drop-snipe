package main

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
)

type TokenRarity struct {
	rarity decimal.Decimal
	id     int
}

// Calculate the probability for an individual token.
//
// Call after any update to token datastore, or token datastore init
func calculateTokenRarity(tokenToEval Token, probabilityMap TraitProbabilityMap, tokenRarityArr []TokenRarity) decimal.Decimal {
	// initialize at 1
	prob := decimal.New(1, initPrecision)

	// iterate thru each trait on the token
	// TODO: parallelize if possible
	for traitGroup, traitValue := range tokenToEval.traits {
		// look up value in the map
		prob = probabilityMap[traitGroup][traitValue].Mul(prob)
	}

	// add to datastore
	tokenRarityArr[tokenToEval.id] = TokenRarity{id: tokenToEval.id, rarity: prob}

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
	for _, token := range tokensToEval {
		// if token's trait struct is empty, ignore
		if len(token.traits) > 0 {
			idStr := fmt.Sprintf("%d", token.id)
			filteredMap[idStr] = token
		}
	}
	filteredTokenArr := make([]Token, len(filteredMap))
	for _, token := range filteredMap {
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
		return tokenRarityArr[i].rarity.LessThan(tokenRarityArr[j].rarity)
	})

	return tokenRarityArr
}
