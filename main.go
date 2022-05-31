package main

import (
	"fmt"
	"log"
	"nftarb/api"
	"strconv"
)

func getPriceFloat(price string) float64 {
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Fatalln(err)
		return -1
	}
	return p / 1000000000000000000
}

func detectArbitrage(bid api.LooksRareOrder, floorPrice float64) {
	price := getPriceFloat(bid.Price)
	if price > floorPrice {
		fmt.Println("Arbitrage Detected, potential profit=", price-floorPrice, "ETH")
	}
}

func main() {
	fmt.Println("Starting NFT Arb")
	floorPrice, _ := api.GetOpenSeaFloor("otherdeed")
	fmt.Println("Floor price:", floorPrice, "ETH")

	address, _ := api.GetOpenSeaCollectionAddress("otherdeed")
	fmt.Println("Address:", address)

	bid, _ := api.GetLooksRareBestBid(address)

	detectArbitrage(*bid, floorPrice)

}
