package consts

// Collection names
const (
	TIPRANK_DIVIDEND_LIST_COLLECTION = "tiprank_dividend_list" // Should match with Colnames's key of AppConf
)

// TipRank available countries
// var TipRankCountries = []string{"Canada", "US", "UK"}
var TipRankCountries = []string{"Canada", "US"}

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
