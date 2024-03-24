package api

import (
	"encoding/json"
	"fmt"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/logger"
	"nudam-ctrader-api/types/ctrader"
	"nudam-ctrader-api/utils"

	"github.com/gorilla/websocket"
)

type CTraderAPI interface {
	GetTrendbars(symbol string) error
	ReadMessage() error
	// SendMsgNewOrder(symbol string, orderType, tradeSide, volume, stopLoss int64) ([]byte, error)
	// SendMsgGetBalance() (float64, error)
	Open() error
	Close() error
}

type CTrader struct {
	ws          *websocket.Conn
	sendChannel chan []byte
}

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
		var protoOASpotEvent ctrader.ProtoOASpotEvent
		if err = json.Unmarshal(baseMsg.Payload, &protoOASpotEvent); err != nil {
			return err
		}
		fmt.Println(string(resp))
		// TODO
	case configs_helper.TraderConfiguration.PayloadTypes["protooagettrendbarsres"]:
		var protoOAGetTrendbarsRes ctrader.ProtoOAGetTrendbarsRes
		if err = json.Unmarshal(baseMsg.Payload, &protoOAGetTrendbarsRes); err != nil {
			return err
		}
		fmt.Println(protoOAGetTrendbarsRes.SymbolId)
		//TODO
	case configs_helper.TraderConfiguration.PayloadTypes["protooasubscribespotsres"]:
		logger.LogMessage("spots subscribed successfully...")
	case configs_helper.TraderConfiguration.PayloadTypes["hearbeatevent"]:
		break
	default:
		return fmt.Errorf("unknown payloadType: %d", baseMsg.PayloadType)
	}

	return nil

	// var closePrices []float64
	// for _, bar := range protoOAGetTrendbarsRes.Payload.Trendbar {
	// 	closePrice := bar.Low + int64(bar.DeltaClose)
	// 	closePrices = append(closePrices, float64(closePrice))
	// } // trendbars
}

func (api *CTrader) writePump() {
	for message := range api.sendChannel {
		if err := api.ws.WriteMessage(websocket.TextMessage, message); err != nil {
			logger.LogError(err, "error sending message")
			return
		}
	}
}

func (api *CTrader) sendMessage(message []byte) {
	api.sendChannel <- message
}

func (api *CTrader) Close() error {
	return api.ws.Close()
}

func NewApi() CTraderAPI {
	cTraderApi := new(CTrader)
	cTraderApi.sendChannel = make(chan []byte, 100)

	go cTraderApi.writePump()

	return cTraderApi
}
