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
type TraitCategory map[string]int
type TraitCategoryProb map[string]float64

// Stats for all traits
type TraitStats map[string]TraitCategoryProb

func getTraitStats(traitCategory TraitCategory, activeTokenCount int) TraitCategoryProb {
	// get sum of all occurences
	var traitOccurrenceSum int
	// m1 dumbfounded + m2 bored + m1 bored + m1 bored unshaven...
	for _, traitValueOccurrence := range traitCategory {
		traitOccurrenceSum += traitValueOccurrence
	}
	fmt.Println("Total: ", activeTokenCount, "sum: ", traitOccurrenceSum)
	// validate this sum equals the number of active tokens
	//
	// if traitOccurrenceSum != activeTokenCount {
	// 	logger.Panic("Trait counts do not line up")
	// }
	// default return

	traitCategoryProb := TraitCategoryProb{}
	// now that we have the sum, apply probability to each trait value
	for traitValue, traitValueOccurrence := range traitCategory {
		prob := float64(traitValueOccurrence) / float64(traitOccurrenceSum)
		traitCategoryProb[traitValue] = prob
		fmt.Println("prob:", traitCategoryProb[traitValue], " = ", traitValueOccurrence, "/", traitOccurrenceSum, "trait:", traitValue)
	}
	return traitCategoryProb
}
func getAllTraitStats(traits CollectionTraits, activeTokenCount int) TraitStats {
	traitStats := TraitStats{}

	// parallelize
	for key, item := range traits {
		fmt.Println("\n\nKey: ", key)
		traitStats[key] = getTraitStats(item, int(activeTokenCount))
	}

	// default return
	return traitStats
}
