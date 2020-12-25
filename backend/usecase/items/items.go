package items

import (
	"context"
	"strings"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetList sequence of items
func (i *Items) GetList(ctx context.Context, page int, row int) (items []model.Item, err error) {
	if items, err = i.database.GetListItems(ctx, row, offset(page, row)); err != nil {
		err = errors.Wrap(err, "GetListItems")
	}
	return
}

// GetByKeyword sequence of items by keyword
func (i *Items) GetByKeyword(ctx context.Context, page int, row int, key string) (items []model.Item, err error) {
	if items, err = i.database.GetListItemsByKeyword(ctx,
		strings.TrimSpace(key),
		row,
		offset(page, row),
	); err != nil {
		err = errors.Wrap(err, "GetListItemsByKeyword")
	}
	return
}

// Get item by code
func (i *Items) Get(ctx context.Context, code string) (item model.Item, err error) {
	if item, err = i.database.GetItemDetail(ctx, code); err != nil {
		err = errors.Wrap(err, "GetItemDetail")
	}
	return
}

// Create new item
func (i *Items) Create(ctx context.Context, item model.Item, actionBy int64) (err error) {
	if err = i.database.UpdateInsertItem(ctx, item, actionBy); err != nil {
		err = errors.Wrap(err, "UpdateInsertItem")
	}
	return
}

// Update item by code
func (i *Items) Update(ctx context.Context, item model.Item, actionBy int64) (err error) {
	var get model.Item
	if get, err = i.Get(ctx, item.Code); err != nil {
		err = errors.Wrapf(err, "Get")
		return
	}
	item.Name = get.Name
	item.Unit = get.Unit
	if err = i.database.UpdateInsertItem(ctx, item, actionBy); err != nil {
		err = errors.Wrap(err, "UpdateInsertItem")
		return
	}
	if item.Price.Sell != 0 && item.Price.Sell != get.Price.Sell {
		if err = i.database.ReplacePriceSell(ctx, item.Code, item.Price.Sell, actionBy); err != nil {
			err = errors.Wrap(err, "ReplacePriceSell")
			return
		}
	}
	if item.Price.Rent != nil {
		err = i.UpdatePriceRent(ctx, item.Code, get.Price.Rent, item.Price.Rent)
	}
	return
}

// Remove item by code
func (i *Items) Remove(ctx context.Context, code string) (err error) {
	// TODO check eligible to delete
	if err = i.database.DeleteItem(ctx, code); err != nil {
		err = errors.Wrap(err, "DeleteItem")
	}
	return
}

func offset(page, row int) int {
	if page < 1 {
		return 0
	}
	return (page - 1) * row
}
