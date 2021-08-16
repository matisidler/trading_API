package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

type Model struct {
	CurrentPrice       json.Number `json:"c"`
	Change             json.Number `json:"d"`
	PercentChange      json.Number `json:"dp"`
	HighPriceOfDay     json.Number `json:"h"`
	LowPriceOfDay      json.Number `json:"l"`
	OpenPriceOfDay     json.Number `json:"o"`
	PreviousPriceClose json.Number `json:"pc"`
}

func main() {
	var status bool
	for {

		cfg := finnhub.NewConfiguration()
		cfg.AddDefaultHeader("X-Finnhub-Token", "c4d4lvqad3icnt8r9eng")
		finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

		ahora := time.Now()
		btc10days, _, err := finnhubClient.CryptoCandles(context.Background()).Symbol("BINANCE:BTCUSDT").Resolution("1").From(ahora.Add(10 * time.Minute * -1).Unix()).To(time.Now().Unix()).Execute()
		if err != nil {
			fmt.Println(err)
		}

		closePrice := btc10days.GetC()
		var mediaMovil10 float32
		for _, valor := range closePrice {
			mediaMovil10 = mediaMovil10 + valor
		}
		mediaMovil10 = mediaMovil10 / float32(len(closePrice))
		fmt.Printf("media movil 10 períodos: %v\n", mediaMovil10)

		btc20days, _, err := finnhubClient.CryptoCandles(context.Background()).Symbol("BINANCE:BTCUSDT").Resolution("1").From(ahora.Add(20 * time.Minute * -1).Unix()).To(time.Now().Unix()).Execute()
		if err != nil {
			fmt.Println(err)
		}
		closePrice = btc20days.GetC()
		var mediaMovil20 float32
		for _, valor := range closePrice {
			mediaMovil20 = mediaMovil20 + valor
		}
		mediaMovil20 = mediaMovil20 / float32(len(closePrice))
		fmt.Printf("media movil 20 períodos: %v\n", mediaMovil20)

		if mediaMovil10 > mediaMovil20 {
			if status != true {
				fmt.Println("COMPRAR")
				status = true
			} else if status == true {
				fmt.Println("YA ESTÁS LONG")
			}

		}

		if mediaMovil10 < mediaMovil20 {
			if status != false {
				fmt.Println("VENDER")
				status = false
			} else if status == false {
				fmt.Println("YA ESTÁS SHORT")
			}

		}

		btcvolumeInfo, _, err := finnhubClient.CryptoCandles(context.Background()).Symbol("BINANCE:BTCUSDT").Resolution("D").From(ahora.Add(24 * time.Hour * -1).Unix()).To(time.Now().Unix()).Execute()
		if err != nil {
			fmt.Println(err)
		}
		btcvolume := btcvolumeInfo.GetV()
		res, err := http.Get("http://finnhub.io/api/v1/quote?symbol=BINANCE:BTCUSDT&token=c4d4lvqad3icnt8r9eng")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		data := Model{}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("----------BINANCE:BTCUSDT-----------\n, current price: %s\n, change: %s\n, percent change: %s\n, highest price of day: %s\n, lowest price of day: %s\n, open price of day: %s\n, previous price close: %s\n, volume: %v\n \n---------------------------\n", data.CurrentPrice, data.Change, data.PercentChange, data.HighPriceOfDay, data.LowPriceOfDay, data.OpenPriceOfDay, data.PreviousPriceClose, btcvolume)

		time.Sleep(5 * time.Second)
	}

}
