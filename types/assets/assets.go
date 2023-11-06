package assets

var Symbols []Symbol

type Symbol struct {
	SymbolId         int     `json:"symbolId"`
	SymbolName       string  `json:"symbolName"`
	Enabled          bool    `json:"enabled"`
	BaseAssetId      int     `json:"baseAssetId"`
	QuoteAssetId     int     `json:"quoteAssetId"`
	SymbolCategoryId int     `json:"symbolCategoryId"`
	Description      string  `json:"description"`
	SortingNumber    float32 `json:"sortingNumber"`
}
