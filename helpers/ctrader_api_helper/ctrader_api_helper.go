package ctrader_api_helper

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func GetClientMsgID() string {
	return uuid.New().String()
}

func LogError(err error) {
	if err != nil {
		log.Printf("Error: " + err.Error())
	}
}

func SendMsg(wsConn *websocket.Conn, msg interface{}) error {
	if err := wsConn.WriteJSON(msg); err != nil {
		LogError(err)
		return err
	}
	return nil
}

func ReadMsg(wsConn *websocket.Conn) ([]byte, error) {
	_, resp, err := wsConn.ReadMessage()
	if err != nil {
		LogError(err)
		return nil, err
	}
	log.Printf(string(resp))
	return resp, nil
}

func CheckResponse(resp []byte, expected int) error {
	if !strings.Contains(string(resp), strconv.Itoa(expected)) {
		err := fmt.Errorf("error receiving response from %s", strconv.Itoa(expected))
		LogError(err)
		return err
	}
	return nil
}
