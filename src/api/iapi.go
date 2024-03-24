package api

import (
	"github.com/gorilla/websocket"
)

type CTraderAPI interface {
	GetTrendbars(symbol string) error
	// SendMsgReadMessage() (*ctrader.Message[ctrader.ProtoOASpotEvent], error)
	// SendMsgNewOrder(symbol string, orderType, tradeSide, volume, stopLoss int64) ([]byte, error)
	// SendMsgGetBalance() (float64, error)
	Close() error
}

type CTrader struct {
	ws *websocket.Conn
}

func (api *CTrader) Close() error {
	return api.ws.Close()
}

func NewApi() (CTraderAPI, error) {
	var err error

	cTraderApi := new(CTrader)
	if err = cTraderApi.initalizeCTrader(); err != nil {
		return nil, err
	}

	return cTraderApi, nil
}
