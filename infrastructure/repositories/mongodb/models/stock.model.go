package models

import (
	"fmt"
	"time"

	"github.com/hthl85/aws-tiprank-dividend-scraper/consts"
	"github.com/hthl85/aws-tiprank-dividend-scraper/entities"
	"github.com/hthl85/aws-tiprank-dividend-scraper/utils/datetime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockModel struct {
	ID               *primitive.ObjectID      `bson:"_id,omitempty"`
	IsActive         bool                     `bson:"isActive,omitempty"`
	CreatedAt        int64                    `bson:"createdAt,omitempty"`
	ModifiedAt       int64                    `bson:"modifiedAt,omitempty"`
	Schema           string                   `bson:"schema,omitempty"`
	Source           string                   `bson:"source,omitempty"`
	Ticker           string                   `bson:"ticker,omitempty"`
	Name             string                   `bson:"name,omitempty"`
	Type             string                   `bson:"type,omitempty"`
	AssetClass       string                   `bson:"assetClass,omitempty"`
	DividendSchedule string                   `bson:"dividendSchedule,omitempty"`
	Currency         string                   `bson:"currency,omitempty"`
	AllocationStock  float64                  `bson:"allocationStock,omitempty"`
	AllocationBond   float64                  `bson:"allocationBond,omitempty"`
	AllocationCash   float64                  `bson:"allocationCash,omitempty"`
	Sectors          []*SectorModel           `bson:"sector,omitempty"`
	Countries        []*CountryModel          `bson:"countries,omitempty"`
	DividendHistory  map[int64]*DividendModel `bson:"dividendHistory,omitempty"`
}

// DividendModel struct
type DividendModel struct {
	PayoutRatio    float64    `bson:"payoutRatio,omitempty"`
	Yield          float64    `bson:"yield,omitempty"`
	Dividend       float64    `bson:"dividend,omitempty"`
	ExDividendDate *time.Time `bson:"exDividendDate,omitempty"`
	RecordDate     *time.Time `bson:"recordDate,omitempty"`
	DividendDate   *time.Time `bson:"payoutDate,omitempty"`
}

// SectorModel struct
type SectorModel struct {
	SectorCode  string  `bson:"sectorCode,omitempty"`
	SectorName  string  `bson:"sectorName,omitempty"`
	FundPercent float64 `bson:"fundPercent,omitempty"`
}

// CountryModel struct
type CountryModel struct {
	CountryCode     string  `bson:"countryCode,omitempty"`
	CountryName     string  `bson:"countryName,omitempty"`
	HoldingStatCode string  `bson:"holdingStatCode,omitempty"`
	FundMktPercent  float64 `bson:"fundMktPercent,omitempty"`
	FundTnaPercent  float64 `bson:"fundTnaPercent,omitempty"`
}

// NewStockModel create stock model
func NewStockModel(e *entities.Stock, countryCode string) (*StockModel, error) {
	var m = &StockModel{}

	m.Source = consts.DATA_SOURCE
	m.Type = consts.SECURITY_TYPE

	m.Ticker = e.Ticker

	if e.Name != "" {
		m.Name = e.Name
	}

	m.AssetClass = consts.ASSET_CLASS

	if v, f := consts.Currencies[countryCode]; f {
		m.Currency = v.Code
	}

	m.AllocationStock = 100

	d, err := newDividendModel(e)
	if err != nil {
		return nil, err
	}

	m.DividendHistory = make(map[int64]*DividendModel)
	dividendTime := d.ExDividendDate.Unix()
	m.DividendHistory[dividendTime] = d

	c, err := newCountryModel(countryCode)
	if err != nil {
		return nil, err
	}

	m.Countries = []*CountryModel{c}

	s, err := newSectorModel(e.Sector)
	m.Sectors = []*SectorModel{s}

	return m, err
}

// newDividendModel create dividend model
func newDividendModel(e *entities.Stock) (*DividendModel, error) {
	exDividendDate, err := datetime.GetDateStartFromString(e.ExDividendDate)
	if err != nil {
		return nil, &DividendModelError{Message: fmt.Sprintf("parse exDividendDate failed : %v", err)}
	}

	recordDate, err := datetime.GetDateStartFromString(e.RecordDate)
	if err != nil {
		return nil, &DividendModelError{Message: fmt.Sprintf("parse recordDate failed : %v", err)}
	}

	dividendDate, err := datetime.GetDateStartFromString(e.DividendDate)
	if err != nil {
		return nil, &DividendModelError{Message: fmt.Sprintf("parse dividendDate failed : %v", err)}
	}

	m := &DividendModel{
		PayoutRatio:    e.PayoutRatio,
		Yield:          e.Yield,
		Dividend:       e.Dividend,
		ExDividendDate: exDividendDate,
		RecordDate:     recordDate,
		DividendDate:   dividendDate,
	}

	return m, err
}

// newSectorModel create sector model
func newSectorModel(code int64) (*SectorModel, error) {
	v, f := consts.Sectors[code]
	if !f {
		return &SectorModel{
			SectorName:  "Other",
			SectorCode:  "OTH",
			FundPercent: 100,
		}, &SectorModelError{Message: fmt.Sprintf("cannot find sector of code %d", code)}
	}

	return &SectorModel{
		SectorName:  v.Name,
		SectorCode:  v.Code,
		FundPercent: 100,
	}, nil
}

// newCountryModel create country model
func newCountryModel(code string) (*CountryModel, error) {
	v, f := consts.Countries[code]
	if !f {
		return nil, &CountryModelError{Message: fmt.Sprintf("cannot find country of code %s", code)}
	}

	return &CountryModel{
		CountryCode:    v.Alpha3Code,
		CountryName:    v.Name,
		FundMktPercent: 100,
		FundTnaPercent: 100,
	}, nil
}
