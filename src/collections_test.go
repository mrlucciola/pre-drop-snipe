package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestRetrieveCollection(t *testing.T) {
	// Load file
	filename := "../test/traits_asuki.json"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Deserialize data from the JSON file
	var traitMap = map[string]map[string]int{}
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

		fmt.Println("EEE:\n", expected, "\n ")
		fmt.Println("OOO:\n", output, "\n ")
		t.Errorf("Failed ! \nexp: \n%v \nout: \n%v", 1, 2)
	} else {
		t.Logf("Success !")
	}
}
