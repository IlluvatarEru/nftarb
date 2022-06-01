package main

import (
	"fmt"
	"nftarb/api"
)

func RunArbitrageDetections() {
	fmt.Println("Starting NFT Arb")
	floorPrice, _ := api.GetOpenSeaFloor("otherdeed")
	fmt.Println("Floor price:", floorPrice, "ETH")

	address, _ := api.GetOpenSeaCollectionAddress("otherdeed")
	fmt.Println("Address:", address)

	bid, _ := api.GetLooksRareBestBid(address)

	if bid != nil {
		detectArbitrage(*bid, floorPrice)
	}
	fmt.Println("---------------")
	fmt.Println("Loading Collections")
	var collections []Collection
	collectionNames := []string{"otherdeed", "goblintownwtf"}
	// add some collection
	collections = AddCollections(collectionNames, collections)
	PrintCollections(collections)
	fmt.Println("Starting to monitor")
	checks := make(chan string)
	for i := range collections {
		go collections[i].MonitorAndUpdateFloorPrice(checks)
	}
	fmt.Println("here")
	MonitorIfArbitrage(checks)
}

func main() {
	RunArbitrageDetections()
}
