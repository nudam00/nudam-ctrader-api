package api

import (
	"encoding/json"
	"fmt"
	"nudam-ctrader-api/external/mongodb"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/logger"
	"nudam-ctrader-api/strategy"
	"nudam-ctrader-api/types/ctrader"
	"nudam-ctrader-api/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// Read messages from websocket.
func (api *CTrader) ReadMessage() error {
	resp, err := utils.ReadMsg(api.ws)
	if err != nil {
		return err
	}

	var baseMsg ctrader.Message[json.RawMessage]
	if err = json.Unmarshal(resp, &baseMsg); err != nil {
		return err
	}

	switch baseMsg.PayloadType {
	case configs_helper.TraderConfiguration.PayloadTypes["protooaspotevent"]:
		err := saveProtoOASpotEvent(baseMsg)
		if err != nil {
			return err
		}
	case configs_helper.TraderConfiguration.PayloadTypes["protooagettrendbarsres"]:
		err := saveProtoOAGetTrendbarsRes(baseMsg)
		if err != nil {
			return err
		}
	case configs_helper.TraderConfiguration.PayloadTypes["protooatraderres"]:
		balance, err := saveProtoOATraderRes(baseMsg)
		if err != nil {
			return err
		}
		if api.onBalanceUpdate != nil {
			api.onBalanceUpdate(balance)
		}
	case configs_helper.TraderConfiguration.PayloadTypes["protooasubscribespotsres"]:
		logger.LogMessage("spots subscribed successfully...")
	case configs_helper.TraderConfiguration.PayloadTypes["hearbeatevent"]:
		break
	default:
		return fmt.Errorf("unknown payloadType: %d", baseMsg.PayloadType)
	}

	return nil
}

// Take available balance.
func saveProtoOATraderRes(baseMsg ctrader.Message[json.RawMessage]) (int64, error) {
	var protoOATraderRes ctrader.ProtoOATraderRes
	if err := json.Unmarshal(baseMsg.Payload, &protoOATraderRes); err != nil {
		return -1, err
	}

	return protoOATraderRes.Trader.Balance, nil
}

// Update bid and ask in mongodb.
func saveProtoOASpotEvent(baseMsg ctrader.Message[json.RawMessage]) error {
	var protoOASpotEvent ctrader.ProtoOASpotEvent
	if err := json.Unmarshal(baseMsg.Payload, &protoOASpotEvent); err != nil {
		return err
	}

	if protoOASpotEvent.Ask != nil && protoOASpotEvent.Bid != nil {
		filter := bson.M{"symbolId": protoOASpotEvent.SymbolId}
		update := bson.M{
			"$set": bson.M{
				"prices.bid": protoOASpotEvent.Bid,
				"prices.ask": protoOASpotEvent.Ask,
			},
		}
		if err := mongodb.UpdateMongo(filter, update); err != nil {
			return err
		}
	}

	return nil
}

// Update close prices in mongodb.
func saveProtoOAGetTrendbarsRes(baseMsg ctrader.Message[json.RawMessage]) error {
	var protoOAGetTrendbarsRes ctrader.ProtoOAGetTrendbarsRes
	if err := json.Unmarshal(baseMsg.Payload, &protoOAGetTrendbarsRes); err != nil {
		return err
	}

	var closePrices []float64
	for _, bar := range protoOAGetTrendbarsRes.Trendbar {
		closePrice := float64(bar.Low + int64(bar.DeltaClose))
		closePrices = append(closePrices, closePrice)
	}

	emas, err := strategy.GetEMAs(closePrices)
	if err != nil {
		return fmt.Errorf("%v for period %d and symbolId %d", err, protoOAGetTrendbarsRes.Period, protoOAGetTrendbarsRes.SymbolId)
	}

	filter := bson.M{"symbolId": protoOAGetTrendbarsRes.SymbolId, "ema": bson.M{
		"$elemMatch": bson.M{"period": protoOAGetTrendbarsRes.Period},
	}}
	update := bson.M{
		"$set": bson.M{
			"ema.$.values": emas,
		},
	}

	if err := mongodb.UpdateMongo(filter, update); err != nil {
		return err
	}

	return nil
}

// Callback for balance update.
func (api *CTrader) SetOnBalanceUpdate(handler func(int64)) {
	api.onBalanceUpdate = handler
}
