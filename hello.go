package main

import (
	"context"
	"fmt"
	"log"

	binance_connector "github.com/binance/binance-connector-go"
)

func TickerPrice() {
	client := binance_connector.NewWebsocketAPIClient("", "", "wss://testnet.binance.vision/ws-api/v3")
    defer client.Close()
	err := client.Connect()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer client.Close()

	response, err := client.NewTickerPriceService().Symbol("ETHUSDT").Do(context.Background())
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println(binance_connector.PrettyPrint(response))


}

func main() {
    TickerPrice()
}