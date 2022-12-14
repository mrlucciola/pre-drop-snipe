package main

import "fmt"

// https://opensea.io/collection/azuki1
var collectionSlug = "we-asuki"
var tokenSlug = "azuki1"

func main() {
	// getTokens_(Collection_{
	// 	count: 10000,
	// 	url:   "azuki1",
	// })
	client := httpClient()
	res := getOsCollection(client, collectionSlug)

	// fmt.Println("res", res)
	fmt.Println("token ct:", res.Count)

	getTokens(tokenSlug, int(2))
	// fmt.Println("done")
	// for i := 0; i < len(tokenArr); i++ {

	// }
	// logger.Println(tokenArr)
}
