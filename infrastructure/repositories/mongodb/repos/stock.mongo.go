package repos

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-tiprank-dividend-scraper/config"
	"github.com/hthl85/aws-tiprank-dividend-scraper/consts"
	"github.com/hthl85/aws-tiprank-dividend-scraper/entities"
	"github.com/hthl85/aws-tiprank-dividend-scraper/infrastructure/repositories/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StockMongo struct
type StockMongo struct {
	db     *mongo.Database
	client *mongo.Client
	log    logger.ContextLog
	conf   *config.MongoConfig
}

// NewStockMongo creates new stock mongo repo
func NewStockMongo(db *mongo.Database, l logger.ContextLog, conf *config.MongoConfig) (*StockMongo, error) {
	if db != nil {
		return &StockMongo{
			db:   db,
			log:  l,
			conf: conf,
		}, nil
	}

	// set context with timeout from the config
	// create new context for the query
	ctx, cancel := createContext(context.Background(), conf.TimeoutMS)
	defer cancel()

	// set mongo client options
	clientOptions := options.Client()

	// set min pool size
	if conf.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(conf.MinPoolSize)
	}

	// set max pool size
	if conf.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(conf.MaxPoolSize)
	}

	// set max idle time ms
	if conf.MaxIdleTimeMS > 0 {
		clientOptions.SetMaxConnIdleTime(time.Duration(conf.MaxIdleTimeMS) * time.Millisecond)
	}

	// construct a connection string from mongo config object
	cxnString := fmt.Sprintf("mongodb+srv://%s:%s@%s", conf.Username, conf.Password, conf.Host)

	// create mongo client by making new connection
	client, err := mongo.Connect(ctx, clientOptions.ApplyURI(cxnString))
	if err != nil {
		return nil, err
	}

	return &StockMongo{
		db:     client.Database(conf.Dbname),
		client: client,
		log:    l,
		conf:   conf,
	}, nil
}

// Close disconnect from database
func (r *StockMongo) Close() {
	ctx := context.Background()
	r.log.Info(ctx, "close mongo client")

	if r.client == nil {
		return
	}

	if err := r.client.Disconnect(ctx); err != nil {
		r.log.Error(ctx, "disconnect mongo failed", "error", err)
	}
}

///////////////////////////////////////////////////////////////////////////////
// Implement interface
///////////////////////////////////////////////////////////////////////////////

// InsertStock insert new stock stock
func (r *StockMongo) InsertStock(ctx context.Context, stock *entities.Stock, countryCode string) error {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	savedStock, err := r.findStockByTicker(ctx, stock.Ticker)
	if err != nil {
		r.log.Error(ctx, "find stock by ticker failed", "error", err, "ticker", stock.Ticker)
		return err
	}

	insertingStock, err := models.NewStockModel(stock, countryCode)
	if err != nil {
		// ignore sector model error, it isn't crucial
		var sectorErr *models.SectorModelError
		if !errors.As(err, &sectorErr) {
			r.log.Error(ctx, "create model failed", "error", err, "ticker", stock.Ticker)
			return err
		}

		// log noncrucial error
		r.log.Warn(ctx, "create model failed but ignored", "error", err, "ticker", stock.Ticker)
	}

	if savedStock != nil {
		// Copy dividend history from saved stock to inserting stock
		for k, v := range savedStock.DividendHistory {
			if _, f := insertingStock.DividendHistory[k]; !f {
				insertingStock.DividendHistory[k] = v
			}
		}
	}

	if err = r.insertStock(ctx, insertingStock); err != nil {
		r.log.Error(ctx, "insert stock failed", "error", err, "ticker", stock.Ticker)
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////
// Implement helper function
///////////////////////////////////////////////////////////

// createContext create a new context with timeout
func createContext(ctx context.Context, t uint64) (context.Context, context.CancelFunc) {
	timeout := time.Duration(t) * time.Millisecond
	return context.WithTimeout(ctx, timeout*time.Millisecond)
}

// insertStock inserts new stock
func (r *StockMongo) insertStock(ctx context.Context, m *models.StockModel) error {
	if m == nil {
		r.log.Error(ctx, "invalid param")
		return fmt.Errorf("invalid param")
	}

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.TIPRANK_STOCK_COL]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	m.IsActive = true
	m.Schema = r.conf.SchemaVersion
	m.ModifiedAt = time.Now().UTC().Unix()

	filter := bson.D{{
		Key:   "ticker",
		Value: m.Ticker,
	}}

	update := bson.D{
		{
			Key:   "$set",
			Value: m,
		},
		{
			Key: "$setOnInsert",
			Value: bson.D{{
				Key:   "createdAt",
				Value: time.Now().UTC().Unix(),
			}},
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.log.Error(ctx, "update one failed", "error", err)
		return err
	}

	return nil
}

// findStockByTicker finds stock of a given ticker
func (r *StockMongo) findStockByTicker(ctx context.Context, ticker string) (*models.StockModel, error) {
	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.TIPRANK_STOCK_COL]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return nil, fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	// filter
	filter := bson.D{
		{
			Key:   "ticker",
			Value: strings.ToUpper(ticker),
		},
	}

	// find options
	findOptions := options.FindOne()

	var stock models.StockModel
	if err := col.FindOne(ctx, filter, findOptions).Decode(&stock); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			r.log.Info(ctx, "stock not found", "ticker", ticker)
			return nil, nil
		}

		r.log.Error(ctx, "decode find one failed", "error", err, "ticker", ticker)
		return nil, err
	}

	return &stock, nil
}
