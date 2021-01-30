package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	krakenapi "github.com/wthorp/kraken-go-api-client"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// read credentials file and connect
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error determining home path")
	}
	credentials, err := ioutil.ReadFile(filepath.Join(home, ".krakenCreds"))
	if err != nil {
		log.Fatalf("Error reading credentials file %s\n", err)
	}
	lines := strings.SplitN(string(credentials), "\n", 2)
	if len(lines) != 2 {
		log.Fatal("Credentials file appears invalid")
	}
	key, host := lines[0], lines[1]

	// query STORJ price in USD
	api := krakenapi.New(key, host)
	ticker, err := api.Ticker("STORJUSD")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bid price: %s\n", (*ticker)["STORJUSD"].Bid[0])
}
