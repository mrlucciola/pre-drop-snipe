package main

import "log"

// https://opensea.io/collection/azuki1
var collectionSlug = "azuki1"

func main() {
	client := httpClient()
	res := getOsCollection(client, collectionSlug)
	log.Println("Res body:", string(res))

	getToken(collectionSlug, 0)
}
