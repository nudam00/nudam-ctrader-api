package api

import (
	"encoding/json"
	"nudam-ctrader-api/external/mongodb"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/types/ctrader"
	"nudam-ctrader-api/utils"
)

// Get trendbars based on given symbol.
func (api *CTrader) GetTrendbars(symbol, period string) error {
	fromTimestamp, toTimestamp := utils.CalculateTimestamps(int(configs_helper.TraderConfiguration.Periods[period].NumberDays)) // then it will get the biggest possible amount of data
	periodId := configs_helper.TraderConfiguration.Periods[period].Value

	symbolId, err := mongodb.FindSymbolId(symbol)
	if err != nil {
		return err
	}
	count := utils.CalculateCountBars(period)

	protoOAGetTrendbarsReq := ctrader.Message[ctrader.ProtoOAGetTrendbarsReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooagettrendbarsreq"],
		Payload: ctrader.ProtoOAGetTrendbarsReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			FromTimestamp:       fromTimestamp,
			ToTimestamp:         toTimestamp,
			Period:              periodId,
			SymbolId:            symbolId,
			Count:               count,
		},
	}

	reqBytes, err := json.Marshal(protoOAGetTrendbarsReq)
	if err != nil {
		return err
	}
	api.sendMessage(reqBytes)

	return nil
}

// // Sends message to create new order.
// func (api *CTrader) SendMsgNewOrder(symbol string, orderType, tradeSide, volume, stopLoss int64) ([]byte, error) {
// 	symbolId, err := utils.FindSymbolId(symbol, api.symbolList)
// 	if err != nil {
// 		return nil, err
// 	}

// 	protoOANewOrderReq := ctrader.Message[ctrader.ProtoOANewOrderReq]{
// 		ClientMsgID: utils.GetClientMsgID(),
// 		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooaneworderreq"],
// 		Payload: ctrader.ProtoOANewOrderReq{
// 			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
// 			SymbolId:            symbolId,
// 			OrderType:           orderType,
// 			TradeSide:           tradeSide,
// 			Volume:              volume,
// 			RelativeStopLoss:    stopLoss,
// 			TrailingStopLoss:    true,
// 		},
// 	}

// 	if err := utils.SendMsg(api.ws, protoOANewOrderReq); err != nil {
// 		return nil, err
// 	}
// 	resp, err := utils.ReadMsg(api.ws)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaexecutionevent"], err); err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// // Gets current balance.
// func (api *CTrader) SendMsgGetBalance() (float64, error) {
// 	protoOATraderReq := ctrader.Message[ctrader.ProtoOATraderReq]{
// 		ClientMsgID: utils.GetClientMsgID(),
// 		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooatraderreq"],
// 		Payload: ctrader.ProtoOATraderReq{
// 			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
// 		},
// 	}

// 	if err := utils.SendMsg(api.ws, protoOATraderReq); err != nil {
// 		return 0, err
// 	}

// 	resp, err := utils.ReadMsg(api.ws)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooatraderres"], err); err != nil {
// 		return 0, err
// 	}

// 	var protoOATraderRes *ctrader.Message[ctrader.ProtoOATraderRes]
// 	if err = json.Unmarshal(resp, &protoOATraderRes); err != nil {
// 		return 0, err
// 	}

// 	return float64(protoOATraderRes.Payload.Trader.Balance) / 100.0, nil
// }
