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
	SendMsgNewOrder(symbol string, orderType, tradeSide, volume, stopLoss int64) ([]byte, error)
	SendMsgGetBalance() (float64, error)
	SaveSymbolEntity(symbol string) (*ctrader.SymbolEntity, error)
	Close() error
}

type CTrader struct {
	wsConn       *websocket.Conn
	symbolList   []ctrader.SymbolList
	SymbolEntity ctrader.SymbolEntity
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
	if err = json.Unmarshal(resp, &protoOASymbolsListRes); err != nil {
		return err
	}
	api.symbolList = protoOASymbolsListRes.Payload.Symbol

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
	symbolId, err := utils.FindSymbolId(symbol, api.symbolList)
	if err != nil {
		return nil, err
	}
	count := utils.CalculateCountBars(period)

	resp, err := api.sendMsgTrendbars(fromTimestamp, toTimestamp, periodId, symbolId, count)
	if err != nil {
		return nil, err
	}

	var protoOAGetTrendbarsRes ctrader.Message[ctrader.ProtoOAGetTrendbarsRes]
	if err = json.Unmarshal(resp, &protoOAGetTrendbarsRes); err != nil {
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
func (api *CTrader) sendMsgTrendbars(fromTimestamp, toTimestamp, periodId, symbolId int64, count uint32) ([]byte, error) {
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
	symbolId, err := utils.FindSymbolId(symbol, api.symbolList)
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
	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaspotevent"], err); err != nil {
		return err
	}

	resp, err = utils.ReadMsg(api.wsConn)
	if err != nil {
		return err
	}
	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooasubscribespotsres"], err); err != nil {
		return err
	}

	return nil
}

// Save symbol entity.
func (api *CTrader) SaveSymbolEntity(symbol string) (*ctrader.SymbolEntity, error) {
	symbolId, err := utils.FindSymbolId(symbol, api.symbolList)
	if err != nil {
		return nil, err
	}

	protoOASymbolByIdReq := ctrader.Message[ctrader.ProtoOASymbolByIdReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooasymbolbyddreq"],
		Payload: ctrader.ProtoOASymbolByIdReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			SymbolId:            symbolId,
		},
	}

	if err := utils.SendMsg(api.wsConn, protoOASymbolByIdReq); err != nil {
		return nil, err
	}

	resp, err := utils.ReadMsg(api.wsConn)
	if err != nil {
		return nil, err
	}
	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooasymbolbyddres"], err); err != nil {
		return nil, err
	}

	var protoOASymbolByIdRes *ctrader.Message[ctrader.ProtoOASymbolByIdRes]
	if err = json.Unmarshal(resp, &protoOASymbolByIdRes); err != nil {
		return nil, err
	}

	return &protoOASymbolByIdRes.Payload.Symbol[0], nil

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
	if err = json.Unmarshal(resp, &protoOASpotEvent); err != nil {
		return nil, err
	}

	return protoOASpotEvent, nil
}

// Sends message to create new order.
func (api *CTrader) SendMsgNewOrder(symbol string, orderType, tradeSide, volume, stopLoss int64) ([]byte, error) {
	symbolId, err := utils.FindSymbolId(symbol, api.symbolList)
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
			RelativeStopLoss:    stopLoss,
			TrailingStopLoss:    true,
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

// Gets current balance.
func (api *CTrader) SendMsgGetBalance() (float64, error) {
	protoOATraderReq := ctrader.Message[ctrader.ProtoOATraderReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooatraderreq"],
		Payload: ctrader.ProtoOATraderReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
		},
	}

	if err := utils.SendMsg(api.wsConn, protoOATraderReq); err != nil {
		return 0, err
	}

	resp, err := utils.ReadMsg(api.wsConn)
	if err != nil {
		return 0, err
	}
	if err = utils.CheckResponse(resp, configs_helper.TraderConfiguration.PayloadTypes["protooatraderres"], err); err != nil {
		return 0, err
	}

	var protoOATraderRes *ctrader.Message[ctrader.ProtoOATraderRes]
	if err = json.Unmarshal(resp, &protoOATraderRes); err != nil {
		return 0, err
	}

	return float64(protoOATraderRes.Payload.Trader.Balance) / 100.0, nil
}

func (api *CTrader) Close() error {
	return api.wsConn.Close()
}

func NewApi() (CTraderAPI, error) {
	var err error

	cTraderApi := new(CTrader)
	if err = cTraderApi.initalizeCTrader(); err != nil {
		return nil, err
	}

	return cTraderApi, nil
}
