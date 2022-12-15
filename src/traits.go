package main

import (
	"fmt"
)

// Calculate the probability for all types of a single trait
//
//	ex. Mouth: {
//	  "m2 bored": 991,
//	  "m1 bored": 3313,
//	  "m1 bored unshaven": 2257,
//	  "m1 dumbfounded": 750,
//	}
//
//	into -> Mouth: {
//	  "m2 bored": 0.1355,
//	  "m1 bored": 0.4531,
//	  "m1 bored unshaven": 0.3087,
//	  "m1 dumbfounded": 0.1025,
//	}
type TraitValueMap map[string]int
type TraitCategoryProb map[string]float64

// Stats for all traits
type TraitProbabilityMap map[string]TraitCategoryProb

func getTraitStats(traitCategory TraitValueMap, activeTokenCount int) TraitCategoryProb {
	// get sum of all occurences
	var traitOccurrenceSum int
	// m1 dumbfounded + m2 bored + m1 bored + m1 bored unshaven...
	for _, traitValueOccurrence := range traitCategory {
		traitOccurrenceSum += traitValueOccurrence
	}
	// fmt.Println("Total: ", activeTokenCount, "sum: ", traitOccurrenceSum)
	// validate this sum equals the number of active tokens
	// if traitOccurrenceSum != activeTokenCount {
	// 	logger.Panic("Trait counts do not line up")
	// }

	traitCategoryProb := TraitCategoryProb{}
	// now that we have the sum, apply probability to each trait value
	for traitValue, traitValueOccurrence := range traitCategory {
		prob := float64(traitValueOccurrence) / float64(traitOccurrenceSum)
		traitCategoryProb[traitValue] = prob
	}
	return traitCategoryProb
}

// OpenSea input
func getAllTraitStats(traits CollectionTraits, activeTokenCount int) TraitProbabilityMap {
	traitStats := TraitProbabilityMap{}

	// parallelize
	for key, item := range traits {
		traitStats[key] = getTraitStats(item, int(activeTokenCount))
	}

	return traitStats
}

type TraitOccurenceMap map[string]TraitValueMap

// type TraitOccurenceArr []struct {
// 	string
// 	TraitValueMap
// }

/*

Token { category1: trait_1x, category2: trait_2x }

[
	{cat: "CLOTHING", [{red: 11}, {blue: 5}, {green: 8}]}
]

*/

/*
Iterate through the list of tokens pulled from the Skip Protocol API
  - Input is `Token` array - from Skip Protocol API
*/
func getAllTraitStatsSkip(tokenArr []*Token, activeTokenCount int) TraitProbabilityMap {
	traitProbabilities := TraitProbabilityMap{}
	traitOccurences := TraitOccurenceMap{}

	// parallelize if possible
	for _, token := range tokenArr {

		// iterate thru the traits, add occurences to the map
		for traitCategory, traitValue := range token.traits {
			if _, found := traitOccurences[traitCategory]; !found {
				traitOccurences[traitCategory] = TraitValueMap{}
			}

			traitOccurences[traitCategory][traitValue] += 1
		}
	}

	// parallelize if possible
	for traitCategory, traitValueMap := range traitOccurences {
		// sum all value counts for the category
		categorySum := 0
		for _, valueCt := range traitValueMap {
			categorySum += valueCt
		}

		// handle empty properties
		if _, found := traitProbabilities[traitCategory]; !found {
			traitProbabilities[traitCategory] = TraitCategoryProb{}
		}

		// iter thru each value, divide to get prob for each value, assign to map
		for traitValue, valueCt := range traitValueMap {
			traitProbabilities[traitCategory][traitValue] = float64(valueCt) / float64(categorySum)
		}
	}

	fmt.Println(traitProbabilities)
	return traitProbabilities
}
