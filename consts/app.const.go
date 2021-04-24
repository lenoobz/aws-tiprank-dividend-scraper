package consts

// Collection names
const (
	TIPRANK_STOCK_COL = "stock" // Should match with Colnames's key of AppConf
)

const (
	SECURITY_TYPE = "STOCK"
	ASSET_CLASS   = "EQUITY"
	DATA_SOURCE   = "TIP_RANK"
)

// TipRank available countries
// var TipRankCountries = []string{"Canada", "US", "UK"}
var TipRankCountries = []string{"Canada"}

/*
Bellow is the list of sector get from tiprank by inspect html page
Utilities: 17351
Basic Materials: 17343
Financial: 17346
Healthcare: 17347
Industrial Goods: 17348
Technology: 17349
Consumer Goods: 18731
Services: 17350
*/

var Sectors = map[int64]struct {
	Name string
	Code string
}{
	17343: {
		Name: "Basic Materials",
		Code: "MTL",
	},
	17346: {
		Name: "Financials",
		Code: "FIN",
	},
	17347: {
		Name: "Health Care",
		Code: "HEC",
	},
	17348: {
		Name: "Industrials",
		Code: "IND",
	},
	17349: {
		Name: "Technology",
		Code: "TEC",
	},
	17350: {
		Name: "Consumer Services",
		Code: "OTH",
	},
	17351: {
		Name: "Utilities",
		Code: "UTL",
	},
	18731: {
		Name: "Consumer Goods",
		Code: "CNS",
	},
}

var Countries = map[string]struct {
	Name       string
	Alpha2Code string
	Alpha3Code string
	NumberCode int
	Latitude   int
	Longitude  int
}{
	"Canada": {
		Name:       "Canada",
		Alpha2Code: "CA",
		Alpha3Code: "CAN",
		NumberCode: 124,
		Latitude:   60,
		Longitude:  -95,
	},
	"US": {
		Name:       "United States",
		Alpha2Code: "US",
		Alpha3Code: "USA",
		NumberCode: 840,
		Latitude:   38,
		Longitude:  -97,
	},
	"UK": {
		Name:       "United Kingdom",
		Alpha2Code: "GB",
		Alpha3Code: "GBR",
		NumberCode: 826,
		Latitude:   54,
		Longitude:  -2,
	},
}

var Currencies = map[string]struct {
	Code string
}{
	"Canada": {
		Code: "CAD",
	},
	"US": {
		Code: "USD",
	},
	"UK": {
		Code: "GBP",
	},
}
