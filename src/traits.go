package main

// Distribution map for trait values of a single group
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
type TraitValueDistMap map[string]int

// Probability map for trait values of a single group.
//
// Trait value probabilities are calculated from the distribution mapping.
//
//	Mouth: {
//	  "m2 bored":          0.1355 ( 991 / 7311),
//	  "m1 bored":          0.4531 (3313 / 7311),
//	  "m1 bored unshaven": 0.3087 (2257 / 7311),
//	  "m1 dumbfounded":    0.1025 ( 750 / 7311),
//	}
type TraitValueProbMap map[string]float64

type TraitDistributionMap map[string]TraitValueDistMap

// Probabilities for all traits
type TraitProbabilityMap map[string]TraitValueProbMap

// Calculate trait value probabilities, by group
//
// Input the distribution for the trait values of a single group & return its probability mapping
func calcTraitValueProbsByGroup(traitGroup TraitValueDistMap, activeTokenCount int) TraitValueProbMap {
	// get sum of all occurences
	var traitOccurrenceSum int

	// TODO: parallelize
	// m1 dumbfounded + m2 bored + m1 bored + m1 bored unshaven...
	for _, traitValueOccurrence := range traitGroup {
		traitOccurrenceSum += traitValueOccurrence
	}

	traitGroupProb := TraitValueProbMap{}
	// now that we have the sum, apply probability to each trait value
	for traitValue, traitValueOccurrence := range traitGroup {
		prob := float64(traitValueOccurrence) / float64(traitOccurrenceSum)
		traitGroupProb[traitValue] = prob
	}

	return traitGroupProb
}

// Calculate the probabilities for all traits for the token collection in its current state.
//
// Creates a probability mapping of all traits, nested within their respective groups.
//
// Sources info from Skip Protocol's database.
//
// Note: OpenSea has a lack of consistency for contract versions with Akuri.
func calcAllTraitValueProbs(traits CollectionTraits, activeTokenCount int) TraitProbabilityMap {
	probMap := TraitProbabilityMap{}

	// TODO: parallelize
	for groupName, group := range traits {
		probMap[groupName] = calcTraitValueProbsByGroup(group, int(activeTokenCount))
	}

	return probMap
}

/*
Iterate through the list of tokens pulled from the Skip Protocol API
  - Input is `Token` array - from Skip Protocol API
*/
func getAllTraitStatsSkip(tokenArr []*Token, activeTokenCount int) TraitProbabilityMap {
	traitProbabilities := TraitProbabilityMap{}
	traitOccurences := TraitDistributionMap{}

	// TODO: parallelize if possible
	for _, token := range tokenArr {

		// iterate thru the traits, add occurences to the map
		for traitGroup, traitValue := range token.traits {
			if _, found := traitOccurences[traitGroup]; !found {
				traitOccurences[traitGroup] = TraitValueDistMap{}
			}

			traitOccurences[traitGroup][traitValue] += 1
		}
	}

	// TODO: parallelize if possible
	for traitGroup, traitValueMap := range traitOccurences {
		// sum all value counts for the group
		groupSum := 0
		for _, valueCt := range traitValueMap {
			groupSum += valueCt
		}

		// handle empty properties
		if _, found := traitProbabilities[traitGroup]; !found {
			traitProbabilities[traitGroup] = TraitValueProbMap{}
		}

		// iter thru each value, divide to get prob for each value, assign to map
		for traitValue, valueCt := range traitValueMap {
			traitProbabilities[traitGroup][traitValue] = float64(valueCt) / float64(groupSum)
		}
	}

	return traitProbabilities
}

// Calculate the probability for an individual token.
//
// Call after any update to token datastore, or token datastore init
func calculateRarity(tokenToEval Token, probabilityMap TraitProbabilityMap, tokenRarityArr []float64) float64 {
	// initialize at 1
	prob := 1.

	// iterate thru each trait on the token
	for traitGroup, traitValue := range tokenToEval.traits {
		// look up value in the map
		prob *= probabilityMap[traitGroup][traitValue]
	}

	// add to datastore
	tokenRarityArr[tokenToEval.id] = roundFloat(prob, 15)

	return prob
}
