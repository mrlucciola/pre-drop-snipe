package main

import "log"

// https://opensea.io/collection/mutant-ape-yacht-club
var collectionSlug = "mutant-ape-yacht-club"

func main() {
	client := httpClient()
	res := getOsCollection(client, collectionSlug)
	log.Println("Res body:", string(res))
}
