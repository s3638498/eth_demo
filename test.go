package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

func BinanceTickerPrice(channel chan float64, client *binance_connector.WebsocketAPIClient) {
	response, err := client.NewTickerPriceService().Symbol("ETHUSDT").Do(context.Background())
	if err != nil {
		log.Fatalf("Error reading price: %v", err.Error())
	}
	result, err := strconv.ParseFloat(response.Result.Price, 64)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
	channel <- (math.Round(result*100) / 100)
	return
}

func UniswapTickerPrice(channel chan float64) {
	// this api https://tradingstrategy.ai/api/explorer/#/Trading%20pair/web_pair_details
	// returning tickerPrice from uniswap exchange using v2 protocol
	response, err := http.Get("https://tradingstrategy.ai/api/pair-details?exchange_slug=sushiswap&chain_slug=ethereum&pair_slug=ETH-USDT")

	if err != nil {
		log.Fatalf("Error: %v", err.Error())
		return
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var dat map[string]map[string]interface{}
	if err := json.Unmarshal(responseData, &dat); err != nil {
		log.Fatal(err)
	}
	result := dat["summary"]["usd_price_latest"].(float64)
	channel <- (math.Round(result*100) / 100)
	return
}

func ComparePrice(uniswap chan float64, binance chan float64) {
	for {
		uniswapPrice, open := <-uniswap
		if !open {
			break
		}
		binancePrice, open := <-binance
		if !open {
			break
		}

		if uniswapPrice > binancePrice {
			fmt.Printf("Buy on binance with price: %.2f, sell on uniswap with price: %.2f\n", binancePrice, uniswapPrice)
		} else if binancePrice > uniswapPrice {
			fmt.Printf("Buy on uniswap with price: %.2f, sell on binance with price: %.2f\n", uniswapPrice, binancePrice)
		} else {
			fmt.Println("They are equal price")
		}
	}
}

func main() {
	uniswap := make(chan float64, 1)
	binance := make(chan float64, 1)
	client := binance_connector.NewWebsocketAPIClient("", "", "wss://testnet.binance.vision/ws-api/v3")
	defer client.Close()
	err := client.Connect()
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	for {
		go UniswapTickerPrice(uniswap)
		go BinanceTickerPrice(binance, client)
		go ComparePrice(uniswap, binance)
		time.Sleep(1 * time.Second)
	}

}
