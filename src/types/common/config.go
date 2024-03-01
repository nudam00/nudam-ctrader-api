package common

type CTraderConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type CTraderAccountConfig struct {
	ClientId            string `json:"clientId"`
	ClientSecret        string `json:"clientSecret"`
	CtidTraderAccountId int64  `json:"ctidTraderAccountId"`
	AccessToken         string `json:"accessToken"`
}

type TraderConfiguration struct {
	PayloadTypes      map[string]int    `json:"payloadTypes"`
	Periods           map[string]Period `json:"periods"`
	QuoteType         map[string]int    `json:"quoteType"`
	OrderType         map[string]int    `json:"orderType"`
	TradeSide         map[string]int    `json:"tradeSide"`
	TimeInForce       map[string]int    `json:"timeInForce"`
	StopTriggerMethod map[string]int    `json:"stopTriggerMethod"`
}

type Period struct {
	Value      int    `json:"value"`
	CountBars  uint32 `json:"countBars"`
	NumberDays uint32 `json:"numberDays"`
}

type Strategy struct {
	Ema []int `json:"ema"`
}
