package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

func GetLooksRareBestBid(address string) (*LooksRareOrder, error) {
	fmt.Println(address)
	resp, err := http.Get("https://api.looksrare.org/api/v1/orders?tokenId=0&isOrderAsk=false&status%5B%5D=VALID&sort=PRICE_ASC&collection=" + address)
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
	if len(looksRareResponse.Data) > 0 {
		return &looksRareResponse.Data[0], nil
	} else {
		fmt.Println("No bid on Looks rare")
		return nil, nil
	}
}
