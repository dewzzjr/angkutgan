package database

import (
	"context"
	"database/sql"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

const qReplaceStock = `INSERT
INTO
	stock (
		code,
		inventory,
		asset,
		modified_by
	)
VALUES ( ?, ?, ?, ? ) ON DUPLICATE KEY
UPDATE inventory = ?, asset = ?, modified_by = ?, update_time = CURRENT_TIMESTAMP
`

// ReplaceStock add stock or update if exists
func (d *Database) ReplaceStock(ctx context.Context, code string, stock model.Stock, actionBy int64) (err error) {
	if _, err = d.DB.ExecContext(ctx, qReplaceStock,
		// INSERT
		code,
		stock.Inventory,
		stock.Asset,
		NullInt64(actionBy),
		// UPDATE
		stock.Inventory,
		stock.Asset,
		NullInt64(actionBy),
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s, %v]", code, stock)
	}
	return
}

const (
	qGetStock = `SELECT asset, inventory
FROM stock
WHERE code = ?
`
	qChangeStock = `UPDATE stock
SET 
	asset = ?, 
	inventory = ?, 
	modified_by = ?, 
	update_time = CURRENT_TIMESTAMP
WHERE code = ?
`
)

// ChangeStock increment or decrement asset or inventory stock
func (d *Database) ChangeStock(ctx context.Context, code string, number int, stype model.StockType, actionBy int64) (err error) {
	var asset, inventory int
	if err = d.DB.QueryRowxContext(ctx, qGetStock, code).Scan(&asset, &inventory); err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "QueryRowxContext [%s]", code)
		return
	}
	switch stype {
	case model.StockAsset:
		asset = asset + number
	case model.StockInventory:
		inventory = inventory + number
	default:
		return
	}
	if _, err = d.DB.ExecContext(ctx, qChangeStock,
		asset,
		inventory,
		NullInt64(actionBy),
		code,
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s, %d, %d]", code, asset, inventory)
	}
	return
}
