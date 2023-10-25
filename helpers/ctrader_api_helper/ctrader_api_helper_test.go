package ctrader_api_helper

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func TestGetClientMsgID(t *testing.T) {
	uuid := GetClientMsgID()
	assert.NotEmpty(t, uuid)
}

func TestLogError(t *testing.T) {
	LogError(errors.New("test"))
}

func setupTestServer(t *testing.T) (*websocket.Conn, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		assert.Nil(t, err)
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				break
			}
		}
	}))

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.Nil(t, err)

	return clientConn, server
}

func TestSendAndReadMsg(t *testing.T) {
	clientConn, server := setupTestServer(t)
	defer server.Close()
	defer clientConn.Close()

	msg := map[string]string{"test": "message"}
	expectedResp := "message"

	err := SendMsg(clientConn, msg)
	assert.Nil(t, err)

	resp, err := ReadMsg(clientConn)
	assert.Nil(t, err)
	assert.Contains(t, string(resp), expectedResp)
}

func TestCheckResponse(t *testing.T) {
	resp := []byte(`{"test":"123"}`)
	expected := 123
	err := CheckResponse(resp, expected)
	assert.Nil(t, err)

	expected = 0
	err = CheckResponse(resp, expected)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), strconv.Itoa(expected))
}
