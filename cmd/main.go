package main

import (
	"log"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-tiprank-dividend-scraper/config"
	"github.com/lenoobz/aws-tiprank-dividend-scraper/infrastructure/repositories/mongodb/repos"
	"github.com/lenoobz/aws-tiprank-dividend-scraper/infrastructure/scraper"
	"github.com/lenoobz/aws-tiprank-dividend-scraper/usecase/tiprank"
)

func main() {
	appConf := config.AppConf

	// create new logger
	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal("create app logger failed")
	}
	defer zap.Close()

	// create new repository
	tiprankDividendRepo, err := repos.NewTipRankDividendMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create TipRank dividend mongo failed")
	}
	defer tiprankDividendRepo.Close()

	// create new service
	tiprankDividendService := tiprank.NewService(tiprankDividendRepo, zap)

	// create new scraper jobs
	jobs := scraper.NewTipRankDividendScraper(tiprankDividendService, zap)
	// jobs.StartSingleDayJob()
	// jobs.StartPreviousWeekJob()
	jobs.StartNextWeekJob()
	jobs.StartPreviousYearJob()
}
