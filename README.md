Kraken GO API Client
====================

A simple API Client for the [Kraken](https://www.kraken.com/ "Kraken") Trading platform.

This is a fork of [kraken-go-api-client](https://github.com/beldur/kraken-go-api-client/) which uses strings instead of typed symbol constants.

Example usage:

```go
package main

import (
	"fmt"
	"log"

	"github.com/wthorp/kraken-go-api-client"
)

func main() {
	api := krakenapi.New("KEY", "SECRET")

	ticker, err := api.Ticker("STORJUSD")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bid price: %s\n", (*ticker)["STORJUSD"].Bid[0])
}
```
