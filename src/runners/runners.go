package runners

import (
	"fmt"
	"log"
	"nudam-ctrader-api/api"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/logger"
	"nudam-ctrader-api/strategy"
	"sync"
	"time"
)

// Start trading routines.
func TradeRoutines() {
	api := api.NewApi()
	err := api.Open()
	if err != nil {
		log.Panic(err)
	}
	defer api.Close()

	go RunnerReadMessage(api)

	go RunnerCheckStrategy()

	RunnerGetTrendbars(api)
}

// Func to start goroutine message reader.
func RunnerReadMessage(api api.CTraderAPI) {
	for {
		if err := api.ReadMessage(); err != nil {
			logger.LogError(err, "error reading message")
		}
	}
}

// Func to start goroutine strategy checker.
func RunnerCheckStrategy() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, symbol := range configs_helper.TraderConfiguration.CurrencyPairs {
			signal, err := strategy.SignalChecker(symbol)
			if err != nil {
				logger.LogError(err, "error getting data from mongodb")
				log.Panic(err)
			}

			if signal == strategy.Short {

			} else if signal == strategy.Long {

			}
		}
	}
}

// Func to start trendbars goroutines.
func RunnerGetTrendbars(api api.CTraderAPI) {
	var wg sync.WaitGroup
	for _, symbol := range configs_helper.TraderConfiguration.CurrencyPairs {
		for period := range configs_helper.TraderConfiguration.Periods {
			wg.Add(1)
			go func(symbol, period string) {
				ticker := time.NewTicker(30 * time.Second)
				defer ticker.Stop()

				for range ticker.C {
					if err := api.GetTrendbars(symbol, period); err != nil {
						logger.LogError(err, fmt.Sprintf("error getting trendbars for %s", symbol))
						log.Panic(err)
					}
				}
			}(symbol, period)
		}
	}
	wg.Wait()
}

// utils.LogMessage(fmt.Sprintf("%s - %s\n%s", symbol, period, resp))

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

// }
