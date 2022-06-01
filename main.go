package main

import (
	"fmt"
	"log"
	"nftarb/api"
	"strconv"
	"time"
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

type Collection struct {
	Name       string
	Address    string
	FloorPrice float64
}

func AddCollections(collectionNames []string, collections []Collection) []Collection {
	for _, collectionName := range collectionNames {
		address, _ := api.GetOpenSeaCollectionAddress(collectionName)
		floorPrice, _ := api.GetOpenSeaFloor(collectionName)
		collection := Collection{
			Name:       collectionName,
			Address:    address,
			FloorPrice: floorPrice,
		}
		collections = append(collections, collection)
	}
	return collections
}

func PrintCollections(collections []Collection) {
	for _, collection := range collections {
		PrintCollection(collection)
	}
}

func PrintCollection(collection Collection) {
	fmt.Println("--------------------")
	fmt.Println(collection.Name)
	fmt.Println(collection.Address)
	fmt.Println(collection.FloorPrice)
	fmt.Println("--------------------")
}

func (collection *Collection) MonitorAndUpdateFloorPrice(checks chan string) {
	for {
		floorPrice, err := api.GetOpenSeaFloor(collection.Name)
		if floorPrice != collection.FloorPrice {
			fmt.Println("Floor Price Update for", collection.Name)
			fmt.Println("Previous Floor price", collection.FloorPrice)
			fmt.Println("New Floor price", floorPrice)
			collection.FloorPrice = floorPrice
			checks <- collection.Name
		}
		if err != nil {
			fmt.Println("ERRRRR")
			fmt.Println(err)
			break
		}
		time.Sleep(2)
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

func main() {
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
	PrintCollections(collections)

}
