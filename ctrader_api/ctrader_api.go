package ctrader_api

import (
	"fmt"
	"log"
	"net/url"
	"nudam-trading-bot/helpers/configs_helper"
	"nudam-trading-bot/helpers/ctrader_api_helper"
	"nudam-trading-bot/types/ctrader_types"

	"github.com/gorilla/websocket"
)

type CTraderAPI struct {
	wsConn *websocket.Conn
}

func NewCTraderAPI() *CTraderAPI {
	return &CTraderAPI{}
}

func (api *CTraderAPI) InitializeWsDialer() error {
	log.Println("initializes ws dialer...")

	var err error
	wsDialer := &websocket.Dialer{}
	wsURL := url.URL{
		Scheme: "wss",
		Host:   fmt.Sprintf("%s:%d", configs_helper.CTraderConfig.Host, configs_helper.CTraderConfig.Port),
	}

	api.wsConn, _, err = wsDialer.Dial(wsURL.String(), nil)
	if err != nil {
		ctrader_api_helper.LogError(err)
		return err
	}

	payloadType := 2100
	protoOAApplicationAuthReq := ctrader_types.Message[ctrader_types.ProtoOAApplicationAuthReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: payloadType,
		Payload: ctrader_types.ProtoOAApplicationAuthReq{
			ClientId:     configs_helper.CTraderAccountConfig.ClientId,
			ClientSecret: configs_helper.CTraderAccountConfig.ClientSecret,
		},
	}

	if err = ctrader_api_helper.SendMsg(api.wsConn, protoOAApplicationAuthReq); err != nil {
		return err
	}

	resp, err := ctrader_api_helper.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}

	if err = ctrader_api_helper.CheckResponse(resp, payloadType+1); err != nil {
		return err
	}

	payloadType = 2102
	protoOAAccountAuthReq := ctrader_types.Message[ctrader_types.ProtoOAAccountAuthReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: payloadType,
		Payload: ctrader_types.ProtoOAAccountAuthReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			AccessToken:         configs_helper.CTraderAccountConfig.AccessToken,
		},
	}

	if err = ctrader_api_helper.SendMsg(api.wsConn, protoOAAccountAuthReq); err != nil {
		return err
	}

	resp, err = ctrader_api_helper.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}

	if err = ctrader_api_helper.CheckResponse(resp, payloadType+1); err != nil {
		return err
	}

	return nil
}
