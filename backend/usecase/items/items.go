package items

import (
	"context"
	"strings"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/package/pagination"
	"github.com/pkg/errors"
)

// GetList sequence of items
func (i *Items) GetList(ctx context.Context, page int, row int) (items []model.Item, err error) {
	if items, err = i.database.GetListItems(ctx, row, pagination.Offset(page, row)); err != nil {
		err = errors.Wrap(err, "GetListItems")
	}
	for i, item := range items {
		for j, rent := range item.Price.Rent {
			items[i].Price.Rent[j].TimeUnitDesc = rent.TimeUnit.String()
		}
	}
	return
}

// GetByKeyword sequence of items by keyword
func (i *Items) GetByKeyword(ctx context.Context, page int, row int, key string) (items []model.Item, err error) {
	if items, err = i.database.GetListItemsByKeyword(ctx,
		strings.TrimSpace(key),
		row,
		pagination.Offset(page, row),
		model.ColumnRents,
	); err != nil {
		err = errors.Wrap(err, "GetListItemsByKeyword")
	}
	for i, item := range items {
		for j, rent := range item.Price.Rent {
			items[i].Price.Rent[j].TimeUnitDesc = rent.TimeUnit.String()
		}
	}
	return
}

// Get item by code
func (i *Items) Get(ctx context.Context, code string) (item model.Item, err error) {
	if item, err = i.database.GetItemDetail(ctx, code); err != nil {
		err = errors.Wrap(err, "GetItemDetail")
	}
	for j, rent := range item.Price.Rent {
		item.Price.Rent[j].TimeUnitDesc = rent.TimeUnit.String()
	}
	return
}

// Create new item
func (i *Items) Create(ctx context.Context, item model.Item, actionBy int64) (err error) {
	if err = i.database.UpdateInsertItem(ctx, item, actionBy); err != nil {
		err = errors.Wrap(err, "UpdateInsertItem")
		return
	}
	if err = i.database.ReplaceStock(ctx, item.Code, model.Stock{
		Inventory: item.Stock,
		Asset:     item.Stock,
	}, actionBy); err != nil {
		err = errors.Wrap(err, "ReplaceStock")
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
	if (item.Name != "" || item.Unit != "") && (item.Name != get.Name || item.Unit != get.Unit) {
		if err = i.database.UpdateInsertItem(ctx, item, actionBy); err != nil {
			err = errors.Wrap(err, "UpdateInsertItem")
			return
		}
	}
	if item.Price.Sell != 0 && item.Price.Sell != get.Price.Sell {
		if err = i.database.ReplacePriceSell(ctx, item.Code, item.Price.Sell, actionBy); err != nil {
			err = errors.Wrap(err, "ReplacePriceSell")
			return
		}
	}
	if item.Price.Rent != nil {
		if err = i.UpdatePriceRent(ctx, item.Code, get.Price.Rent, item.Price.Rent, actionBy); err != nil {
			err = errors.Wrap(err, "UpdatePriceRent")
			return
		}
	}
	if item.Stock != 0 && get.Available.Asset != item.Stock {
		if err = i.database.ReplaceStock(ctx, item.Code, model.Stock{
			Asset:     item.Stock,
			Inventory: get.Available.Inventory + item.Stock - get.Available.Asset,
		}, actionBy); err != nil {
			err = errors.Wrap(err, "ReplaceStock")
			return
		}
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
