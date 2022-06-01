package main

import (
	"fmt"
	"nftarb/api"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Collection struct {
	Name                string
	Address             string
	FloorPrice          float64
	HighestBiddingPrice float64
	Bidders             []Bidder
	Sellers             []Seller
}

type Bidder struct {
	Address             common.Address
	BiddingPrice        float64
	Platform            string
	IsBiddingForAnyItem bool
}

type Seller struct {
	Address          common.Address
	AskingPrice      float64
	Platform         string
	IsSellingAtFloor bool
}

func (collection *Collection) GetHighestBid() float64 {
	highestBid := 0.0
	for _, bidder := range collection.Bidders {
		if bidder.BiddingPrice > highestBid {
			highestBid = bidder.BiddingPrice
		}
	}
	return highestBid
}

func (collection *Collection) GetLowestAsk() float64 {
	lowestAsk := 0.0
	for _, seller := range collection.Sellers {
		if seller.AskingPrice < lowestAsk {
			lowestAsk = seller.AskingPrice
		}
	}
	return lowestAsk
}

func (collection *Collection) CheckIfArbitrage() bool {
	return collection.GetHighestBid() > collection.GetLowestAsk()
}

func (collection *Collection) AddBidder(bidder Bidder) {
	collection.Bidders = append(collection.Bidders, bidder)
	if bidder.BiddingPrice > collection.HighestBiddingPrice {
		collection.HighestBiddingPrice = bidder.BiddingPrice
	}
}

func (collection *Collection) AddSeller(seller Seller) {
	collection.Sellers = append(collection.Sellers, seller)
	if seller.AskingPrice < collection.FloorPrice {
		collection.FloorPrice = seller.AskingPrice
	}
}

/*
How to organize all of that ?
Should each collection have a list of buyer and a list of seller and then we update theses lists?
*/
func MaintainListOfSellers(sellers []Seller) {
	/*
		Maintain a list of sellers by looking at new listings
	*/
}

func MaintainListOfBidders(bidders []Bidder) {
	/*
		We want to maintain a list of bidders from different platforms
		Listen for new orders on looks rare etc

	*/
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
