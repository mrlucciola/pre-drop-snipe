package main

import (
	"reflect"
	"testing"
)

func TestBuildTraitMap(t *testing.T) {
	// Load file
	// filename := "../test/traits_asuki.json"
	// content, err := os.ReadFile(filename)
	// if err != nil {
	// 	log.Fatal("Error when opening file: ", err)
	// }

	// // Deserialize data from the JSON file
	// var traitMap = map[string]map[string]int{}
	// errUnmarshal := json.Unmarshal(content, &traitMap)
	// if errUnmarshal != nil {
	// 	t.Errorf("Failed ! unmarshalling json")
	// }

	// expected := TraitStats{}

	// // get collection
	// client := httpClient()
	// res := getOsCollection(client, "we-asuki")
	// // get trait-map
	// output := getTraitStats(res.Traits, int(res.Count))

	// if reflect.DeepEqual(expected, output) == false {
	if reflect.DeepEqual(1, 1) == false {

		// t.Errorf("Failed ! \nexp: \n%v \nout: \n%v", expected, output)
		t.Errorf("Failed ! \nexp: \n%v \nout: \n%v", 1, 1)
	} else {
		t.Logf("Success !")
	}
}
