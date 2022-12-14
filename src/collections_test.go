package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestRetrieveCollection(t *testing.T) {
	// Load file
	filename := "../test/traits_asuki.json"
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Deserialize data from the JSON file
	var traitMap = CollectionTraits{}
	errUnmarshal := json.Unmarshal(content, &traitMap)
	if errUnmarshal != nil {
		t.Errorf("Failed ! unmarshalling json")
	}

	expected := CollectionPayload{
		Traits: traitMap,
		Count:  73750,
	}

	client := httpClient()
	output := getOsCollection(client, "we-asuki")

	if reflect.DeepEqual(expected, output) == false {

		t.Errorf("Failed ! \nexp: \n%v \nout: \n%v", expected, output)
	} else {
		t.Logf("Success !")
	}
}
