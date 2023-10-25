package ctrader_types

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
