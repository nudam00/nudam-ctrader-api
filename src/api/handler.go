package api

import (
	"encoding/json"
	"fmt"
	"nudam-ctrader-api/external/mongodb"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/logger"
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
	case configs_helper.TraderConfiguration.PayloadTypes["protooasubscribespotsres"]:
		logger.LogMessage("spots subscribed successfully...")
	case configs_helper.TraderConfiguration.PayloadTypes["hearbeatevent"]:
		break
	default:
		return fmt.Errorf("unknown payloadType: %d", baseMsg.PayloadType)
	}

	return nil

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

	logger.LogMessage(fmt.Sprintf("bid and ask received for: %v", protoOASpotEvent.SymbolId))

	return nil
}

// Update close prices in mongodb.
func saveProtoOAGetTrendbarsRes(baseMsg ctrader.Message[json.RawMessage]) error {
	var protoOAGetTrendbarsRes ctrader.ProtoOAGetTrendbarsRes
	if err := json.Unmarshal(baseMsg.Payload, &protoOAGetTrendbarsRes); err != nil {
		return err
	}
	for _, bar := range protoOAGetTrendbarsRes.Trendbar {
		closePrice := bar.Low + int64(bar.DeltaClose)
		protoOAGetTrendbarsRes.ClosePrices = append(protoOAGetTrendbarsRes.ClosePrices, float64(closePrice))
	}

	protoOAGetTrendbarsRes.Trendbar = nil

	filter := bson.M{"symbolId": protoOAGetTrendbarsRes.SymbolId}
	update := bson.M{"$set": bson.M{"closePrices": protoOAGetTrendbarsRes.ClosePrices}}

	if err := mongodb.UpdateMongo(filter, update); err != nil {
		return err
	}

	logger.LogMessage(fmt.Sprintf("close prices received for: %v", protoOAGetTrendbarsRes.SymbolId))

	return nil
}
