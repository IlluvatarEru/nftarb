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

type OpeanSeaCollectionResponse struct {
	Stats *OpenSeaStats `json:"stats"`
}

type OpenSeaStats struct {
	FloorPrice float64 `json:"floor_price"`
}
type Asset struct {
	TokenID  string `json:"token_id"`
	NumSales int64  `json:"num_sales"`
}

func getOpenSeaFloor(slug string) (float64, error) {
	resp, err := http.Get("https://api.opensea.io/api/v1/collection/" + slug + "/stats")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return -1, err
	}
	collection := new(OpeanSeaCollectionResponse)
	json.Unmarshal(body, collection)

	return collection.Stats.FloorPrice, nil
}

func getOpenSeaAsset() ([]byte, error) {
	resp, err := http.Get("https://api.opensea.io/api/v1/asset/0xb47e3cd837ddf8e4c57f05d70ab865de6e193bbb/1/?include_orders=false")
	if err != nil {
		fmt.Println("ERR1\n")
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERR2\n")
		log.Fatalln(err)
		return nil, err
	}

	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
	return body, nil
}

func main() {
	fmt.Println("Starting NFT Arb")
	floorPrice, _ := getOpenSeaFloor("trippin-ape-tribe-solana")
	fmt.Println(floorPrice)

}
