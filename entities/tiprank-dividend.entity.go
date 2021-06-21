package entities

// TipRankDividend struct
type TipRankDividend struct {
	Ticker         string  `json:"ticker,omitempty"`
	Name           string  `json:"company,omitempty"`
	Yield          float64 `json:"yield,omitempty"`
	Amount         float64 `json:"amount,omitempty"`
	ExDividendDate string  `json:"exDate,omitempty"`
	RecordDate     string  `json:"recDate,omitempty"`
	DividendDate   string  `json:"payDate,omitempty"`
}
