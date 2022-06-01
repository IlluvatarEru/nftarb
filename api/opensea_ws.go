package api

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

/*
Stream new listings from https://docs.opensea.io/reference/stream-api-overview
wss://stream.openseabeta.com/socket/websocket?token=

*/

func GetOpenSeaWSHost(live bool) string {
	if live {
		return "stream.openseabeta.com"
	} else {
		return "testnets-stream.openseabeta.com"
	}
}
func StartSimpleExample() {
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
func StartWSConnection(live bool) {
	openSeaHost := GetOpenSeaWSHost(live)
	openSeaPath := "/socket/websocket"
	//apiKey := "ae36b1b4131c421e8c84088ad48abe9b"

	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "wss", Host: "marketdata.tradermade.com", Path: "/feedadv"}
	u2 := url.URL{Scheme: "wss", Host: openSeaHost, Path: openSeaPath}

	log.Printf("connecting to %s", u.String())
	log.Printf("connecting to %s", u2.String())
	//c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	c, resp, err := websocket.DefaultDialer.Dial(u2.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}

	//When the program closes close the connection
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			if string(message) == "Connected" {
				log.Printf("Send Sub Details: %s", message)
				messageOut <- `{"userKey":"wsChqp-7X80Q0jGuCSWg", "symbol":"EURUSD"}`
			}
		}

	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case m := <-messageOut:
			log.Printf("Send Message %s", m)
			err := c.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
