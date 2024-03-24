package utils

import (
	"fmt"
	"nudam-ctrader-api/logger"
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

// Sends message to api with body.
func SendMsg(wsConn *websocket.Conn, msg interface{}) error {
	if err := wsConn.WriteJSON(msg); err != nil {
		logger.LogError(err, fmt.Sprintln(msg))
		return err
	}
	return nil
}

// Reads message response.
func ReadMsg(wsConn *websocket.Conn) ([]byte, error) {
	_, resp, err := wsConn.ReadMessage()
	if err != nil {
		logger.LogError(err, string(resp))
		return nil, err
	}
	return resp, nil
}

// Checks response from message.
func CheckResponse(resp []byte, expected int, err error) error {
	if !strings.Contains(string(resp), strconv.Itoa(expected)) {
		err := fmt.Errorf("error receiving response from %s; error: %s", strconv.Itoa(expected), err.Error())
		logger.LogError(err, string(resp))
		return err
	}
	return nil
}

// Calculates fromTimestamp and toTimestamp.
func CalculateTimestamps(numberDays int) (int64, int64) {
	now := time.Now()
	fromTime := now.AddDate(0, 0, -numberDays)
	fromTimestamp := fromTime.UnixNano() / int64(time.Millisecond)
	toTimestamp := now.UnixNano() / int64(time.Millisecond)
	return fromTimestamp, toTimestamp
}

// // Calculates the amount of bars based on given period.
// func CalculateCountBars(period string) uint32 {
// 	return configs_helper.TraderConfiguration.Periods[period].CountBars * configs_helper.TraderConfiguration.Periods[period].NumberDays
// }
