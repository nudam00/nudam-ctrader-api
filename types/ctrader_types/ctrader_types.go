package ctrader_types

import "nudam-ctrader-api/types/assets"

var PayloadTypes = map[string]int{
	"ProtoOAApplicationAuthReq": 2100,
	"ProtoOAApplicationAuthRes": 2101,
	"ProtoOAAccountAuthReq":     2102,
	"ProtoOAAccountAuthRes":     2103,
	"ProtoOASymbolsListReq":     2114,
	"ProtoOASymbolsListRes":     2115,
}

type Payload interface{}

type Message[T Payload] struct {
	ClientMsgID string `json:"clientMsgId"`
	PayloadType int    `json:"payloadType"`
	Payload     T      `json:"payload"`
}

type ProtoOAApplicationAuthReq struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type ProtoOAAccountAuthReq struct {
	CtidTraderAccountId int64  `json:"ctidTraderAccountId"`
	AccessToken         string `json:"accessToken"`
}

type ProtoOASymbolsListReq struct {
	CtidTraderAccountId    int64 `json:"ctidTraderAccountId"`
	IncludeArchivedSymbols bool  `json:"includeArchivedSymbols"`
}

type ProtoOASymbolsListRes struct {
	CtidTraderAccountId int64           `json:"ctidTraderAccountId"`
	Symbol              []assets.Symbol `json:"symbol"`
}
