package database

import (
	"context"

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
