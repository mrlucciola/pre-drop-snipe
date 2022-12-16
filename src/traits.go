package main

import (
	"sync"

	"github.com/shopspring/decimal"
)

// # Frequency map for trait values of a single category
//
// Shows the occurrences of each trait.
//
// Used to create the trait value probability map.
//
//	ex. Mouth: {
//	  "m2 bored":          991,
//	  "m1 bored":          3313,
//	  "m1 bored unshaven": 2257,
//	  "m1 dumbfounded":    750,
//	}
type TraitValueFreqMap struct {
	mu     sync.RWMutex
	values map[string]int
}

// # A collection of trait value frequencies, ordered by category
//
// Trait-category maps contain maps for trait-values map[traitValueStr]int
// Trait-value maps contain each value's frequency - the amount of appearances in a list of tokens.
type TraitFrequencyMap struct {
	mu         sync.RWMutex
	categories map[string]*TraitValueFreqMap
}

// # Holds the rarity scores for each individual trait value.
//
// 1. Multiply the trait frequency by # of trait values within category
// 2. Take the inverse of this value
//
//	given category Mouth: {
//	  "m2 bored":          0.0002522704339, = 1 / (3,964    = 991 * 4)
//	  "m1 bored":          0.0000754603078, = 1 / (13,252   = 3313 * 4)
//	  "m1 bored unshaven": 0.0001107665042, = 1 / (9,028    = 2257 * 4)
//	  "m1 dumbfounded":    0.0003333333333, = 1 / (3,000    = 750 * 4)
//	}
type TraitValueScoreMap map[string]decimal.Decimal

// # Trait rarity scores, by category.
type TraitScoreMap map[string]TraitValueScoreMap

// # Probability map for trait values of a single category.
//
// Trait value probabilities are calculated from the frequency mapping.
//
//	Mouth: {
//	  "m2 bored":          0.1355 ( 991 / 7311),
//	  "m1 bored":          0.4531 (3313 / 7311),
//	  "m1 bored unshaven": 0.3087 (2257 / 7311),
//	  "m1 dumbfounded":    0.1025 ( 750 / 7311),
//	}
type TraitValueProbMap map[string]decimal.Decimal

// # Probabilities for all traits.
type TraitProbabilityMap map[string]TraitValueProbMap

// # Parent struct. Holds all stats about tokens.
type TraitStatMap struct {
	freq  TraitFrequencyMap
	prob  TraitProbabilityMap
	score TraitScoreMap
}

const initPrecision = 0

/*
Build a mapping of all probabilities for all possible traits assignable to a token.

The map is a nested data structure that holds
  - "Trait Categories" at the top level
  - "Trait Values" within trait categories

Iterate through the list of tokens pulled from the Skip Protocol API and
build a frequency map for each trait as they appear.
  - Input is `Token` array - from Skip Protocol API
*/
func buildTraitProbabilityMap(traitOccurences *TraitFrequencyMap) TraitProbabilityMap {

	// TODO: parallelize if possible
	traitProbabilities := make(TraitProbabilityMap)
	for traitCategory, traitValueFreqMap := range traitOccurences.categories {
		// sum all value counts for the category
		categorySum := 0
		// TODO: parallelize if possible
		for _, valueCt := range traitValueFreqMap.values {
			categorySum += valueCt
		}
		// handle empty properties
		if _, found := traitProbabilities[traitCategory]; !found {
			traitProbabilities[traitCategory] = make(TraitValueProbMap)
		}

		// iter thru each value, divide to get prob for each value, assign to map
		// TODO: parallelize if possible
		for traitValue, valueCt := range traitValueFreqMap.values {
			val := decimal.NewFromInt(int64(valueCt)).Div(decimal.NewFromInt(int64(categorySum)))
			traitProbabilities[traitCategory][traitValue] = val
		}
	}

	return traitProbabilities
}

/*
# Build a mapping of all rarity scores for all possible traits assignable to a token.

The map is a nested data structure that holds
  - "Trait Categoriess" at the top level
  - "Trait Values" within trait categories

Iterate through the list of tokens pulled from the Skip Protocol API and
build a frequency map for each trait as they appear.
  - Input is `Token` array - from Skip Protocol API

Relies on the trait frequency map to be calculated
*/
func buildTraitScoreMap(freqMap *TraitFrequencyMap) TraitScoreMap {

	// TODO: parallelize if possible
	traitScores := make(TraitScoreMap)
	for traitCategory, traitValueFreqMap := range freqMap.categories {
		// get the number of unique trait values
		uniqueTraitCt := len(traitValueFreqMap.values)
		// handle empty properties
		if _, found := traitScores[traitCategory]; !found {
			traitScores[traitCategory] = make(TraitValueScoreMap)
		}
		// TODO: parallelize if possible
		// multiply the amount of times this value occurs, by the amount of unique traits (values) within the category
		for traitValue, valueFreqCt := range traitValueFreqMap.values {
			traitValueRarityScore := decimal.New(1, 0).Div(decimal.NewFromInt(int64(valueFreqCt)).Mul(decimal.NewFromInt(int64(uniqueTraitCt))))
			traitScores[traitCategory][traitValue] = traitValueRarityScore
		}

	}

	return traitScores
}

// DEPRECATED
//
// # Build the trait frequency map
// TODO: parallelize if possible
func buildTraitFrequencyMap(tokenMap *TokensMap) TraitFrequencyMap {
	freqMap := TraitFrequencyMap{categories: make(map[string]*TraitValueFreqMap)}

	for _, token := range tokenMap.v {

		// iterate thru the traits, add occurences to the map
		for traitCategory, traitValue := range token.traits {
			if _, found := freqMap.categories[traitCategory]; !found {
				freqMap.categories[traitCategory] = &TraitValueFreqMap{values: make(map[string]int)}
			}

			// handle concurrent access
			freqMap.categories[traitCategory].mu.Lock()
			freqMap.categories[traitCategory].values[traitValue] += 1
			freqMap.categories[traitCategory].mu.Unlock()
		}
	}
	return freqMap
}

// DEPRECATED
//
// # Calculate trait value probabilities, by category
//
// Input the frequency for the trait values of a single category & return its probability mapping
// func calcTraitValueProbsByCategory(traitCategory TraitValueFreqMap, activeTokenCount int) TraitValueProbMap {
// 	// get sum of all occurences
// 	var traitOccurrenceSum int

// 	// TODO: parallelize
// 	// m1 dumbfounded + m2 bored + m1 bored + m1 bored unshaven...
// 	for _, traitValueOccurrence := range traitCategory {
// 		traitOccurrenceSum += traitValueOccurrence
// 	}

// 	traitCategoryProb := TraitValueProbMap{}
// 	// now that we have the sum, apply probability to each trait value
// 	for traitValue, traitValueOccurrence := range traitCategory {
// 		prob := float64(traitValueOccurrence) / float64(traitOccurrenceSum)
// 		traitCategoryProb[traitValue] = prob
// 	}

// 	return traitCategoryProb
// }

// DEPRECATED: update test
// Calculate the probabilities for all traits for the token collection in its current state.
//
// Creates a probability mapping of all traits, nested within their respective categories.
//
// Sources info from Skip Protocol's database.
//
// Note: OpenSea has a lack of consistency for contract versions with Akuri.
// func calcAllTraitValueProbs(traits TraitFrequencyMap, activeTokenCount int) TraitProbabilityMap {
// 	probMap := TraitProbabilityMap{}

// 	// TODO: parallelize
// 	for categoryName, category := range traits {
// 		probMap[categoryName] = calcTraitValueProbsByCategory(category, int(activeTokenCount))
// 	}

// 	return probMap
// }
