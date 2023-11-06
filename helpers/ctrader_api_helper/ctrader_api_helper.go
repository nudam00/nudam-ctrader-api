package ctrader_api_helper

import (
	"fmt"
	"log"
	"nudam-ctrader-api/types/assets"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Generates uuid message id.
func GetClientMsgID() string {
	return uuid.New().String()
}

// Logs error.
func LogError(err error) {
	if err != nil {
		log.Printf("Error: " + err.Error())
	}
}

// Sends message to api with body.
func SendMsg(wsConn *websocket.Conn, msg interface{}) error {
	if err := wsConn.WriteJSON(msg); err != nil {
		LogError(err)
		return err
	}
	return nil
}

// Reads message response.
func ReadMsg(wsConn *websocket.Conn) ([]byte, error) {
	_, resp, err := wsConn.ReadMessage()
	if err != nil {
		LogError(err)
		return nil, err
	}
	log.Printf(string(resp))
	return resp, nil
}

// Checks response from message.
func CheckResponse(resp []byte, expected int) error {
	if !strings.Contains(string(resp), strconv.Itoa(expected)) {
		err := fmt.Errorf("error receiving response from %s", strconv.Itoa(expected))
		LogError(err)
		return err
	}
	return nil
}

// Finds symbol id based on given name.
func FindSymbolId(symbolName string) (int64, error) {
	for _, symbol := range assets.Symbols {
		if symbol.SymbolName == symbolName {
			return symbol.SymbolId, nil
		}
	}
	return 0, fmt.Errorf("symbol %s not found", symbolName)
}
