package assets

var Periods = map[string]int{
	"M1":  1,
	"M2":  2,
	"M3":  3,
	"M4":  4,
	"M5":  5,
	"M10": 6,
	"M15": 7,
	"M30": 8,
	"H1":  9,
	"H4":  10,
	"H12": 11,
	"D1":  12,
	"W1":  13,
	"MN1": 14,
}

var Symbols []Symbol

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

type Trendbar struct {
	Volume                int64  `json:"volume"`
	Period                int    `json:"period"`
	Low                   int64  `json:"low"`
	DeltaOpen             uint64 `json:"deltaOpen"`
	DeltaClose            uint64 `json:"deltaClose"`
	DeltaHigh             uint64 `json:"deltaHigh"`
	UTCTimestampInMinutes uint32 `json:"utcTimestampInMinutes"`
}
