package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nudam-ctrader-api/ctrader_api"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/helpers/ctrader_api_helper"
	"nudam-ctrader-api/strategy"
	"sync"
	"time"
)

var (
	config_path  = "./configs"
	symbolPeriod = map[string]string{"EURUSD": "m1", "AUDUSD": "m1"}
)

func tradeRoutines() {
	var wg sync.WaitGroup
	for symbol, period := range symbolPeriod {
		wg.Add(1)
		go func(symbol, period string) {
			trader := strategy.NewTrader()
			apiCurrentPrice, err := ctrader_api.NewApi()
			if err != nil {
				log.Panic(err)
			}
			err = apiCurrentPrice.SendMsgSubscribeSpot(symbol)
			if err != nil {
				log.Panic(err)
			}
			apiTrendbars, err := ctrader_api.NewApi()
			if err != nil {
				log.Panic(err)
			}

			for {
				prices, err := apiCurrentPrice.SendMsgReadMessage()
				if err != nil {
					log.Panic(err)
				}

				closePrices, err := apiTrendbars.GetTrendbars(symbol, period)
				if err != nil {
					log.Panic(err)
				}

				outPrices, err := json.Marshal(prices)
				if err != nil {
					panic(err)
				}

				ctrader_api_helper.LogMessage(fmt.Sprintf("%s; %s; %s", symbol, period, trader.GetEMAs(closePrices)))
				ctrader_api_helper.LogMessage(fmt.Sprintf("%s", string(outPrices)))

				time.Sleep(5 * time.Second)
			}
			// 	// Przetwórz przychodzącą wiadomość
			// 	if spotEvent, ok := message.(*ctrader.ProtoOASpotEvent); ok && spotEvent.Symbol == symbol {
			// 		currentPrice := float64(spotEvent.Ask) // Dostosuj do faktycznej struktury danych
			// 		trader.UpdatePrice(currentPrice)       // Aktualizuj cenę w traderze
			// 		// Dodatkowa logika handlowa...
			// 	}

			// re, err := api.SendMsgNewOrder(int64(configs_helper.TraderConfiguration.OrderType["market"]), int64(configs_helper.TraderConfiguration.TradeSide["buy"]), int64(100000))
			// fmt.Println(string(re))
		}(symbol, period)
	}
	wg.Wait()
}

func main() {
	err := configs_helper.InitializeConfig(config_path)
	if err != nil {
		log.Panic(err)
	}

	tradeRoutines()
}
