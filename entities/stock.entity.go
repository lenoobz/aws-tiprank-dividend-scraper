package entities

// Stock represents a stock entity
type Stock struct {
	Ticker         string  `json:"ticker,omitempty"`
	Name           string  `json:"company,omitempty"`
	Yield          float64 `json:"yield,omitempty"`
	Dividend       float64 `json:"amount,omitempty"`
	ExDividendDate string  `json:"exDate,omitempty"`
	RecordDate     string  `json:"recDate,omitempty"`
	DividendDate   string  `json:"payDate,omitempty"`
}
