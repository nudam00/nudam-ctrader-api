package ctrader

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
	CtidTraderAccountId int64    `json:"ctidTraderAccountId"`
	Symbol              []Symbol `json:"symbol"`
}

type Symbol struct {
	SymbolId         int64   `json:"symbolId"`
	SymbolName       string  `json:"symbolName"`
	Enabled          bool    `json:"enabled"`
	BaseAssetId      int64   `json:"baseAssetId"`
	QuoteAssetId     int64   `json:"quoteAssetId"`
	SymbolCategoryId int64   `json:"symbolCategoryId"`
	Description      string  `json:"description"`
	SortingNumber    float64 `json:"sortingNumber"`
}

type ProtoOAGetTrendbarsReq struct {
	CtidTraderAccountId int64  `json:"ctidTraderAccountId"`
	FromTimestamp       int64  `json:"fromTimestamp"`
	ToTimestamp         int64  `json:"toTimestamp"`
	Period              int    `json:"period"`
	SymbolId            int64  `json:"symbolId"`
	Count               uint32 `json:"count"`
}

type ProtoOAGetTrendbarsRes struct {
	CtidTraderAccountId int64      `json:"ctidTraderAccountId"`
	Period              int        `json:"period"`
	Timestamp           int64      `json:"timestamp"`
	Trendbar            []Trendbar `json:"trendbar"`
	SymbolId            int64      `json:"symbolId"`
}

type Trendbar struct {
	Volume                int64  `json:"volume"`
	Period                int    `json:"period"`
	Low                   int64  `json:"low"`
	DeltaOpen             uint64 `json:"deltaOpen"`
	DeltaClose            uint64 `json:"deltaClose"`
	DeltaHigh             uint64 `json:"deltaHigh"`
	UTCTimestampInMinutes uint32 `json:"utcTimestampInMinutes"`
}
