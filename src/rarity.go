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
	for traitCategory, traitValue := range tokenToEval.traits {
		// look up value in the map
		prob = probabilityMap[traitCategory][traitValue].Mul(prob)
	}

	// add to datastore
	tokenRarityArr[tokenToEval.id] = TokenRarity{id: tokenToEval.id, rarity: prob}

	return prob
}

// Calculate the "Rarity Score" for an individual token.
//
// Call after any update to token datastore, or token datastore init
func calculateTokenRarityScore(tokenToEval Token, rarityScoreMap TraitScoreMap, tokenRarityArr []TokenRarity) decimal.Decimal {
	// initialize at 0
	rarityScore := decimal.NewFromInt(int64(0))

	// iterate thru each trait on the token
	// TODO: parallelize if possible
	for traitCategory, traitValue := range tokenToEval.traits {
		// look up value in the map
		rarityScore = rarityScoreMap[traitCategory][traitValue].Add(rarityScore)
	}

	// add to datastore
	tokenRarityArr[tokenToEval.id] = TokenRarity{id: tokenToEval.id, rarity: rarityScore}

	return rarityScore
}

// Build list of rarity/probability values for an array of tokens.
//
// Loop over calculate-token-rarity fxn.
//
// # Returns array of token rarities, in the order indexible by token.id
//
// Array is not sorted
func calculateTokensRarity(probabilityMap TraitProbabilityMap, tokensToEval *TokensMap) []TokenRarity {
	// temporary: filter empty tokens out of array
	filteredMap := make(map[string]Token)
	for _, token := range tokensToEval.v {
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
	for _, token := range tokensToEval.v {
		calculateTokenRarity(token, probabilityMap, tokenRarityArr)
	}

	return tokenRarityArr
}

// # Build list of rarity scores for an array of tokens.
//
// Loop over calculate-token-rarity fxn.
//
// # Returns array of token rarity scores, in the order indexible by token.id
//
// Array is not sorted
func calculateTokensRarityScores(rarityScoreMap TraitScoreMap, tokenMap *TokensMap) []TokenRarity {
	// init arr
	tokenRarityArr := make([]TokenRarity, len(tokenMap.v))
	// parallelize?
	// we don't need to get a read lock for this. no writes for the rest of the program.
	for _, token := range tokenMap.v {
		calculateTokenRarityScore(token, rarityScoreMap, tokenRarityArr)
	}

	return tokenRarityArr
}

// Standard sorting function
func sortRarityArr(tokenRarityArr []TokenRarity) []TokenRarity {
	sort.Slice(tokenRarityArr[:], func(i, j int) bool {
		return tokenRarityArr[i].rarity.GreaterThan(tokenRarityArr[j].rarity)
	})

	return tokenRarityArr
}
