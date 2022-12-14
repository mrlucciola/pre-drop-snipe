package main

import "fmt"

// https://opensea.io/collection/azuki1
var collectionSlug = "we-asuki"
var tokenSlug = "azuki1"

const COLOR_YELLOW = "\033[33m"
const COLOR_GREEN = "\033[32m"
const COLOR_RED = "\033[31m"
const COLOR_RESET = "\033[0m"

func main() {
	client := httpClient()
	res := getOsCollection(client, collectionSlug)

	fmt.Println("token ct:", res.Count)

	getTokens(tokenSlug, int(2))
}
