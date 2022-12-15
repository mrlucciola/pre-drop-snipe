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
	expected := []*Token{
		&Token{
			id: 0,
			traits: map[string]string{
				"Background": "Off White A",
				"Clothing":   "Pink Oversized Kimono",
				"Eyes":       "Striking",
				"Offhand":    "Monkey King Staff",
				"Type":       "Human",
				"Hair":       "Water",
				"Mouth":      "Frown",
			},
		},
		&Token{
			id: 0,
			traits: map[string]string{
				"Background": "Off White D",
				"Clothing":   "White Qipao with Fur",
				"Eyes":       "Daydreaming",
				"Offhand":    "Gloves",
				"Type":       "Human",
				"Hair":       "Pink Hairband",
				"Mouth":      "Lipstick",
			},
		},
	}

	output := getTokens("azuki1", 2)

	if reflect.DeepEqual(expected, output) {
		t.Errorf("Failed ! got %v want %v", output, expected)
	} else {
		t.Logf("Success !")
	}
}
