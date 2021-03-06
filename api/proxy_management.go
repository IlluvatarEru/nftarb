package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	resourceUrl = "https://ipinfo.io"
	proxyHost   = "gate.smartproxy.com:7000"
	username    = "username"
	password    = "password"
)

func TestProxy() {
	proxyUrl := &url.URL{
		Scheme: "http",
		User:   url.UserPassword(username, password),
		Host:   proxyHost,
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	resp, err := client.Get(resourceUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}

	fmt.Println(body)
}
