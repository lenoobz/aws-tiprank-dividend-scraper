package main

import (
	"log"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-dividend-scraper/config"
	"github.com/hthl85/aws-tiprank-dividend-scraper/infrastructure/repositories/mongodb/repos"
	"github.com/hthl85/aws-tiprank-dividend-scraper/infrastructure/scraper"
	"github.com/hthl85/aws-tiprank-dividend-scraper/usecase/stock"
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
	repo, err := repos.NewStockMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create stock mongo repo failed")
	}
	defer repo.Close()

	// create new service
	fs := stock.NewService(repo, zap)

	// create new scraper jobs
	jobs := scraper.NewStockScraper(fs, zap)
	// jobs.StartSingleDayJob()
	// jobs.StartPreviousWeekJob()
	jobs.StartNextWeekJob()
	jobs.StartPreviousYearJob()
}
