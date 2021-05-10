package scraper

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/google/uuid"
	corid "github.com/hthl85/aws-lambda-corid"
	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-dividend-scraper/config"
	"github.com/hthl85/aws-tiprank-dividend-scraper/consts"
	"github.com/hthl85/aws-tiprank-dividend-scraper/entities"
	"github.com/hthl85/aws-tiprank-dividend-scraper/usecase/stock"
)

// StockScraper struct
type StockScraper struct {
	StockJob     *colly.Collector
	stockService *stock.Service
	log          logger.ContextLog
}

// NewStockScraper create new stock scraper
func NewStockScraper(ss *stock.Service, l logger.ContextLog) *StockScraper {
	sj := newScraperJob()

	return &StockScraper{
		StockJob:     sj,
		stockService: ss,
		log:          l,
	}
}

// newScraperJob creates a new colly collector with some custom configs
func newScraperJob() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(config.AllowDomain),
		colly.Async(true),
	)

	// Overrides the default timeout (10 seconds) for this collector
	c.SetRequestTimeout(30 * time.Second)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  config.DomainGlob,
		Parallelism: 2,
		RandomDelay: 10 * time.Second,
	})

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	return c
}

// configJobs configs on error handler and on response handler for scaper jobs
func (s *StockScraper) configJobs() {
	s.StockJob.OnError(s.errorHandler)
	s.StockJob.OnResponse(s.processDividendResponse)
}

// StartSingleDayJob start job
func (s *StockScraper) StartSingleDayJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		date := time.Now()

		url := config.GetDividendStockByDateURL(countryCode, date)

		s.log.Info(ctx, "scraping dividend stock", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

		if err := s.StockJob.Request("GET", url, nil, reqContext, nil); err != nil {
			s.log.Error(ctx, "scrape dividend stock list failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
		}
	}

	s.StockJob.Wait()
}

// StartPreviousWeekJob start job
func (s *StockScraper) StartPreviousWeekJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		for i := 0; i <= 7; i++ {
			date := time.Now().AddDate(0, 0, -i)

			url := config.GetDividendStockByDateURL(countryCode, date)

			s.log.Info(ctx, "scraping dividend stock", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

			if err := s.StockJob.Request("GET", url, nil, reqContext, nil); err != nil {
				s.log.Error(ctx, "scrape dividend stock list failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
			}
		}
	}

	s.StockJob.Wait()
}

// StartNextWeekJob start job
func (s *StockScraper) StartNextWeekJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		for i := 0; i <= 7; i++ {
			date := time.Now().AddDate(0, 0, i)

			url := config.GetDividendStockByDateURL(countryCode, date)

			s.log.Info(ctx, "scraping dividend stock", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

			if err := s.StockJob.Request("GET", url, nil, reqContext, nil); err != nil {
				s.log.Error(ctx, "scrape dividend stock list failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
			}
		}
	}

	s.StockJob.Wait()
}

// StartPreviousYearJob start job
func (s *StockScraper) StartPreviousYearJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		for i := 0; i <= 365; i++ {
			date := time.Now().AddDate(0, 0, -i)

			url := config.GetDividendStockByDateURL(countryCode, date)

			s.log.Info(ctx, "scraping dividend stock", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

			if err := s.StockJob.Request("GET", url, nil, reqContext, nil); err != nil {
				s.log.Error(ctx, "scrape dividend stock list failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
			}
		}
	}

	s.StockJob.Wait()
}

// StartDailyJob start job
func (s *StockScraper) StartDailyJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		// scrape dividend stock of next week
		date := time.Now().AddDate(0, 0, 7)

		url := config.GetDividendStockByDateURL(countryCode, date)

		s.log.Info(ctx, "scraping dividend stock", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

		if err := s.StockJob.Request("GET", url, nil, reqContext, nil); err != nil {
			s.log.Error(ctx, "scrape dividend stock list failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
		}
	}

	s.StockJob.Wait()
}

///////////////////////////////////////////////////////////
// Scraper Handler
///////////////////////////////////////////////////////////

// errorHandler generic error handler for all scaper jobs
func (s *StockScraper) errorHandler(r *colly.Response, err error) {
	ctx := context.Background()
	s.log.Error(ctx, "failed to request url", "url", r.Request.URL, "error", err)
}

func (s *StockScraper) processDividendResponse(r *colly.Response) {
	// create correlation if for processing fund list
	id, _ := uuid.NewRandom()
	ctx := corid.NewContext(context.Background(), id)
	s.log.Info(ctx, "scraping dividend stock")

	var stocks []*entities.Stock

	// unmarshal response data to stock struct
	if err := json.Unmarshal(r.Body, &stocks); err != nil {
		s.log.Error(ctx, "unmarshal response failed", "error", err)
		return
	}

	countryCode := r.Request.Ctx.Get("country")

	for _, stock := range stocks {
		if err := s.stockService.AddStock(ctx, stock, countryCode); err != nil {
			s.log.Error(ctx, "add dividend stock failed", "error", err, "ticker", stock.Ticker)
		}
	}
}
