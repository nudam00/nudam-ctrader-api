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

type MongoDbConfig struct {
	Uri          string `json:"uri"`
	DatabaseName string `json:"databaseName"`
	Collection   string `json:"collection"`
}

type TraderConfiguration struct {
	PayloadTypes      map[string]int      `json:"payloadTypes"`
	Periods           map[string]Period   `json:"periods"`
	QuoteType         map[string]int64    `json:"quoteType"`
	OrderType         map[string]int64    `json:"orderType"`
	TradeSide         map[string]int64    `json:"tradeSide"`
	TimeInForce       map[string]int64    `json:"timeInForce"`
	StopTriggerMethod map[string]int64    `json:"stopTriggerMethod"`
	Pips              map[string]Dividers `json:"pips"`
	CurrencyPairs     []string            `json:"currencyPairs"`
}

type Period struct {
	Value      int64  `json:"value"`
	CountBars  uint32 `json:"countBars"`
	NumberDays uint32 `json:"numberDays"`
}

type Strategy struct {
	Ema      []float64 `json:"ema"`
	Risk     float64   `json:"risk"`
	Leverage int64     `json:"leverage"`
}

type Dividers struct {
	Pips    uint64 `json:"pips"`
	Price   uint64 `json:"price"`
	LotUnit uint64 `json:"lotUnit"`
}
