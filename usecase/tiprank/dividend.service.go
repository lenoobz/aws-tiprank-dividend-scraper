package tiprank

import (
	"context"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-dividend-scraper/entities"
	"github.com/hthl85/aws-tiprank-dividend-scraper/utils/currency"
)

// Service sector
type Service struct {
	tiprankDividendRepo Repo
	log                 logger.ContextLog
}

// NewService create new service
func NewService(tiprankDividendRepo Repo, log logger.ContextLog) *Service {
	return &Service{
		tiprankDividendRepo: tiprankDividendRepo,
		log:                 log,
	}
}

// AddTipRankDividend add TipRank dividend
func (s *Service) AddTipRankDividend(ctx context.Context, tiprankDividend *entities.TipRankDividend, country string) error {
	s.log.Info(ctx, "adding TipRank dividend", "ticker", tiprankDividend.Ticker)

	currency, err := currency.GetCountryCurrency(country)
	if err != nil {
		s.log.Error(ctx, "get country currency failed", "country", country)
	}

	return s.tiprankDividendRepo.InsertTipRankDividend(ctx, tiprankDividend, currency)
}
