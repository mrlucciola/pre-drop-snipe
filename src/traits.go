package main

import "github.com/shopspring/decimal"

// Frequency map for trait values of a single group
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
type TraitValueFreqMap map[string]int

// Calculated like this:
//
//	given category Mouth: {
//	  "m2 bored":          0.0002522704339, = 1 / (3,964    = 991 * 4)
//	  "m1 bored":          0.0000754603078, = 1 / (13,252   = 3313 * 4)
//	  "m1 bored unshaven": 0.0001107665042, = 1 / (9,028    = 2257 * 4)
//	  "m1 dumbfounded":    0.0003333333333, = 1 / (3,000    = 750 * 4)
//	}
type TraitValueScoreMap map[string]decimal.Decimal

// Probability map for trait values of a single group.
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

type TraitFrequencyMap map[string]TraitValueFreqMap

// Probabilities for all traits
type TraitProbabilityMap map[string]TraitValueProbMap

// Trait rarity scores, by category
type TraitScoreMap map[string]TraitValueScoreMap

const initPrecision = 0

// Build the trait frequency map
// TODO: parallelize if possible
func buildTraitFrequencyMap(tokenArr []Token) TraitFrequencyMap {
	traitOccurences := make(TraitFrequencyMap)

	for _, token := range tokenArr {

		// iterate thru the traits, add occurences to the map
		for traitGroup, traitValue := range token.traits {
			if _, found := traitOccurences[traitGroup]; !found {
				traitOccurences[traitGroup] = TraitValueFreqMap{}
			}

			traitOccurences[traitGroup][traitValue] += 1
		}
	}
	return traitOccurences
}

/*
Build a mapping of all probabilities for all possible traits assignable to a token.

The map is a nested data structure that holds
  - "Trait Groups" at the top level
  - "Trait Values" within trait groups

Iterate through the list of tokens pulled from the Skip Protocol API and
build a frequency map for each trait as they appear.
  - Input is `Token` array - from Skip Protocol API
*/
func buildTraitProbabilityMap(tokenArr []Token, traitOccurences TraitFrequencyMap) TraitProbabilityMap {

	// TODO: parallelize if possible
	traitProbabilities := make(TraitProbabilityMap)
	for traitGroup, traitValueFreqMap := range traitOccurences {
		// sum all value counts for the group
		groupSum := 0
		// TODO: parallelize if possible
		for _, valueCt := range traitValueFreqMap {
			groupSum += valueCt
		}
		// handle empty properties
		if _, found := traitProbabilities[traitGroup]; !found {
			traitProbabilities[traitGroup] = make(TraitValueProbMap)
		}

		// iter thru each value, divide to get prob for each value, assign to map
		// TODO: parallelize if possible
		for traitValue, valueCt := range traitValueFreqMap {
			val := decimal.NewFromInt(int64(valueCt)).Div(decimal.NewFromInt(int64(groupSum)))
			traitProbabilities[traitGroup][traitValue] = val
		}
	}

	return traitProbabilities
}
func buildTraitScoreMap(tokenArr []Token, traitOccurences TraitFrequencyMap) TraitScoreMap {

	// TODO: parallelize if possible
	traitScores := make(TraitScoreMap)
	for traitGroup, traitValueFreqMap := range traitOccurences {
		// get the number of unique trait values
		uniqueTraitCt := len(traitValueFreqMap)
		// handle empty properties
		if _, found := traitScores[traitGroup]; !found {
			traitScores[traitGroup] = make(TraitValueScoreMap)
		}
		// TODO: parallelize if possible
		// multiply the amount of times this value occurs, by the amount of unique traits (values) within the group
		for traitValue, valueFreqCt := range traitValueFreqMap {
			traitValueRarityScore := decimal.New(1, 0).Div(decimal.NewFromInt(int64(valueFreqCt)).Mul(decimal.NewFromInt(int64(uniqueTraitCt))))
			traitScores[traitGroup][traitValue] = traitValueRarityScore
		}

	}

	return traitScores
}

// DEPRECATED
//
// # Calculate trait value probabilities, by group
//
// Input the frequency for the trait values of a single group & return its probability mapping
// func calcTraitValueProbsByGroup(traitGroup TraitValueFreqMap, activeTokenCount int) TraitValueProbMap {
// 	// get sum of all occurences
// 	var traitOccurrenceSum int

// 	// TODO: parallelize
// 	// m1 dumbfounded + m2 bored + m1 bored + m1 bored unshaven...
// 	for _, traitValueOccurrence := range traitGroup {
// 		traitOccurrenceSum += traitValueOccurrence
// 	}

// 	traitGroupProb := TraitValueProbMap{}
// 	// now that we have the sum, apply probability to each trait value
// 	for traitValue, traitValueOccurrence := range traitGroup {
// 		prob := float64(traitValueOccurrence) / float64(traitOccurrenceSum)
// 		traitGroupProb[traitValue] = prob
// 	}

// 	return traitGroupProb
// }

// DEPRECATED: update test
// Calculate the probabilities for all traits for the token collection in its current state.
//
// Creates a probability mapping of all traits, nested within their respective groups.
//
// Sources info from Skip Protocol's database.
//
// Note: OpenSea has a lack of consistency for contract versions with Akuri.
// func calcAllTraitValueProbs(traits TraitFrequencyMap, activeTokenCount int) TraitProbabilityMap {
// 	probMap := TraitProbabilityMap{}

// 	// TODO: parallelize
// 	for groupName, group := range traits {
// 		probMap[groupName] = calcTraitValueProbsByGroup(group, int(activeTokenCount))
// 	}

// 	return probMap
// }
