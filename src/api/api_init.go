package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nudam-ctrader-api/external/mongodb"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/logger"
	"nudam-ctrader-api/types/ctrader"
	"nudam-ctrader-api/utils"

	"github.com/gorilla/websocket"
)

// Initialize cTrader connection with available symbols.
func (api *CTrader) Open() error {
	logger.LogMessage("initializes ctrader connection...")

	if err := api.initializeWsDialer(); err != nil {
		return err
	}

	if err := api.authenticate(); err != nil {
		return err
	}

	if err := api.saveAvailableSymbols(); err != nil {
		return err
	}

	if err := api.saveSymbolEntity(); err != nil {
		return err
	}

	if err := api.sendMsgSubscribeSpot(); err != nil {
		return err
	}

	return nil
}

// Initialize websocket connection.
func (api *CTrader) initializeWsDialer() error {
	var err error
	var resp *http.Response
	wsDialer := &websocket.Dialer{}
	wsURL := url.URL{
		Scheme: "wss",
		Host:   fmt.Sprintf("%s:%d", configs_helper.CTraderConfig.Host, configs_helper.CTraderConfig.Port),
	}

	api.ws, resp, _ = wsDialer.Dial(wsURL.String(), nil)
	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(err, string(respB))
		return err
	}

	logger.LogMessage("ws dialer initalized successfully...")

	return nil
}

// Initialize cTrader account.
func (api *CTrader) authenticate() error {
	protoOAApplicationAuthReq := ctrader.Message[ctrader.ProtoOAApplicationAuthReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooaapplicationauthreq"],
		Payload: ctrader.ProtoOAApplicationAuthReq{
			ClientId:     configs_helper.CTraderAccountConfig.ClientId,
			ClientSecret: configs_helper.CTraderAccountConfig.ClientSecret,
		},
	}

	if err := utils.SendMsg(api.ws, protoOAApplicationAuthReq); err != nil {
		return err
	}
	resp, err := utils.ReadMsg(api.ws)
	if err != nil {
		return err
	}
	if err = utils.CheckResponseContains(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaapplicationauthres"], err); err != nil {
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

	if err = utils.SendMsg(api.ws, protoOAAccountAuthReq); err != nil {
		return err
	}
	resp, err = utils.ReadMsg(api.ws)
	if err != nil {
		return err
	}
	if err = utils.CheckResponseContains(resp, configs_helper.TraderConfiguration.PayloadTypes["protooaaccountauthres"], err); err != nil {
		return err
	}

	logger.LogMessage("cTrader account initalized successfully...")

	return nil
}

// Save available symbols to MongoDb.
func (api *CTrader) saveAvailableSymbols() error {
	logger.LogMessage("getting available symbols...")

	protoOASymbolsListReq := ctrader.Message[ctrader.ProtoOASymbolsListReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooasymbolslistreq"],
		Payload: ctrader.ProtoOASymbolsListReq{
			CtidTraderAccountId:    configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			IncludeArchivedSymbols: false,
		},
	}

	if err := utils.SendMsg(api.ws, protoOASymbolsListReq); err != nil {
		return err
	}
	resp, err := utils.ReadMsg(api.ws)
	if err != nil {
		return err
	}
	if err = utils.CheckResponseContains(resp, configs_helper.TraderConfiguration.PayloadTypes["protooasymbolslistres"], err); err != nil {
		return err
	}

	var protoOASymbolsListRes ctrader.Message[ctrader.ProtoOASymbolsListRes]
	if err = json.Unmarshal(resp, &protoOASymbolsListRes); err != nil {
		return err
	}

	if err = mongodb.SaveToMongo(protoOASymbolsListRes, configs_helper.TraderConfiguration.PayloadTypes["protooasymbolslistres"]); err != nil {
		return err
	}

	logger.LogMessage("available symbols saved successfully...")

	return nil
}

// Save symbol entities.
func (api *CTrader) saveSymbolEntity() error {
	symbolIds, err := mongodb.FindSymbolIds(configs_helper.TraderConfiguration.CurrencyPairs, configs_helper.TraderConfiguration.PayloadTypes["protooasymbolslistres"])
	if err != nil {
		return err
	}

	protoOASymbolByIdReq := ctrader.Message[ctrader.ProtoOASymbolByIdReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooasymbolbyddreq"],
		Payload: ctrader.ProtoOASymbolByIdReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			SymbolId:            symbolIds,
		},
	}

	if err := utils.SendMsg(api.ws, protoOASymbolByIdReq); err != nil {
		return err
	}
	resp, err := utils.ReadMsg(api.ws)
	if err != nil {
		return err
	}
	if err = utils.CheckResponseContains(resp, configs_helper.TraderConfiguration.PayloadTypes["protooasymbolbyddres"], err); err != nil {
		return err
	}

	var protoOASymbolByIdRes ctrader.Message[ctrader.ProtoOASymbolByIdRes]
	if err = json.Unmarshal(resp, &protoOASymbolByIdRes); err != nil {
		return err
	}

	if err = mongodb.SaveToMongo(protoOASymbolByIdRes, configs_helper.TraderConfiguration.PayloadTypes["protooasymbolbyddres"]); err != nil {
		return err
	}

	logger.LogMessage("available symbol entities saved successfully...")

	return nil
}

// Subscribe spots to get current price.
func (api *CTrader) sendMsgSubscribeSpot() error {
	symbolIds, err := mongodb.FindSymbolIds(configs_helper.TraderConfiguration.CurrencyPairs, configs_helper.TraderConfiguration.PayloadTypes["protooasymbolslistres"])
	if err != nil {
		return err
	}

	protoOASubscribeSpotsReq := ctrader.Message[ctrader.ProtoOASubscribeSpotsReq]{
		ClientMsgID: utils.GetClientMsgID(),
		PayloadType: configs_helper.TraderConfiguration.PayloadTypes["protooasubscribespotsreq"],
		Payload: ctrader.ProtoOASubscribeSpotsReq{
			CtidTraderAccountId: configs_helper.CTraderAccountConfig.CtidTraderAccountId,
			SymbolId:            symbolIds,
		},
	}

	if err := utils.SendMsg(api.ws, protoOASubscribeSpotsReq); err != nil {
		return err
	}
	_, err = utils.ReadMsg(api.ws)
	if err != nil {
		return err
	}

	return nil
}
