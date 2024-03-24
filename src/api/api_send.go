package api

import (
	"encoding/json"
	"fmt"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/logger"
	"nudam-ctrader-api/types/ctrader"
)

// Get trendbars based on given symbol.
func (api *CTrader) GetTrendbars(symbol string) error {
	logger.LogMessage("getting trendbars...")

	fromTimestamp, toTimestamp := utils.CalculateTimestamps(int(configs_helper.TraderConfiguration.Periods["mn1"].NumberDays))
	periodId := configs_helper.TraderConfiguration.Periods[period].Value
	symbolId, err := utils.FindSymbolId(symbol, api.symbolList)
	if err != nil {
		return nil, err
	}
	count := utils.CalculateCountBars(period)

	resp, err := api.sendMsgTrendbars(fromTimestamp, toTimestamp, periodId, symbolId, count)
	if err != nil {
		return nil, err
	}

	var protoOAGetTrendbarsRes ctrader.Message[ctrader.ProtoOAGetTrendbarsRes]
	if err = json.Unmarshal(resp, &protoOAGetTrendbarsRes); err != nil {
		return nil, err
	}

	var closePrices []float64
	for _, bar := range protoOAGetTrendbarsRes.Payload.Trendbar {
		closePrice := bar.Low + int64(bar.DeltaClose)
		closePrices = append(closePrices, float64(closePrice))
	}

	return closePrices, nil
}

// Sends message to get current trendbars.
func (api *CTrader) sendMsgTrendbars(fromTimestamp, toTimestamp, periodId, symbolId int64, count uint32) ([]byte, error) {
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

	if err := utils.SendMsg(api.ws, protoOAGetTrendbarsReq); err != nil {
		return nil, err
	}
	resp, err := utils.ReadMsg(api.ws)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(resp))
	resp, err = utils.ReadMsg(api.ws)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp))

	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooagettrendbarsres"], err); err != nil {
		return nil, err
	}

	return resp, nil
}

// Sends message to get current price.
// func (api *CTrader) SendMsgReadMessage() (*ctrader.Message[ctrader.ProtoOASpotEvent], error) {
// 	resp, err := utils.ReadMsg(api.ws)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaspotevent"], err); err != nil {
// 		return nil, err
// 	}

// 	var protoOASpotEvent *ctrader.Message[ctrader.ProtoOASpotEvent]
// 	if err = json.Unmarshal(resp, &protoOASpotEvent); err != nil {
// 		return nil, err
// 	}

// 	return protoOASpotEvent, nil
// }

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
