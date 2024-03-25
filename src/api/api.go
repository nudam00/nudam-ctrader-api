package api

import (
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

func NewApi() CTraderAPI {
	cTraderApi := new(CTrader)
	cTraderApi.sendChannel = make(chan []byte, 100)

	go cTraderApi.writePump()

	return cTraderApi
}
