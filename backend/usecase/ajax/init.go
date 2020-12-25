package ajax

import (
	"context"
)

// Ajax usecase object
type Ajax struct {
	database iDatabase
}

// New initiate usecase/ajax
func New(database iDatabase) *Ajax {
	return &Ajax{
		database: database,
	}
}

type iDatabase interface {
	IsValidUsername(ctx context.Context, username string) (bool, error)
}
