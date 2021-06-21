package tiprank

import (
	"context"

	"github.com/hthl85/aws-tiprank-dividend-scraper/entities"
)

///////////////////////////////////////////////////////////
// Stock Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
}

// Writer interface
type Writer interface {
	InsertTipRankDividend(ctx context.Context, tiprankDividend *entities.TipRankDividend, currency string) error
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
