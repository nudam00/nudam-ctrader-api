package runners

import (
	"log"
	"nudam-ctrader-api/api"
)

type IRunner interface {
	StartRoutines()
}

type Runner struct {
	api     api.CTraderAPI
	handler IHandler
}

func NewRunner() *Runner {
	return &Runner{
		api: api.NewApi(),
	}
}

// Start trading goroutines.
func (r *Runner) StartRoutines() {
	err := r.api.Open()
	if err != nil {
		log.Panic(err)
	}
	defer r.api.Close()

	r.handler = NewHandler()

	go r.handler.HandlerReadMessage(r.api)

	go r.handler.HandlerStrategy()

	r.handler.HandlerGetTrendbars(r.api)
}

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
