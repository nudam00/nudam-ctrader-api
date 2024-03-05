package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/types/ctrader"
	"nudam-ctrader-api/utils"

	"github.com/gorilla/websocket"
)

type CTraderAPI interface {
	GetTrendbars(symbol, period string) ([]float64, error)
	SendMsgSubscribeSpot(symbol string) error
	SendMsgReadMessage() (*ctrader.Message[ctrader.ProtoOASpotEvent], error)
	SendMsgNewOrder(symbol string, orderType, tradeSide, volume int64, stopLoss, takeProfit *float64, clientOrderId *string, traillingStopLoss *bool) ([]byte, error)
}

type CTrader struct {
	wsConn  *websocket.Conn
	symbols []ctrader.Symbol
}

// Initialize cTrader connection with available symbols.
func (api *CTrader) initalizeCTrader() error {
	utils.LogMessage("initializes ctrader connection...")

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
func (api *CTrader) initializeWsDialer() error {
	utils.LogMessage("initializes ws dialer...")

	var err error
	var resp *http.Response
	wsDialer := &websocket.Dialer{}
	wsURL := url.URL{
		Scheme: "wss",
		Host:   fmt.Sprintf("%s:%d", configs_helper.CTraderConfig.Host, configs_helper.CTraderConfig.Port),
	}

	api.wsConn, resp, _ = wsDialer.Dial(wsURL.String(), nil)
	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err != nil {
		utils.LogError(err, string(respB))
		return err
	}

	utils.LogMessage("ws dialer initalized successfully...")

	return nil
}

// Initializes cTrader account.
func (api *CTrader) authenticate() error {
	protoOAApplicationAuthReq := ctrader.Message[ctrader.ProtoOAApplicationAuthReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooaapplicationauthreq"],
		Payload: ctrader.ProtoOAApplicationAuthReq{
			ClientId:     configs_helper.CTraderAccountConfig.ClientId,
			ClientSecret: configs_helper.CTraderAccountConfig.ClientSecret,
		},
	}

	if err := utils.SendMsg(api.wsConn, protoOAApplicationAuthReq); err != nil {
		return err
	}
	resp, err := utils.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}
	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaapplicationauthres"], err); err != nil {
		return err
	}

	protoOAAccountAuthReq := ctrader.Message[ctrader.ProtoOAAccountAuthReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooaaccountauthreq"],
		Payload: ctrader.ProtoOAAccountAuthReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			AccessToken:         configs_helper.CTraderAccountConfig.AccessToken,
		},
	}

	if err = utils.SendMsg(api.wsConn, protoOAAccountAuthReq); err != nil {
		return err
	}
	resp, err = utils.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}
	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaaccountauthres"], err); err != nil {
		return err
	}

	utils.LogMessage("cTrader account initalized successfully...")

	return nil
}

// Saves available symbols to variable in assets.go.
func (api *CTrader) saveAvailableSymbols() error {
	utils.LogMessage("getting available symbols...")

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

	utils.LogMessage("available symbols saved successfully...")

	return nil
}

// Sends message to receive available symbols.
func (api *CTrader) getAvailableSymbols() ([]byte, error) {
	protoOASymbolsListReq := ctrader.Message[ctrader.ProtoOASymbolsListReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooasymbolslistreq"],
		Payload: ctrader.ProtoOASymbolsListReq{
			CtidTraderAccountId:    configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			IncludeArchivedSymbols: false,
		},
	}

	if err := utils.SendMsg(api.wsConn, protoOASymbolsListReq); err != nil {
		return nil, err
	}
	resp, err := utils.ReadMsg(api.wsConn)
	if err != nil {
		return nil, err
	}
	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooasymbolslistres"], err); err != nil {
		return nil, err
	}

	return resp, nil
}

// Get trendbars based on given symbol.
func (api *CTrader) GetTrendbars(symbol, period string) ([]float64, error) {
	utils.LogMessage("getting trendbars...")

	fromTimestamp, toTimestamp := utils.CalculateTimestamps(int(configs_helper.TraderConfiguration.Periods[period].NumberDays))
	periodId := configs_helper.TraderConfiguration.Periods[period].Value
	symbolId, err := utils.FindSymbolId(symbol, api.symbols)
	if err != nil {
		return nil, err
	}
	count := utils.CalculateCountBars(period)

	resp, err := api.sendMsgTrendbars(fromTimestamp, toTimestamp, periodId, symbolId, count)
	if err != nil {
		return nil, err
	}

	var protoOAGetTrendbarsRes ctrader.Message[ctrader.ProtoOAGetTrendbarsRes]
	err = json.Unmarshal(resp, &protoOAGetTrendbarsRes)
	if err != nil {
		return nil, err
	}

	var closePrices []float64
	for _, bar := range protoOAGetTrendbarsRes.Payload.Trendbar {
		closePrice := bar.Low + int64(bar.DeltaClose)
		closePrices = append(closePrices, float64(closePrice))
	}

	return closePrices, nil
}

// Sends message to get current trendbars.
func (api *CTrader) sendMsgTrendbars(fromTimestamp int64, toTimestamp int64, periodId int, symbolId int64, count uint32) ([]byte, error) {
	protoOAGetTrendbarsReq := ctrader.Message[ctrader.ProtoOAGetTrendbarsReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooagettrendbarsreq"],
		Payload: ctrader.ProtoOAGetTrendbarsReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			FromTimestamp:       fromTimestamp,
			ToTimestamp:         toTimestamp,
			Period:              periodId,
			SymbolId:            symbolId,
			Count:               count,
		},
	}

	if err := utils.SendMsg(api.wsConn, protoOAGetTrendbarsReq); err != nil {
		return nil, err
	}
	resp, err := utils.ReadMsg(api.wsConn)
	if err != nil {
		return nil, err
	}
	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooagettrendbarsres"], err); err != nil {
		return nil, err
	}

	return resp, nil
}

// Subscribes spot to get current price.
func (api *CTrader) SendMsgSubscribeSpot(symbol string) error {
	symbolId, err := utils.FindSymbolId(symbol, api.symbols)
	if err != nil {
		return err
	}

	protoOASubscribeSpotsReq := ctrader.Message[ctrader.ProtoOASubscribeSpotsReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooasubscribespotsreq"],
		Payload: ctrader.ProtoOASubscribeSpotsReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			SymbolId:            symbolId,
		},
	}

	if err := utils.SendMsg(api.wsConn, protoOASubscribeSpotsReq); err != nil {
		return err
	}

	resp, err := utils.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}
	err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaspotevent"], err)
	if err != nil {
		return err
	}

	resp, err = utils.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}
	err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooasubscribespotsres"], err)
	if err != nil {
		return err
	}

	return nil
}

// Sends message to get current price.
func (api *CTrader) SendMsgReadMessage() (*ctrader.Message[ctrader.ProtoOASpotEvent], error) {
	resp, err := utils.ReadMsg(api.wsConn)
	if err != nil {
		return nil, err
	}
	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaspotevent"], err); err != nil {
		return nil, err
	}

	var protoOASpotEvent *ctrader.Message[ctrader.ProtoOASpotEvent]
	err = json.Unmarshal(resp, &protoOASpotEvent)
	if err != nil {
		return nil, err
	}

	return protoOASpotEvent, nil
}

// Sends message to create new order.
func (api *CTrader) SendMsgNewOrder(symbol string, orderType, tradeSide, volume int64, stopLoss, takeProfit *float64, clientOrderId *string, traillingStopLoss *bool) ([]byte, error) {
	symbolId, err := utils.FindSymbolId(symbol, api.symbols)
	if err != nil {
		return nil, err
	}

	protoOANewOrderReq := ctrader.Message[ctrader.ProtoOANewOrderReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooaneworderreq"],
		Payload: ctrader.ProtoOANewOrderReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			SymbolId:            symbolId,
			OrderType:           orderType,
			TradeSide:           tradeSide,
			Volume:              volume,
			StopLoss:            stopLoss,
			TakeProfit:          takeProfit,
			ClientOrderId:       clientOrderId,
			TrailingStopLoss:    traillingStopLoss,
		},
	}

	if err := utils.SendMsg(api.wsConn, protoOANewOrderReq); err != nil {
		return nil, err
	}
	resp, err := utils.ReadMsg(api.wsConn)
	if err != nil {
		return nil, err
	}

	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaexecutionevent"], err); err != nil {
		return nil, err
	}

	return resp, nil
}

func NewApi() (CTraderAPI, error) {
	var err error

	CTrader := new(CTrader)
	err = CTrader.initalizeCTrader()
	if err != nil {
		return nil, err
	}

	return CTrader, nil
}
