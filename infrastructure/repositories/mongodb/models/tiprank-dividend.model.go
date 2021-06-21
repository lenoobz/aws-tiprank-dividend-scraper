package models

import (
	"context"
	"time"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-dividend-scraper/entities"
	"github.com/hthl85/aws-tiprank-dividend-scraper/utils/datetime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TipRankDividendModel struct {
	ID              *primitive.ObjectID             `bson:"_id,omitempty"`
	CreatedAt       int64                           `bson:"createdAt,omitempty"`
	ModifiedAt      int64                           `bson:"modifiedAt,omitempty"`
	Enabled         bool                            `bson:"enabled"`
	Deleted         bool                            `bson:"deleted"`
	Schema          string                          `bson:"schema,omitempty"`
	Ticker          string                          `bson:"ticker,omitempty"`
	Name            string                          `bson:"name,omitempty"`
	Yield           float64                         `bson:"yield,omitempty"`
	Currency        string                          `bson:"currency,omitempty"`
	DividendHistory map[int64]*DividendHistoryModel `bson:"dividendHistory,omitempty"`
}

// DividendHistoryModel struct
type DividendHistoryModel struct {
	Dividend       float64    `bson:"dividend,omitempty"`
	ExDividendDate *time.Time `bson:"exDividendDate,omitempty"`
	RecordDate     *time.Time `bson:"recordDate,omitempty"`
	DividendDate   *time.Time `bson:"payoutDate,omitempty"`
}

// NewTipRankDividendModel create stock model
func NewTipRankDividendModel(ctx context.Context, log logger.ContextLog, tiprankDividend *entities.TipRankDividend, currency string, schemaVersion string) (*TipRankDividendModel, error) {
	var tiprankDividendModel = &TipRankDividendModel{
		ModifiedAt:      time.Now().UTC().Unix(),
		Enabled:         true,
		Deleted:         false,
		Schema:          schemaVersion,
		Ticker:          tiprankDividend.Ticker,
		Name:            tiprankDividend.Name,
		Yield:           tiprankDividend.Yield,
		Currency:        currency,
		DividendHistory: map[int64]*DividendHistoryModel{},
	}

	dividendHistoryModel, err := newDividendHistoryModel(ctx, log, tiprankDividend)

	tiprankDividendModel.DividendHistory = make(map[int64]*DividendHistoryModel)
	dividendTime := dividendHistoryModel.ExDividendDate.Unix()
	tiprankDividendModel.DividendHistory[dividendTime] = dividendHistoryModel

	return tiprankDividendModel, err
}

// newDividendHistoryModel create dividend history model
func newDividendHistoryModel(ctx context.Context, log logger.ContextLog, tiprankDividend *entities.TipRankDividend) (*DividendHistoryModel, error) {
	dividendHistoryModel := &DividendHistoryModel{
		Dividend: tiprankDividend.Amount,
	}

	exDividendDate, err := datetime.GetStarDateFromString(tiprankDividend.ExDividendDate)
	if err != nil {
		log.Error(ctx, "parse exDividendDate failed", "error", err)
	} else {
		dividendHistoryModel.ExDividendDate = exDividendDate
	}

	recordDate, err := datetime.GetStarDateFromString(tiprankDividend.RecordDate)
	if err != nil {
		log.Error(ctx, "parse recordDate failed", "error", err)
	} else {
		dividendHistoryModel.RecordDate = recordDate
	}

	dividendDate, err := datetime.GetStarDateFromString(tiprankDividend.DividendDate)
	if err != nil {
		log.Error(ctx, "parse dividendDate failed", "error", err)
	} else {
		dividendHistoryModel.DividendDate = dividendDate
	}

	return dividendHistoryModel, err
}
