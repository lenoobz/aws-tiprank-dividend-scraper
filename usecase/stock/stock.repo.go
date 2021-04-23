package stock

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
	InsertStock(context.Context, *entities.Stock, string) error
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
