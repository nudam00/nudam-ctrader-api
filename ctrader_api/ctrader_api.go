package ctrader_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/helpers/ctrader_api_helper"
	"nudam-ctrader-api/types/assets"
	"nudam-ctrader-api/types/ctrader_types"
	"time"

	"github.com/gorilla/websocket"
)

type CTraderAPI struct {
	wsConn *websocket.Conn
}

// Initalizes new cTrader api.
func NewCTraderAPI() *CTraderAPI {
	return &CTraderAPI{}
}

// Initializes cTrader account using websockets.
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

	protoOAApplicationAuthReq := ctrader_types.Message[ctrader_types.ProtoOAApplicationAuthReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: ctrader_types.PayloadTypes["ProtoOAApplicationAuthReq"],
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

	if err = ctrader_api_helper.CheckResponse(resp, ctrader_types.PayloadTypes["ProtoOAApplicationAuthRes"]); err != nil {
		return err
	}

	protoOAAccountAuthReq := ctrader_types.Message[ctrader_types.ProtoOAAccountAuthReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: ctrader_types.PayloadTypes["ProtoOAAccountAuthReq"],
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

	if err = ctrader_api_helper.CheckResponse(resp, ctrader_types.PayloadTypes["ProtoOAAccountAuthRes"]); err != nil {
		return err
	}

	return nil
}

// Saves available symbols to variable in assets.go.
func (api *CTraderAPI) SaveAvailableSymbols() error {
	log.Println("getting available symbols...")
	protoOASymbolsListReq := ctrader_types.Message[ctrader_types.ProtoOASymbolsListReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: ctrader_types.PayloadTypes["ProtoOASymbolsListReq"],
		Payload: ctrader_types.ProtoOASymbolsListReq{
			CtidTraderAccountId:    configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			IncludeArchivedSymbols: false,
		},
	}

	if err := ctrader_api_helper.SendMsg(api.wsConn, protoOASymbolsListReq); err != nil {
		return err
	}

	resp, err := ctrader_api_helper.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}

	if err = ctrader_api_helper.CheckResponse(resp, ctrader_types.PayloadTypes["ProtoOASymbolsListRes"]); err != nil {
		return err
	}

	var protoOASymbolsListRes ctrader_types.Message[ctrader_types.ProtoOASymbolsListRes]
	err = json.Unmarshal(resp, &protoOASymbolsListRes)
	if err != nil {
		return err
	}
	assets.Symbols = protoOASymbolsListRes.Payload.Symbol

	return nil
}

// Get trendbars based on given symbol.
func (api *CTraderAPI) GetTrendbars(numberDays int, period string, symbol string, countBars uint32) error {
	log.Println("getting trendbars...")

	now := time.Now()
	fromTime := now.AddDate(0, 0, -numberDays)
	fromTimestamp := fromTime.UnixNano() / int64(time.Millisecond)
	toTimestamp := now.UnixNano() / int64(time.Millisecond)
	periodId := assets.Periods[period]
	symbolId, err := ctrader_api_helper.FindSymbolId(symbol)
	if err != nil {
		return err
	}
	count := countBars

	protoOAGetTrendbarsReq := ctrader_types.Message[ctrader_types.ProtoOAGetTrendbarsReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: ctrader_types.PayloadTypes["ProtoOAGetTrendbarsReq"],
		Payload: ctrader_types.ProtoOAGetTrendbarsReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			FromTimestamp:       fromTimestamp,
			ToTimestamp:         toTimestamp,
			Period:              periodId,
			SymbolId:            symbolId,
			Count:               count,
		},
	}

	if err := ctrader_api_helper.SendMsg(api.wsConn, protoOAGetTrendbarsReq); err != nil {
		return err
	}

	resp, err := ctrader_api_helper.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}

	if err = ctrader_api_helper.CheckResponse(resp, ctrader_types.PayloadTypes["ProtoOAGetTrendbarsRes"]); err != nil {
		return err
	}

	var protoOAGetTrendbarsRes ctrader_types.Message[ctrader_types.ProtoOAGetTrendbarsRes]
	err = json.Unmarshal(resp, &protoOAGetTrendbarsRes)
	if err != nil {
		return err
	}

	var closePrices []int64
	for _, bar := range protoOAGetTrendbarsRes.Payload.Trendbar {
		closePrice := bar.Low + int64(bar.DeltaClose)
		closePrices = append(closePrices, closePrice)
	}

	log.Println("Close Prices: ", closePrices)

	return nil
}
