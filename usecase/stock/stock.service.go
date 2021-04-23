package stock

import (
	"context"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-dividend-scraper/entities"
)

// Service sector
type Service struct {
	repo Repo
	log  logger.ContextLog
}

// NewService create new service
func NewService(r Repo, l logger.ContextLog) *Service {
	return &Service{
		repo: r,
		log:  l,
	}
}

// AddStock creates new stock
func (s *Service) AddStock(ctx context.Context, stock *entities.Stock, countryCode string) error {
	s.log.Info(ctx, "adding stock", "ticker", stock.Ticker, "country", countryCode)
	return s.repo.InsertStock(ctx, stock, countryCode)
}
