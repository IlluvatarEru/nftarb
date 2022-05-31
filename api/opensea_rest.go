package api

import (
	"encoding/json"
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

func GetOpenSeaAsset() ([]byte, error) {
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

func GetOpenSeaCollectionAddress(slug string) (string, error) {
	body, err := ExecuteHttpRequest("https://api.opensea.io/api/v1/collection/" + slug)

	collection := new(OpenSeaCollectionResponse)
	json.Unmarshal(body, collection)
	address := collection.Collection.PrimaryAssets[0].CollectionAddress
	return address, err
}

func GetOpenSeaFloor(slug string) (float64, error) {
	body, err := ExecuteHttpRequest("https://api.opensea.io/api/v1/collection/" + slug + "/stats")

	collectionStats := new(OpenSeaCollectionStatsResponse)
	json.Unmarshal(body, collectionStats)

	return collectionStats.Stats.FloorPrice, err
}
