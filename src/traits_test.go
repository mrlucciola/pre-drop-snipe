package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

// Check if the output has the same fields as input
func TestBuildTraitMap(t *testing.T) {
	// Load file
	filename := "../test/trait_prob_asuki.json"
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Deserialize data from the JSON file
	var expected = TraitProbabilityMap{}
	errUnmarshal := json.Unmarshal(content, &expected)
	if errUnmarshal != nil {
		t.Errorf("Failed ! unmarshalling json")
	}

	// get collection
	client := httpClient()
	res := getOsCollection(client, "we-asuki")
	// get trait-prob-map
	output := calcAllTraitValueProbs(res.Traits, int(res.Count))

	if reflect.DeepEqual(expected, output) == false {

		t.Errorf("Failed ! \nexp: \n%v \nout: \n%v", expected, output)
	} else {
		t.Logf("Success !")
	}
}
