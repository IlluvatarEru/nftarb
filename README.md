# nftarb

### Goal
The goal is to detect if someone on a platform makes a bid at a price that is higher than the floor price and then to bundle a buy and sell txs together to arb it.

### How

For each collection we need to maintain the floor price. It does not change too often so we can query it via OpenSea REST api.
Then we need to stream collection bids from Looks Rare, we can do REST queries for that.