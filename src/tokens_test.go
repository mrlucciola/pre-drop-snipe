package main

import (
	"reflect"
	"testing"
)

func TestRetrieveToken(t *testing.T) {
	expected := Token{
		id: 0,
		traits: map[string]string{
			"Clothing":   "Pink Oversized Kimono",
			"Eyes":       "Striking",
			"Offhand":    "Monkey King Staff",
			"Type":       "Human",
			"Hair":       "Water",
			"Mouth":      "Frown",
			"Background": "Off White A",
		},
	}
	output := getToken("we-azuki", 0)

	if reflect.DeepEqual(expected, output) {
		t.Errorf("Failed ! got %v want %v", output, expected)
	} else {
		t.Logf("Success !")
	}
}

func TestRetrieveTokens(t *testing.T) {
	expected := Tokens{
		&Token{
			id: 0,
			traits: map[string]string{
				"Clothing":   "Pink Oversized Kimono",
				"Eyes":       "Striking",
				"Offhand":    "Monkey King Staff",
				"Type":       "Human",
				"Hair":       "Water",
				"Mouth":      "Frown",
				"Background": "Off White A",
			},
		},
	}

	output := getToken("azuki1", 0)

	if reflect.DeepEqual(false, output) {
		t.Errorf("Failed ! got %v want %v", output, expected)
	} else {
		t.Logf("Success !")
	}
}
