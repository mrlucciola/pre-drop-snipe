package main

import (
	"reflect"
	"testing"
)

func TestRetrieveToken(t *testing.T) {
	expected := Token{
		id: 0,
		traits: TokenTraitMap{
			"Background": "Off White A",
			"Clothing":   "Pink Oversized Kimono",
			"Eyes":       "Striking",
			"Offhand":    "Monkey King Staff",
			"Type":       "Human",
			"Hair":       "Water",
			"Mouth":      "Frown",
		},
	}

	client := httpClient()
	output := getToken(client, "azuki1", 0)

	if reflect.DeepEqual(expected, output) == true {
		t.Errorf("Failed ! got %v want %v", output, expected)
	} else {
		t.Logf("Success !")
	}
}

func TestRetrieveTokens(t *testing.T) {
	expected := []*Token{
		{
			id: 0,
			traits: TokenTraitMap{
				"Background": "Off White A",
				"Clothing":   "Pink Oversized Kimono",
				"Eyes":       "Striking",
				"Offhand":    "Monkey King Staff",
				"Type":       "Human",
				"Hair":       "Water",
				"Mouth":      "Frown",
			},
		},
		{
			id: 1,
			traits: TokenTraitMap{
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

	if reflect.DeepEqual(expected, output) == true {
		t.Errorf("Failed ! got %v want %v", output, expected)
	} else {
		t.Logf("Success !")
	}
}
