package ctrader

type Payload interface{}

type Message[T Payload] struct {
	ClientMsgID string `json:"clientMsgId"`
	PayloadType int    `json:"payloadType"`
	Payload     T      `json:"payload"`
}

// Auth.
type ProtoOAApplicationAuthReq struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// Auth.
type ProtoOAAccountAuthReq struct {
	CtidTraderAccountId int64  `json:"ctidTraderAccountId"`
	AccessToken         string `json:"accessToken"`
}

// Get available symbols request message.
type ProtoOASymbolsListReq struct {
	CtidTraderAccountId    int64 `json:"ctidTraderAccountId"`
	IncludeArchivedSymbols bool  `json:"includeArchivedSymbols"`
}

// Get available symbols response message.
type ProtoOASymbolsListRes struct {
	CtidTraderAccountId int64        `json:"ctidTraderAccountId"`
	Symbol              []SymbolList `json:"symbol"`
}

type SymbolList struct {
	SymbolId         int64    `json:"symbolId"`
	SymbolName       *string  `json:"symbolName"`
	Enabled          *bool    `json:"enabled"`
	BaseAssetId      *int64   `json:"baseAssetId"`
	QuoteAssetId     *int64   `json:"quoteAssetId"`
	SymbolCategoryId *int64   `json:"symbolCategoryId"`
	Description      *string  `json:"description"`
	SortingNumber    *float64 `json:"sortingNumber"`
}

// Get trendbars request message.
type ProtoOAGetTrendbarsReq struct {
	CtidTraderAccountId int64  `json:"ctidTraderAccountId"`
	FromTimestamp       int64  `json:"fromTimestamp"`
	ToTimestamp         int64  `json:"toTimestamp"`
	Period              int64  `json:"period"`
	SymbolId            int64  `json:"symbolId"`
	Count               uint32 `json:"count"`
}

// Get trendbars response message.
type ProtoOAGetTrendbarsRes struct {
	CtidTraderAccountId int64      `json:"ctidTraderAccountId"`
	Period              int64      `json:"period"`
	Timestamp           *int64     `json:"timestamp"`
	Trendbar            []Trendbar `json:"trendbar"`
	SymbolId            int64      `json:"symbolId"`
	ClosePrices         []float64  `json:"closePrices"`
}

type Trendbar struct {
	Volume                int64   `json:"volume"`
	Period                *int64  `json:"period"`
	Low                   int64   `json:"low"`
	DeltaOpen             *uint64 `json:"deltaOpen"`
	DeltaClose            uint64  `json:"deltaClose"`
	DeltaHigh             *uint64 `json:"deltaHigh"`
	UTCTimestampInMinutes *uint32 `json:"utcTimestampInMinutes"`
}

// Get current price request message.
type ProtoOASubscribeSpotsReq struct {
	CtidTraderAccountId int64   `json:"ctidTraderAccountId"`
	SymbolId            []int64 `json:"symbolId"`
}

// Get current price response message.
type ProtoOASpotEvent struct {
	CtidTraderAccountId int64   `json:"ctidTraderAccountId"`
	SymbolId            int64   `json:"symbolId"`
	Bid                 *uint64 `json:"bid"`
	Ask                 *uint64 `json:"ask"`
}

// Send new order request message.
type ProtoOANewOrderReq struct {
	CtidTraderAccountId int64    `json:"ctidTraderAccountId"`
	SymbolId            int64    `json:"symbolId"`
	OrderType           int64    `json:"orderType"`
	TradeSide           int64    `json:"tradeSide"`
	Volume              int64    `json:"volume"`
	LimitPrice          *float64 `json:"limitPrice"`
	StopPrice           *float64 `json:"stopPrice"`
	TimeInForce         *int64   `json:"timeInForce"`
	ExpirationTimestamp *int64   `json:"expirationTimestamp"`
	StopLoss            *float64 `json:"stopLoss"`
	TakeProfit          *float64 `json:"takeProfit"`
	Comment             *string  `json:"comment"`
	BaseSlippagePrice   *float64 `json:"baseSlippagePrice"`
	SlippageInPoints    *int64   `json:"slippageInPoints"`
	Label               *string  `json:"label"`
	PositionId          *int64   `json:"positionId"`
	ClientOrderId       *string  `json:"clientOrderId"`
	RelativeStopLoss    *int64   `json:"relativeStopLoss"`
	RelativeTakeProfit  *int64   `json:"relativeTakeProfit"`
	GuaranteedStopLoss  *bool    `json:"guaranteedStopLoss"`
	TrailingStopLoss    *bool    `json:"trailingStopLoss"`
	StopTriggerMethod   *int64   `json:"stopTriggerMethod"`
}

// Get current trader's informations request message.
type ProtoOATraderReq struct {
	CtidTraderAccountId int64 `json:"ctidTraderAccountId"`
}

// Get current trader's informations response message.
type ProtoOATraderRes struct {
	CtidTraderAccountId int64         `json:"ctidTraderAccountId"`
	Trader              ProtoOATrader `json:"trader"`
}

type ProtoOATrader struct {
	Balance int64 `json:"balance"`
}

// Get symbol's entity request message.
type ProtoOASymbolByIdReq struct {
	CtidTraderAccountId int64   `json:"ctidTraderAccountId"`
	SymbolId            []int64 `json:"symbolId"`
}

// Get symbol's entity response message.
type ProtoOASymbolByIdRes struct {
	CtidTraderAccountId int64          `json:"ctidTraderAccountId"`
	Symbol              []SymbolEntity `json:"symbol"`
}

type SymbolEntity struct {
	SymbolId    int64  `json:"symbolId"`
	PipPosition int32  `json:"pipPosition"`
	StepVolume  *int64 `json:"stepVolume"`
	LotSize     *int64 `json:"lotSize"`
}
