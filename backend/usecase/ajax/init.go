package ajax

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
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
	IsValidItemCode(ctx context.Context, code string) (bool, error)
	GetListItemsByKeyword(ctx context.Context, keyword string, limit, offset int) (items []model.Item, err error)
}
