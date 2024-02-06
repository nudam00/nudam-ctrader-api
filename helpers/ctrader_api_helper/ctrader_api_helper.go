package ctrader_api_helper

import (
	"fmt"
	"log"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/types/ctrader"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Generates uuid message id.
func GetClientMsgID() string {
	return uuid.New().String()
}

// Logs error.
func LogError(err error, msg string) {
	if err != nil {
		log.Printf("Msg: " + msg)
		log.Printf("Error: " + err.Error())
	}
}

// Logs message.
func LogMessage(msg string) {
	log.Printf(msg)
}

// Sends message to api with body.
func SendMsg(wsConn *websocket.Conn, msg interface{}) error {
	if err := wsConn.WriteJSON(msg); err != nil {
		LogError(err, fmt.Sprintln(msg))
		return err
	}
	return nil
}

// Reads message response.
func ReadMsg(wsConn *websocket.Conn) ([]byte, error) {
	_, resp, err := wsConn.ReadMessage()
	if err != nil {
		LogError(err, string(resp))
		return nil, err
	}
	return resp, nil
}

// Checks response from message.
func CheckResponse(resp []byte, expected int) error {
	if !strings.Contains(string(resp), strconv.Itoa(expected)) {
		err := fmt.Errorf("error receiving response from %s", strconv.Itoa(expected))
		LogError(err, string(resp))
		return err
	}
	return nil
}

// Finds symbol id based on given name.
func FindSymbolId(symbolName string, symbols []ctrader.Symbol) (int64, error) {
	for _, symbol := range symbols {
		if symbol.SymbolName == symbolName {
			return symbol.SymbolId, nil
		}
	}
	return 0, fmt.Errorf("symbol %s not found", symbolName)
}

// Calculates fromTimestamp and toTimestamp.
func CalculateTimestamps(numberDays int) (int64, int64) {
	now := time.Now()
	fromTime := now.AddDate(0, 0, -numberDays)
	fromTimestamp := fromTime.UnixNano() / int64(time.Millisecond)
	toTimestamp := now.UnixNano() / int64(time.Millisecond)
	return fromTimestamp, toTimestamp
}

// Calculates the amount of bars based on given period.
func CalculateCountBars(period string) uint32 {
	return configs_helper.TraderConfiguration.Periods[period].CountBars * configs_helper.TraderConfiguration.Periods[period].NumberDays
}
