package models

import (
	"fmt"
	"time"

	"github.com/hthl85/aws-tiprank-dividend-scraper/entities"
	"github.com/hthl85/aws-tiprank-dividend-scraper/utils/datetime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockModel struct {
	ID              *primitive.ObjectID      `bson:"_id,omitempty"`
	IsActive        bool                     `bson:"isActive,omitempty"`
	CreatedAt       int64                    `bson:"createdAt,omitempty"`
	ModifiedAt      int64                    `bson:"modifiedAt,omitempty"`
	Schema          string                   `bson:"schema,omitempty"`
	Ticker          string                   `bson:"ticker,omitempty"`
	Name            string                   `bson:"name,omitempty"`
	Yield           float64                  `bson:"yield,omitempty"`
	DividendHistory map[int64]*DividendModel `bson:"dividendHistory,omitempty"`
}

// DividendModel struct
type DividendModel struct {
	Dividend       float64    `bson:"dividend,omitempty"`
	ExDividendDate *time.Time `bson:"exDividendDate,omitempty"`
	RecordDate     *time.Time `bson:"recordDate,omitempty"`
	DividendDate   *time.Time `bson:"payoutDate,omitempty"`
}

// NewStockModel create stock model
func NewStockModel(e *entities.Stock, countryCode string) (*StockModel, error) {
	var m = &StockModel{}

	m.Ticker = e.Ticker

	if e.Name != "" {
		m.Name = e.Name
	}

	m.Yield = e.Yield

	d, err := newDividendModel(e)
	if err != nil {
		return nil, err
	}

	m.DividendHistory = make(map[int64]*DividendModel)
	dividendTime := d.ExDividendDate.Unix()
	m.DividendHistory[dividendTime] = d

	return m, err
}

// newDividendModel create dividend model
func newDividendModel(e *entities.Stock) (*DividendModel, error) {
	exDividendDate, err := datetime.GetStarDateFromString(e.ExDividendDate)
	if err != nil {
		return nil, &DividendModelError{Message: fmt.Sprintf("parse exDividendDate failed : %v", err)}
	}

	recordDate, err := datetime.GetStarDateFromString(e.RecordDate)
	if err != nil {
		return nil, &DividendModelError{Message: fmt.Sprintf("parse recordDate failed : %v", err)}
	}

	dividendDate, err := datetime.GetStarDateFromString(e.DividendDate)
	if err != nil {
		return nil, &DividendModelError{Message: fmt.Sprintf("parse dividendDate failed : %v", err)}
	}

	m := &DividendModel{
		Dividend:       e.Amount,
		ExDividendDate: exDividendDate,
		RecordDate:     recordDate,
		DividendDate:   dividendDate,
	}

	return m, err
}
