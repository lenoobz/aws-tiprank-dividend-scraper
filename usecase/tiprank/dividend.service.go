package tiprank

import (
	"context"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-dividend-scraper/entities"
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
func (s *Service) AddTipRankDividend(ctx context.Context, tiprankDividend *entities.TipRankDividend) error {
	s.log.Info(ctx, "adding TipRank dividend", "ticker", tiprankDividend.Ticker)
	return s.tiprankDividendRepo.InsertTipRankDividend(ctx, tiprankDividend)
}
