package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/getlantern/systray"
	krakenapi "github.com/wthorp/kraken-go-api-client"
)

func main() {
	systray.Run(onReady, func() {})
}

func onReady() {
	// read credentials file and connect
	credentials, err := ioutil.ReadFile("credentials.txt")
	if err != nil {
		systray.SetTitle("ERR")
		systray.SetTooltip(err.Error())
	}
	lines := strings.SplitN(string(credentials), "\n", 2)
	if len(lines) != 2 {
		systray.SetTitle("ERR")
		systray.SetTooltip(err.Error())
	}
	key, host := lines[0], lines[1]
	api := krakenapi.New(key, host)

	// set icon and allow quit
	iconBytes, err := getIcon("https://storj.io/img/favicon/favicon.ico")
	if err == nil {
		systray.SetTemplateIcon(iconBytes, iconBytes)
	}
	mQuitOrig := systray.AddMenuItem("Quit", "Quit Storj Ticker")

	// Update price on interval
	go func() {
		for {
			showPrice(api)
			select {
			case <-time.After(2 * time.Minute):
			case <-mQuitOrig.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func showPrice(api *krakenapi.KrakenAPI) {
	ticker, err := api.Ticker("STORJUSD")
	if err != nil {
		systray.SetTitle("ERR")
		systray.SetTooltip(err.Error())
	}
	systray.SetTitle((*ticker)["STORJUSD"].Bid[0])
}

func getIcon(url string) ([]byte, error) {
	var client http.Client
	resp, err := client.Get("https://storj.io/img/favicon/favicon.ico")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}
	return nil, fmt.Errorf("Bad response code")
}
