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
	"github.com/hthl85/aws-tiprank-dividend-scraper/usecase/tiprank"
)

// TipRankDividendScraper struct
type TipRankDividendScraper struct {
	ScrapeTipRankDividendJob *colly.Collector
	tiprankDividendService   *tiprank.Service
	log                      logger.ContextLog
	errorTickers             []string
}

// NewTipRankDividendScraper create new TipRank dividend scraper
func NewTipRankDividendScraper(tiprankDividendService *tiprank.Service, log logger.ContextLog) *TipRankDividendScraper {
	scrapeTipRankDividendJob := newScraperJob()

	return &TipRankDividendScraper{
		ScrapeTipRankDividendJob: scrapeTipRankDividendJob,
		tiprankDividendService:   tiprankDividendService,
		log:                      log,
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
func (s *TipRankDividendScraper) configJobs() {
	s.ScrapeTipRankDividendJob.OnError(s.errorHandler)
	s.ScrapeTipRankDividendJob.OnResponse(s.processDividendResponse)
}

// StartSingleDayJob start job
func (s *TipRankDividendScraper) StartSingleDayJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		date := time.Now()

		url := config.GetDividendStockByDateURL(countryCode, date)

		s.log.Info(ctx, "scraping TipRank dividend", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

		if err := s.ScrapeTipRankDividendJob.Request("GET", url, nil, reqContext, nil); err != nil {
			s.log.Error(ctx, "scrape TipRank dividend failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
		}
	}

	s.ScrapeTipRankDividendJob.Wait()
}

// StartPreviousWeekJob start job
func (s *TipRankDividendScraper) StartPreviousWeekJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		for i := 0; i <= 7; i++ {
			date := time.Now().AddDate(0, 0, -i)

			url := config.GetDividendStockByDateURL(countryCode, date)

			s.log.Info(ctx, "scraping TipRank dividend", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

			if err := s.ScrapeTipRankDividendJob.Request("GET", url, nil, reqContext, nil); err != nil {
				s.log.Error(ctx, "scrape TipRank dividend failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
			}
		}
	}

	s.ScrapeTipRankDividendJob.Wait()
}

// StartNextWeekJob start job
func (s *TipRankDividendScraper) StartNextWeekJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		for i := 0; i <= 7; i++ {
			date := time.Now().AddDate(0, 0, i)

			url := config.GetDividendStockByDateURL(countryCode, date)

			s.log.Info(ctx, "scraping TipRank dividend", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

			if err := s.ScrapeTipRankDividendJob.Request("GET", url, nil, reqContext, nil); err != nil {
				s.log.Error(ctx, "scrape TipRank dividend failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
			}
		}
	}

	s.ScrapeTipRankDividendJob.Wait()
}

// StartPreviousYearJob start job
func (s *TipRankDividendScraper) StartPreviousYearJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		for i := 0; i <= 365; i++ {
			date := time.Now().AddDate(0, 0, -i)

			url := config.GetDividendStockByDateURL(countryCode, date)

			s.log.Info(ctx, "scraping TipRank dividend", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

			if err := s.ScrapeTipRankDividendJob.Request("GET", url, nil, reqContext, nil); err != nil {
				s.log.Error(ctx, "scrape TipRank dividend failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
			}
		}
	}

	s.ScrapeTipRankDividendJob.Wait()
}

// StartDailyJob start job
func (s *TipRankDividendScraper) StartDailyJob() {
	ctx := context.Background()

	s.configJobs()

	for _, countryCode := range consts.TipRankCountries {
		reqContext := colly.NewContext()
		reqContext.Put("country", countryCode)

		// scrape dividend stock of next week
		date := time.Now().AddDate(0, 0, 7)

		url := config.GetDividendStockByDateURL(countryCode, date)

		s.log.Info(ctx, "scraping TipRank dividend", "country", countryCode, "date", date.Format("2006-01-02"), "url", url)

		if err := s.ScrapeTipRankDividendJob.Request("GET", url, nil, reqContext, nil); err != nil {
			s.log.Error(ctx, "scrape TipRank dividend failed", "error", err, "country", countryCode, "date", date.Format("2006-01-02"))
		}
	}

	s.ScrapeTipRankDividendJob.Wait()
}

///////////////////////////////////////////////////////////
// Scraper Handler
///////////////////////////////////////////////////////////

// errorHandler generic error handler for all scaper jobs
func (s *TipRankDividendScraper) errorHandler(r *colly.Response, err error) {
	ctx := context.Background()
	s.log.Error(ctx, "failed to request url", "url", r.Request.URL, "error", err)
}

func (s *TipRankDividendScraper) processDividendResponse(r *colly.Response) {
	// create correlation if for processing fund list
	id, _ := uuid.NewRandom()
	ctx := corid.NewContext(context.Background(), id)
	s.log.Info(ctx, "processDividendResponse")

	var tiprankDividends []*entities.TipRankDividend

	// unmarshal response data
	if err := json.Unmarshal(r.Body, &tiprankDividends); err != nil {
		s.log.Error(ctx, "unmarshal response failed", "error", err)
		return
	}

	for _, tiprankDividend := range tiprankDividends {
		if err := s.tiprankDividendService.AddTipRankDividend(ctx, tiprankDividend); err != nil {
			s.log.Error(ctx, "add TipRank dividend failed", "error", err, "ticker", tiprankDividend.Ticker)
		} else {
			s.errorTickers = append(s.errorTickers, tiprankDividend.Ticker)
		}
	}
}

// Close scraper
func (s *TipRankDividendScraper) Close() []string {
	s.log.Info(context.Background(), "DONE - SCRAPING TIPRANK DIVIDENDS", "tickers", s.errorTickers)
	return s.errorTickers
}
