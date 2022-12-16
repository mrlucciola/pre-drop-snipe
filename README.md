# About

This application calculates the rarity of collectible NFTs (probability of tokens, provided their attributes) relative to all other NFTs within the same collection.

This application performs two main tasks:

1. Fetch all tokens of a single collection from a single source (Skip Protocol's database at the time of writing)
1. Calculate probabilities/rarities for each token relative to the rest of the collection

________________________________

## Running this code

1. Navigate to repository home directory\
   `cd ~/path_to_repo/`
1. Run command: `go run ./src`

Docker deployment soon.

________________________________

## Environment

| software    | command          | version                 |
| ----------- | ---------------- | ----------------------- |
| macOS       | `About this Mac` | `Ventura 13.0.1`        |
| Go (Golang) | `go version`     | `go1.19.4 darwin/amd64` |

________________________________

## Concurrency Error

I was receiving +300,000-line error logs when attempting concurrently call 10,000 endpoints, with and without
breaking down these requests into 1000-call chunks.

log: 
```
runtime/cgo: pthread_create failed: Resource temporarily unavailable
SIGABRT: abort
PC=0x7ff80e2c530e m=2256 sigcode=0

goroutine 0 [idle]:
runtime: g 0: unknown pc 0x7ff80e2c530e
stack: frame={sp:0x7000532b3b28, fp:0x0} stack=[0x700053234340,0x7000532b3f40)
```



## Assignment

### Goal:

> Write a program to download the trait metadata for a collection, compute the rarity scores for all tokens, sort the list by rarity, and output the rarity scores of the top 5 tokens in golang.

### Stub program:

The provided stub program includes data structures to represent tokens and the rarity scores for tokens, as well as a helper method for retrieving a token’s metadata from our server.

### Hints:

1. You should leverage concurrency and golang’s built in concurrency primitives.
2. If you do use concurrency, you should make the maximum number of threads configurable
3. As a sanity check, the most rare token is 6088 with a rarity score of 0.00856

<!-- rarityc,t =∑∑nc 1(vj,i ==vt,i)⋅oi -->
