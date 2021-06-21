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
	InsertTipRankDividend(context.Context, *entities.TipRankDividend) error
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
