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
	result, err := api.Query("Ticker", map[string]string{
		"pair": "XXBTZEUR",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %+v\n", result)

	// There are also some strongly typed methods available
	ticker, err := api.Ticker("XXBTZEUR")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ticker["XXBTZEUR"].OpeningPrice)
}
```
