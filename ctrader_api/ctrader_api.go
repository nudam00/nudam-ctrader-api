package ctrader_api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/helpers/ctrader_api_helper"
	"nudam-ctrader-api/strategy"
	"nudam-ctrader-api/types/constants"
	"nudam-ctrader-api/types/ctrader"

	"github.com/gorilla/websocket"
)

type CTraderAPI struct {
	wsConn     *websocket.Conn
	symbols    []ctrader.Symbol
	numberDays int
	period     string
	symbol     string
}

// Initalizes new cTrader api.
func NewCTraderAPI(numberDays int, period string, symbol string) *CTraderAPI {
	return &CTraderAPI{numberDays: numberDays, period: period, symbol: symbol}
}

// Initialize cTrader connection with available symbols.
func (api *CTraderAPI) InitalizeCTrader() error {
	ctrader_api_helper.LogMessage("initializes ctrader connection...")

	if err := api.initializeWsDialer(); err != nil {
		return err
	}

	if err := api.authenticate(); err != nil {
		return err
	}

	if err := api.saveAvailableSymbols(); err != nil {
		return err
	}

	return nil
}

// Initializes websocket connection.
func (api *CTraderAPI) initializeWsDialer() error {
	ctrader_api_helper.LogMessage("initializes ws dialer...")

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

	ctrader_api_helper.LogMessage("ws dialer initalized successfully...")

	return nil
}

// Initializes cTrader account.
func (api *CTraderAPI) authenticate() error {
	protoOAApplicationAuthReq := ctrader.Message[ctrader.ProtoOAApplicationAuthReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: constants.PayloadTypes["ProtoOAApplicationAuthReq"],
		Payload: ctrader.ProtoOAApplicationAuthReq{
			ClientId:     configs_helper.CTraderAccountConfig.ClientId,
			ClientSecret: configs_helper.CTraderAccountConfig.ClientSecret,
		},
	}

	if err := ctrader_api_helper.SendMsg(api.wsConn, protoOAApplicationAuthReq); err != nil {
		return err
	}
	resp, err := ctrader_api_helper.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}
	if err = ctrader_api_helper.CheckResponse(resp, constants.PayloadTypes["ProtoOAApplicationAuthRes"]); err != nil {
		return err
	}

	protoOAAccountAuthReq := ctrader.Message[ctrader.ProtoOAAccountAuthReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: constants.PayloadTypes["ProtoOAAccountAuthReq"],
		Payload: ctrader.ProtoOAAccountAuthReq{
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
	if err = ctrader_api_helper.CheckResponse(resp, constants.PayloadTypes["ProtoOAAccountAuthRes"]); err != nil {
		return err
	}

	ctrader_api_helper.LogMessage("cTrader account initalized successfully...")

	return nil
}

// Saves available symbols to variable in assets.go.
func (api *CTraderAPI) saveAvailableSymbols() error {
	ctrader_api_helper.LogMessage("getting available symbols...")

	resp, err := api.getAvailableSymbols()
	if err != nil {
		return err
	}

	var protoOASymbolsListRes ctrader.Message[ctrader.ProtoOASymbolsListRes]
	err = json.Unmarshal(resp, &protoOASymbolsListRes)
	if err != nil {
		return err
	}
	api.symbols = protoOASymbolsListRes.Payload.Symbol

	ctrader_api_helper.LogMessage("available symbols saved successfully...")

	return nil
}

// Sends message to receive available symbols.
func (api *CTraderAPI) getAvailableSymbols() ([]byte, error) {
	protoOASymbolsListReq := ctrader.Message[ctrader.ProtoOASymbolsListReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: constants.PayloadTypes["ProtoOASymbolsListReq"],
		Payload: ctrader.ProtoOASymbolsListReq{
			CtidTraderAccountId:    configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			IncludeArchivedSymbols: false,
		},
	}

	if err := ctrader_api_helper.SendMsg(api.wsConn, protoOASymbolsListReq); err != nil {
		return nil, err
	}
	resp, err := ctrader_api_helper.ReadMsg(api.wsConn)
	if err != nil {
		return nil, err
	}
	if err = ctrader_api_helper.CheckResponse(resp, constants.PayloadTypes["ProtoOASymbolsListRes"]); err != nil {
		return nil, err
	}

	return resp, nil
}

// Get trendbars based on given symbol.
func (api *CTraderAPI) GetTrendbars() error {
	ctrader_api_helper.LogMessage("getting trendbars...")

	fromTimestamp, toTimestamp := ctrader_api_helper.CalculateTimestamps(api.numberDays)
	periodId := constants.Periods[api.period]
	symbolId, err := ctrader_api_helper.FindSymbolId(api.symbol, api.symbols)
	if err != nil {
		return err
	}
	count := ctrader_api_helper.CalculateCountBars(api.period, api.numberDays)

	resp, err := api.sendMsgTrendbars(fromTimestamp, toTimestamp, periodId, symbolId, count)
	if err != nil {
		return err
	}

	var protoOAGetTrendbarsRes ctrader.Message[ctrader.ProtoOAGetTrendbarsRes]
	err = json.Unmarshal(resp, &protoOAGetTrendbarsRes)
	if err != nil {
		return err
	}

	var closePrices []float64
	for _, bar := range protoOAGetTrendbarsRes.Payload.Trendbar {
		closePrice := bar.Low + int64(bar.DeltaClose)
		closePrices = append(closePrices, float64(closePrice))
	}

	strategy.GetSignal(closePrices)

	return nil
}

// Sends message to get current trendbars.
func (api *CTraderAPI) sendMsgTrendbars(fromTimestamp int64, toTimestamp int64, periodId int, symbolId int64, count uint32) ([]byte, error) {
	protoOAGetTrendbarsReq := ctrader.Message[ctrader.ProtoOAGetTrendbarsReq]{
		ClientMsgID: ctrader_api_helper.GetClientMsgID(),
		PayloadType: constants.PayloadTypes["ProtoOAGetTrendbarsReq"],
		Payload: ctrader.ProtoOAGetTrendbarsReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			FromTimestamp:       fromTimestamp,
			ToTimestamp:         toTimestamp,
			Period:              periodId,
			SymbolId:            symbolId,
			Count:               count,
		},
	}

	if err := ctrader_api_helper.SendMsg(api.wsConn, protoOAGetTrendbarsReq); err != nil {
		return nil, err
	}
	resp, err := ctrader_api_helper.ReadMsg(api.wsConn)
	if err != nil {
		return nil, err
	}
	if err = ctrader_api_helper.CheckResponse(resp, constants.PayloadTypes["ProtoOAGetTrendbarsRes"]); err != nil {
		return nil, err
	}

	return resp, nil
}
