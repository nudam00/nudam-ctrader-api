package constants

var PayloadTypes = map[string]int{
	"ProtoOAApplicationAuthReq": 2100,
	"ProtoOAApplicationAuthRes": 2101,
	"ProtoOAAccountAuthReq":     2102,
	"ProtoOAAccountAuthRes":     2103,
	"ProtoOASymbolsListReq":     2114,
	"ProtoOASymbolsListRes":     2115,
	"ProtoOAGetTrendbarsReq":    2137,
	"ProtoOAGetTrendbarsRes":    2138,
}

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

var CountBars = map[string]uint32{
	"M1":  1440,
	"M2":  720,
	"M3":  480,
	"M4":  360,
	"M5":  288,
	"M10": 144,
	"M15": 96,
	"M30": 48,
	"H1":  24,
	"H4":  6,
	"H12": 2,
	"D1":  1,
	"W1":  1 / 7,
	"MN1": 1 / 30,
}
