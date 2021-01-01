package items

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Items usecase object
type Items struct {
	database iDatabase
}

// New initiate usecase/items
func New(database iDatabase) *Items {
	return &Items{
		database: database,
	}
}

type iDatabase interface {
	GetItemDetail(ctx context.Context, code string) (item model.Item, err error)
	GetListItems(ctx context.Context, limit, offset int) (items []model.Item, err error)
	GetListItemsByKeyword(ctx context.Context, keyword string, limit, offset int) (items []model.Item, err error)
	UpdateInsertItem(ctx context.Context, item model.Item, actionBy int64) (err error)
	GetPriceRent(ctx context.Context, code string) (prices []model.PriceRent, err error)
	ReplacePriceSell(ctx context.Context, code string, value int, actionBy int64) (err error)
	InsertDeleteRentPrice(ctx context.Context, code string, insert []model.PriceRent, delete []int64) (err error)
	DeleteItem(ctx context.Context, code string) (err error)

	ReplaceStock(ctx context.Context, code string, stock model.Stock, actionBy int64) (err error)
}
