package runners

import (
	"fmt"
	"log"
	"nudam-ctrader-api/api"
	"nudam-ctrader-api/external/mongodb"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/logger"
	"nudam-ctrader-api/strategy"
	"sync"
	"time"
)

type IHandler interface {
	HandlerReadMessage()
	HandlerStrategy()
	HandlerGetTrendbars()
}

type Handler struct {
	api api.CTraderAPI
}

func NewHandler(api api.CTraderAPI) IHandler {
	handler := new(Handler)
	handler.api = api
	return handler
}

// Func to start goroutine signal checker.
func (h *Handler) HandlerStrategy() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, symbol := range configs_helper.TraderConfiguration.CurrencyPairs {
			position, err := mongodb.FindPosition(symbol)
			if err != nil {
				logger.LogError(err, "error getting data from mongodb")
				log.Panic(err)
			}
			if position {
				continue
			}

			signal, err := strategy.SignalChecker(symbol)
			if err != nil {
				logger.LogError(err, "error getting data from mongodb")
				log.Panic(err)
			}

			switch signal {
			case strategy.Short, strategy.Long:
				h.openPosition(symbol, signal)
			default:
				continue
			}
		}
	}
}

func (h *Handler) openPosition(symbol string, signal strategy.Signal) {
	h.api.SendMsgGetBalance()
	h.api.SetOnBalanceUpdate(func(balance int64) {

	})
}

// Func to start goroutine message reader.
func (h *Handler) HandlerReadMessage() {
	for {
		if err := h.api.ReadMessage(); err != nil {
			logger.LogError(err, "error reading message")
		}
	}
}

// Func to start goroutine trendbar receiver.
func (h *Handler) HandlerGetTrendbars() {
	var wg sync.WaitGroup
	for _, symbol := range configs_helper.TraderConfiguration.CurrencyPairs {
		for period := range configs_helper.TraderConfiguration.Periods {
			wg.Add(1)
			go func(symbol, period string) {
				ticker := time.NewTicker(30 * time.Second)
				defer ticker.Stop()

				for range ticker.C {
					if err := h.api.GetTrendbars(symbol, period); err != nil {
						logger.LogError(err, fmt.Sprintf("error getting trendbars for %s", symbol))
						log.Panic(err)
					}
				}
			}(symbol, period)
		}
	}
	wg.Wait()
}

// func (h *Handler) closePosition(symbol string) {
// 	h.mu.Lock()
// 	defer h.mu.Unlock()
// 	h.positions[symbol] = false
// }

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
