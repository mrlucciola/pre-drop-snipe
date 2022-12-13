# Assignment

### Goal:

> Write a program to download the trait metadata for a collection, compute the rarity scores for all tokens, sort the list by rarity, and output the rarity scores of the top 5 tokens in golang.

### Stub program:

The provided stub program includes data structures to represent tokens and the rarity scores for tokens, as well as a helper method for retrieving a token’s metadata from our server.

### Hints:

1. You should leverage concurrency and golang’s built in concurrency primitives.
2. If you do use concurrency, you should make the maximum number of threads configurable
3. As a sanity check, the most rare token is 6088 with a rarity score of 0.00856

<!-- rarityc,t =∑∑nc 1(vj,i ==vt,i)⋅oi -->