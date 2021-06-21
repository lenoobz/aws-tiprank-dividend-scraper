package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-dividend-scraper/config"
	"github.com/hthl85/aws-tiprank-dividend-scraper/infrastructure/repositories/mongodb/repos"
	"github.com/hthl85/aws-tiprank-dividend-scraper/infrastructure/scraper"
	"github.com/hthl85/aws-tiprank-dividend-scraper/usecase/tiprank"
)

func main() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context) ([]string, error) {
	log.Println("lambda handler is called")

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
	jobs.StartDailyJob()

	tickers := jobs.Close()
	return tickers, nil
}
