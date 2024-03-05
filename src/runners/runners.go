package runners

import (
	"fmt"
	"log"
	"nudam-ctrader-api/api"
	"nudam-ctrader-api/strategy"
	"nudam-ctrader-api/types/ctrader"
	"nudam-ctrader-api/utils"
	"sync"
	"time"
)

// Starts trading routines.
func TradeRoutines(symbolPeriod map[string]string) {
	var wg sync.WaitGroup
	for symbol, period := range symbolPeriod {
		wg.Add(1)
		go func(symbol, period string) {
			defer wg.Done()

			apiCurrentPrice, err := api.NewApi()
			if err != nil {
				log.Panic(err)
			}
			apiTrendbars, err := api.NewApi()
			if err != nil {
				log.Panic(err)
			}
			defer apiCurrentPrice.Close()
			defer apiTrendbars.Close()

			time.Sleep(5 * time.Second)
			err = apiCurrentPrice.SendMsgSubscribeSpot(symbol)
			if err != nil {
				log.Panic(err)
			}

			for {
				prices, err := apiCurrentPrice.SendMsgReadMessage()
				if err != nil {
					if err.Error() == "websocket: close 1000 (normal): Bye" {
						apiCurrentPrice, prices = reconnectApiCurrentPrice(symbol, period)
					} else {
						log.Panic(err)
					}
				}

				closePrices, err := apiTrendbars.GetTrendbars(symbol, period)
				if err != nil {
					log.Panic(err)
				}

				resp, resBool := strategy.AreIntervalsTrendMatching(apiTrendbars, symbol, period)
				utils.LogMessage(fmt.Sprintf("%s - %s\n%s", symbol, period, resp))

				if resBool {
					currentTrend := strategy.GetTrendForPeriod(apiTrendbars, symbol, period)
					EMAs := strategy.GetEMAs(closePrices)

					trader := strategy.NewTrader(EMAs)
					signal := trader.CheckPriceBetweenEma26Ema50(float64(prices.Payload.Bid), float64(prices.Payload.Ask))

					if currentTrend == strategy.Uptrend && signal == strategy.Short {
						utils.LogMessage("short")
						break
					} else if currentTrend == strategy.Downtrend && signal == strategy.Long {
						utils.LogMessage("long")
						break
					}
				}

				time.Sleep(5 * time.Second)
			}
			// 	if spotEvent, ok := message.(*ctrader.ProtoOASpotEvent); ok && spotEvent.Symbol == symbol {
			// 		currentPrice := float64(spotEvent.Ask)
			// 		trader.UpdatePrice(currentPrice)
			// 	}

			// re, err := api.SendMsgNewOrder(int64(configs_helper.TraderConfiguration.OrderType["market"]), int64(configs_helper.TraderConfiguration.TradeSide["buy"]), int64(100000))
			// fmt.Println(string(re))
		}(symbol, period)
	}
	wg.Wait()
}

// Reconnects connection to subscribe spot.
func reconnectApiCurrentPrice(symbol, period string) (api.CTraderAPI, *ctrader.Message[ctrader.ProtoOASpotEvent]) {
	utils.LogMessage(fmt.Sprintf("%s; %s; %s", symbol, period, "attempting to reconnect due to normal WebSocket closure..."))
	apiCurrentPrice, err := api.NewApi()
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(5 * time.Second)
	err = apiCurrentPrice.SendMsgSubscribeSpot(symbol)
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(5 * time.Second)
	prices, err := apiCurrentPrice.SendMsgReadMessage()
	if err != nil {
		log.Panic(err)
	}

	return apiCurrentPrice, prices
}
