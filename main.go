package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"nftarb/api"
	"strconv"

	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")

}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Starting NFT Arb")
	floorPrice, _ := api.GetOpenSeaFloor("otherdeed")
	fmt.Println("Floor price:", floorPrice, "ETH")

	address, _ := api.GetOpenSeaCollectionAddress("otherdeed")
	fmt.Println("Address:", address)

	bid, _ := api.GetLooksRareBestBid(address)

	detectArbitrage(*bid, floorPrice)

	// Create a url.URL to connect to. `ws://` is non-encrypted websocket.
	urlStr := "wss://testnets-stream.openseabeta.com/socket?token=ae36b1b4131c421e8c84088ad48abe9b&vsn=2.0.0"
	endPoint, _ := url.Parse(urlStr)
	log.Println("Connecting to", endPoint)
	c, resp, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}
	_ = c
}
