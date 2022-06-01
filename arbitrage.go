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

func MonitorIfArbitrage(checks chan string) {
	for {
		res, ok := <-checks
		if ok {
			fmt.Println(res)
			add, _ := api.GetOpenSeaCollectionAddress(res)
			bid, _ := api.GetLooksRareBestBid(add)
			floor, _ := api.GetOpenSeaFloor(res)
			if bid != nil {
				detectArbitrage(*bid, floor)
			}
		} else {
			fmt.Println("ERR")
			break
		}
	}
}
