package repos

import (
	"context"
	"fmt"
	"strings"
	"time"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-tiprank-dividend-scraper/config"
	"github.com/lenoobz/aws-tiprank-dividend-scraper/consts"
	"github.com/lenoobz/aws-tiprank-dividend-scraper/entities"
	"github.com/lenoobz/aws-tiprank-dividend-scraper/infrastructure/repositories/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TipRankDividendMongo struct
type TipRankDividendMongo struct {
	db     *mongo.Database
	client *mongo.Client
	log    logger.ContextLog
	conf   *config.MongoConfig
}

// NewTipRankDividendMongo creates new stock mongo repo
func NewTipRankDividendMongo(db *mongo.Database, log logger.ContextLog, conf *config.MongoConfig) (*TipRankDividendMongo, error) {
	if db != nil {
		return &TipRankDividendMongo{
			db:   db,
			log:  log,
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

	return &TipRankDividendMongo{
		db:     client.Database(conf.Dbname),
		client: client,
		log:    log,
		conf:   conf,
	}, nil
}

// Close disconnect from database
func (r *TipRankDividendMongo) Close() {
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

// InsertTipRankDividend insert new Tiprank dividend
func (r *TipRankDividendMongo) InsertTipRankDividend(ctx context.Context, tiprankDividend *entities.TipRankDividend, currency string) error {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	savedTipRankDividend, err := r.findTipRankDividendByTicker(ctx, tiprankDividend.Ticker)
	if err != nil {
		r.log.Error(ctx, "find stock by ticker failed", "error", err, "ticker", tiprankDividend.Ticker)
		return err
	}

	newTipRankDividend, err := models.NewTipRankDividendModel(ctx, r.log, tiprankDividend, currency, r.conf.SchemaVersion)
	if err != nil {
		// log noncrucial error
		r.log.Warn(ctx, "create model failed but ignored", "error", err, "ticker", tiprankDividend.Ticker)
	}

	if savedTipRankDividend != nil {
		// Copy dividend history from saved stock to inserting stock
		for k, v := range savedTipRankDividend.DividendHistory {
			if _, f := newTipRankDividend.DividendHistory[k]; !f {
				newTipRankDividend.DividendHistory[k] = v
			}
		}
	}

	if err = r.insertTipRankDividend(ctx, newTipRankDividend); err != nil {
		r.log.Error(ctx, "insert TipRank dividend failed", "error", err, "ticker", tiprankDividend.Ticker)
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////
// Implement helper function
///////////////////////////////////////////////////////////

// findTipRankDividendByTicker finds TipRank dividend of a given ticker
func (r *TipRankDividendMongo) findTipRankDividendByTicker(ctx context.Context, ticker string) (*models.TipRankDividendModel, error) {
	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.TIPRANK_DIVIDEND_LIST_COLLECTION]
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

	var tiprankDividendModel models.TipRankDividendModel
	if err := col.FindOne(ctx, filter, findOptions).Decode(&tiprankDividendModel); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			r.log.Info(ctx, "TipRank dividend not found", "ticker", ticker)
			return nil, nil
		}

		r.log.Error(ctx, "decode find one failed", "error", err, "ticker", ticker)
		return nil, err
	}

	return &tiprankDividendModel, nil
}

// insertTipRankDividend inserts TipRank dividend
func (r *TipRankDividendMongo) insertTipRankDividend(ctx context.Context, tiprankDividendModel *models.TipRankDividendModel) error {
	if tiprankDividendModel == nil {
		r.log.Error(ctx, "invalid param")
		return fmt.Errorf("invalid param")
	}

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.TIPRANK_DIVIDEND_LIST_COLLECTION]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	filter := bson.D{{
		Key:   "ticker",
		Value: tiprankDividendModel.Ticker,
	}}

	update := bson.D{
		{
			Key:   "$set",
			Value: tiprankDividendModel,
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
