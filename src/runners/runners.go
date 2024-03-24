package runners

import (
	"log"
	"nudam-ctrader-api/api"
	"sync"
)

// Starts trading routines.
func TradeRoutines(symbolPeriod map[string]string) {
	var wg sync.WaitGroup
	for symbol, period := range symbolPeriod {
		wg.Add(1)
		go func(symbol, period string) {
			defer wg.Done()

			// apiCurrentPrice, err := api.NewApi()
			// if err != nil {
			// 	log.Panic(err)
			// }
			apiWhatever, err := api.NewApi()
			if err != nil {
				log.Panic(err)
			}
			// defer apiCurrentPrice.Close()
			defer apiWhatever.Close()

			// symbolEntity, err := apiWhatever.SaveSymbolEntity(symbol)
			// if err != nil {
			// 	log.Panic(err)
			// }

			// time.Sleep(5 * time.Second)
			// if err = apiWhatever.SendMsgSubscribeSpot(symbol); err != nil {
			// 	log.Panic(err)
			// }

			// for {
			// 	prices, err := apiWhatever.SendMsgReadMessage()
			// 	if err != nil {
			// 		if err.Error() == "websocket: close 1000 (normal): Bye" {
			// 			apiWhatever, prices = reconnectApiCurrentPrice(symbol, period)
			// 		} else {
			// 			log.Panic(err)
			// 		}
			// 	}

			// 	closePrices, err := apiWhatever.GetTrendbars(symbol, period)
			// 	if err != nil {
			// 		log.Panic(err)
			// 	}

			// 	resp, resBool := strategy.AreIntervalsTrendMatching(apiWhatever, symbol, period)
			// 	utils.LogMessage(fmt.Sprintf("%s - %s\n%s", symbol, period, resp))

			// 	if resBool {
			// 		currentTrend := strategy.GetTrendForPeriod(apiWhatever, symbol, period)
			// 		EMAs := strategy.GetEMAs(closePrices)

			// 		signal := strategy.CheckPriceBetweenEma26Ema50(float64(prices.Payload.Bid), float64(prices.Payload.Ask), EMAs)

			// 		if currentTrend == strategy.Downtrend && signal == strategy.Short {
			// 			balance, err := apiWhatever.SendMsgGetBalance()
			// 			if err != nil {
			// 				log.Panic(err)
			// 			}
			// 			if prices.Payload.Bid != 0 {
			// 				stopLossPips, volume := strategy.GetStopLossVolume(balance, EMAs, prices.Payload.Bid, symbolEntity)
			// 				utils.LogMessage(fmt.Sprintf("Opening position:\n%s - %s\n%v", symbol, period, volume))
			// 				re, err := apiWhatever.SendMsgNewOrder(symbol, int64(configs_helper.TraderConfiguration.OrderType["market"]), int64(configs_helper.TraderConfiguration.TradeSide["sell"]), volume, stopLossPips)
			// 				if err != nil {
			// 					utils.LogError(err, fmt.Sprintf("cant open position: %s", symbol))
			// 				}
			// 				utils.LogMessage(string(re))
			// 			}
			// 			break //read message
			// 		} else if currentTrend == strategy.Uptrend && signal == strategy.Long {
			// 			balance, err := apiWhatever.SendMsgGetBalance()
			// 			if err != nil {
			// 				log.Panic(err)
			// 			}
			// 			if prices.Payload.Ask != 0 {
			// 				stopLossPips, volume := strategy.GetStopLossVolume(balance, EMAs, prices.Payload.Ask, symbolEntity)
			// 				utils.LogMessage(fmt.Sprintf("Opening position:\n%s - %s\n%v", symbol, period, volume))
			// 				re, err := apiWhatever.SendMsgNewOrder(symbol, int64(configs_helper.TraderConfiguration.OrderType["market"]), int64(configs_helper.TraderConfiguration.TradeSide["buy"]), volume, stopLossPips)
			// 				if err != nil {
			// 					utils.LogError(err, fmt.Sprintf("cant open position: %s", symbol))
			// 				}
			// 				utils.LogMessage(string(re))
			// 			}

			// 			break
			// 		}
			// 	}

			// 	time.Sleep(5 * time.Second)
			// }
		}(symbol, period)
	}
	wg.Wait()
}

// Reconnects connection to subscribe spot.
// func reconnectApiCurrentPrice(symbol, period string) (api.CTraderAPI, *ctrader.Message[ctrader.ProtoOASpotEvent]) {
// 	utils.LogMessage(fmt.Sprintf("%s; %s; %s", symbol, period, "attempting to reconnect due to normal WebSocket closure..."))
// 	apiCurrentPrice, err := api.NewApi()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	time.Sleep(5 * time.Second)
// 	if err = apiCurrentPrice.SendMsgSubscribeSpot(symbol); err != nil {
// 		log.Panic(err)
// 	}
// 	time.Sleep(5 * time.Second)
// 	prices, err := apiCurrentPrice.SendMsgReadMessage()
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	return apiCurrentPrice, prices
// }
