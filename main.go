/*
TODO:
- connect to opensea api
- connect to nftx
- connect to so rare api

*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type OpenSeaCollectionStatsResponse struct {
	Stats *OpenSeaStats `json:"stats"`
}

type OpenSeaCollectionResponse struct {
	Collection *OpenSeaCollection `json:"collection"`
}

type OpenSeaCollection struct {
	PrimaryAssets []OpenSeaPrimaryAssetsContract `json:"primary_asset_contracts"`
	TwitterName   string                         `json:"twitter_username"`
}

type OpenSeaPrimaryAssetsContract struct {
	CollectionAddress string `json:"address"`
	Name              string `json:"name"`
}

type OpenSeaStats struct {
	FloorPrice float64 `json:"floor_price"`
}

type LooksRareOrdersResponse struct {
	Data []LooksRareOrder `json:"data"`
}

type LooksRareOrder struct {
	Hash      string `json:"hash"`
	Price     string `json:"price"`
	Amount    int    `json:"amount"`
	StartTime int    `json:"startTime"`
	EndTime   int    `json:"endTime"`
}

func executeHttpRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return body, nil
}

func getOpenSeaCollectionAddress(slug string) (string, error) {
	body, err := executeHttpRequest("https://api.opensea.io/api/v1/collection/" + slug)

	collection := new(OpenSeaCollectionResponse)
	json.Unmarshal(body, collection)
	address := collection.Collection.PrimaryAssets[0].CollectionAddress
	return address, err
}

func getOpenSeaFloor(slug string) (float64, error) {
	body, err := executeHttpRequest("https://api.opensea.io/api/v1/collection/" + slug + "/stats")

	collectionStats := new(OpenSeaCollectionStatsResponse)
	json.Unmarshal(body, collectionStats)

	return collectionStats.Stats.FloorPrice, err
}

func getOpenSeaAsset() ([]byte, error) {
	resp, err := http.Get("https://api.opensea.io/api/v1/asset/0xb47e3cd837ddf8e4c57f05d70ab865de6e193bbb/1/?include_orders=false")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
	return body, nil
}

func getLooksRareBids(address string) ([]LooksRareOrder, error) {
	resp, err := http.Get("https://api.looksrare.org/api/v1/orders?tokenId=0&isOrderAsk=false&status%5B%5D=VALID&collection=" + address)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	looksRareResponse := new(LooksRareOrdersResponse)
	json.Unmarshal(body, looksRareResponse)

	return looksRareResponse.Data, nil
}

func main() {
	fmt.Println("Starting NFT Arb")
	floorPrice, _ := getOpenSeaFloor("otherdeed")
	fmt.Println("Floor price:", floorPrice, "ETH")

	address, _ := getOpenSeaCollectionAddress("otherdeed")
	fmt.Println("Address:", address)

	bids, _ := getLooksRareBids(address)
	for _, bid := range bids {
		fmt.Println("    Amount:", bid.Amount)
		fmt.Println("    Price", bid.Price)
	}

}
