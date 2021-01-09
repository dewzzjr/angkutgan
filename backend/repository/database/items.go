package database

import (
	"context"
	"database/sql"
	"reflect"
	"strings"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const qGetListItems = `SELECT
	items.code, 
	name, 
	unit, 
	COALESCE(value, 0),
	COALESCE(inventory, 0),
	COALESCE(asset, 0)
FROM
	items 
LEFT JOIN 
	price_sell ON items.code = price_sell.code
LEFT JOIN
	stock ON items.code = stock.code
ORDER BY 
	items.create_time DESC
LIMIT ? OFFSET ?
`

// GetListItems using pagination
func (d *Database) GetListItems(ctx context.Context, limit, offset int) (items []model.Item, err error) {
	var rows *sqlx.Rows
	if rows, err = d.DB.QueryxContext(ctx, qGetListItems, limit, offset); err != nil {
		err = errors.Wrapf(err, "QueryxContext [%d, %d]", limit, offset)
		return
	}
	if items, err = d.GetRentByCodes(ctx, rows); err != nil {
		err = errors.Wrapf(err, "GetRentByCodes [%d, %d]", limit, offset)
	}
	return
}

const qGetListItemsByKeyword = `SELECT
	items.code, 
	name, 
	unit, 
	COALESCE(value, 0),
	COALESCE(inventory, 0),
	COALESCE(asset, 0)
FROM
	items 
LEFT JOIN 
	price_sell ON items.code = price_sell.code
LEFT JOIN
	stock ON items.code = stock.code
WHERE
	items.code = ? OR UPPER(name) LIKE CONCAT('%', ?, '%')
ORDER BY
	items.create_time DESC
LIMIT ? OFFSET ?
`

// GetListItemsByKeyword by keyword using pagination
func (d *Database) GetListItemsByKeyword(ctx context.Context, keyword string, limit, offset int) (items []model.Item, err error) {
	var rows *sqlx.Rows
	if rows, err = d.DB.QueryxContext(ctx, qGetListItemsByKeyword,
		strings.ToUpper(keyword),
		strings.ToUpper(keyword),
		limit,
		offset,
	); err != nil {
		err = errors.Wrapf(err, "QueryxContext [%d, %d, %s]", limit, offset, keyword)
		return
	}
	if items, err = d.GetRentByCodes(ctx, rows); err != nil {
		err = errors.Wrapf(err, "GetRentByCodes [%d, %d, %s]", limit, offset, keyword)
	}
	return
}

const qGetRentByCodes = `SELECT
	id,
	code,
	description,
	duration,
	time_unit,
	value
FROM
	price_rent
WHERE
	code IN (?)
ORDER BY 
	code ASC
`

// GetRentByCodes bulk multiple code
func (d *Database) GetRentByCodes(ctx context.Context, rows *sqlx.Rows) (items []model.Item, err error) {
	var i int
	index := make(map[string]int)
	items = make([]model.Item, 0)
	for rows.Next() {
		var item model.Item
		if err = rows.Scan(
			&item.Code,
			&item.Name,
			&item.Unit,
			&item.Price.Sell,
			&item.Available.Inventory,
			&item.Available.Asset,
		); err != nil {
			err = errors.Wrapf(err, "Scan")
			continue
		}
		item.Price.Rent = make([]model.PriceRent, 0)
		items = append(items, item)
		index[item.Code] = i
		i++
	}

	if len(items) == 0 {
		return
	}

	q, in, _ := sqlx.In(qGetRentByCodes, reflect.ValueOf(index).MapKeys())
	if rows, err = d.DB.QueryxContext(ctx, q, in...); err != nil {
		err = errors.Wrapf(err, "QueryxContext [%v]", in)
		return
	}

	for rows.Next() {
		var rent model.PriceRent
		var code string
		if err = rows.Scan(
			&rent.ID,
			&code,
			&rent.Description,
			&rent.Duration,
			&rent.TimeUnit,
			&rent.Value,
		); err != nil {
			err = errors.Wrapf(err, "Scan [%v]", in)
			continue
		}
		if i, ok := index[code]; ok {
			items[i].Price.Rent = append(items[i].Price.Rent, rent)
		}
	}
	return
}

const qUpdateInsertItem = `INSERT
INTO
	items (
		code,
		name,
		unit,
		modified_by
	)
VALUES ( ?, ?, ?, ? ) ON DUPLICATE KEY
UPDATE name = ?, unit = ?, modified_by = ?, update_time = CURRENT_TIMESTAMP
`

// UpdateInsertItem insert item or update if exists
func (d *Database) UpdateInsertItem(ctx context.Context, item model.Item, actionBy int64) (err error) {
	if _, err = d.DB.ExecContext(ctx, qUpdateInsertItem,
		// INSERT
		item.Code,
		item.Name,
		item.Unit,
		NullInt64(actionBy),
		// UPDATE
		item.Name,
		item.Unit,
		NullInt64(actionBy),
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [%v]", item)
	}
	return
}

const (
	qDeleteRentPrice = `DELETE
FROM
	price_rent
WHERE
	id = ? AND code = ?
`
	qInsertRentPrice = `INSERT
INTO
	price_rent (
		code,
		description,
		duration,
		time_unit,
		value,
		modified_by
	)
VALUES ( ?, ?, ?, ?, ?, ? )
`
)

// InsertDeleteRentPrice insert and delete rent price transaction
func (d *Database) InsertDeleteRentPrice(ctx context.Context, code string, insert []model.PriceRent, delete []int64, actionBy int64) (err error) {
	var tx *sqlx.Tx
	if tx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	for _, d := range delete {
		if _, err = tx.ExecContext(ctx, qDeleteRentPrice, d, code); err != nil {
			err = errors.Wrapf(err, "ExecContext [%s, %d]", code, d)
			tx.Rollback()
			return
		}
	}
	for _, i := range insert {
		if _, err = tx.ExecContext(ctx, qInsertRentPrice,
			code,
			i.Description,
			i.Duration,
			i.TimeUnit,
			i.Value,
			NullInt64(actionBy),
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [%s, %v]", code, i)
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		err = errors.Wrapf(err, "Commit [%s, %v, %v]", code, delete, insert)
	}
	return
}

const qGetPrintRent = `SELECT
	id,
	description,
	duration,
	time_unit,
	value
FROM
	price_rent
WHERE
	code = ?
`

// GetPriceRent list of price rent by code item
func (d *Database) GetPriceRent(ctx context.Context, code string) (prices []model.PriceRent, err error) {
	prices = make([]model.PriceRent, 0)
	if err = d.DB.SelectContext(ctx, &prices, qGetPrintRent, code); err != nil {
		err = errors.Wrapf(err, "SelectContext [%s]", code)
	}
	return
}

const qReplacePriceSell = `INSERT
INTO
	price_sell (
		code,
		value,
		modified_by
	)
VALUES ( ?, ?, ? ) ON DUPLICATE KEY
UPDATE value = ?, modified_by = ?, update_time = CURRENT_TIMESTAMP
`

// ReplacePriceSell add price sell or update if exists
func (d *Database) ReplacePriceSell(ctx context.Context, code string, value int, actionBy int64) (err error) {
	if _, err = d.DB.ExecContext(ctx, qReplacePriceSell,
		// INSERT
		code,
		value,
		NullInt64(actionBy),
		// UPDATE
		value,
		NullInt64(actionBy),
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s, %d]", code, value)
	}
	return
}

const qGetItemDetail = `SELECT
	items.code, 
	name, 
	unit, 
	COALESCE(value, 0),
	COALESCE(inventory, 0),
	COALESCE(asset, 0)
FROM
	items 
LEFT JOIN 
	price_sell ON items.code = price_sell.code
LEFT JOIN
	stock ON items.code = stock.code
WHERE
	items.code = ?
`

// GetItemDetail get item detail by code
// detail including price sell and rent
func (d *Database) GetItemDetail(ctx context.Context, code string) (item model.Item, err error) {
	if err = d.DB.QueryRowxContext(ctx, qGetItemDetail, code).Scan(
		&item.Code,
		&item.Name,
		&item.Unit,
		&item.Price.Sell,
		&item.Available.Inventory,
		&item.Available.Asset,
	); err != nil {
		err = errors.Wrapf(err, "QueryRowxContext [%s]", code)
		return
	}
	if item.Price.Rent, err = d.GetPriceRent(ctx, code); err != nil {
		err = errors.Wrap(err, "GetPriceRent")
	}
	return
}

const (
	qDeleteItem = `DELETE
FROM
	items
WHERE
	code = ?
`
	qDeleteSell = `DELETE
FROM
	price_sell
WHERE
	code = ?
`
	qDeleteRent = `DELETE
FROM
	price_rent
WHERE
	code = ?
`
)

// DeleteItem delete item including price
func (d *Database) DeleteItem(ctx context.Context, code string) (err error) {
	var tx *sqlx.Tx
	if tx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	if _, err = tx.ExecContext(ctx, qDeleteRent, code); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s]", code)
		tx.Rollback()
		return
	}
	if _, err = tx.ExecContext(ctx, qDeleteSell, code); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s]", code)
		tx.Rollback()
		return
	}
	if _, err = tx.ExecContext(ctx, qDeleteItem, code); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s]", code)
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		err = errors.Wrapf(err, "Commit [%s]", code)
	}
	return
}

const qIsValidItemCode = `SELECT
	code
FROM
	items
WHERE
	code = ?
`

// IsValidItemCode check is code is valid or not to be use
func (d *Database) IsValidItemCode(ctx context.Context, code string) (bool, error) {
	err := d.DB.QueryRowxContext(ctx, qIsValidItemCode, code).Scan(&code)
	if err == nil {
		return false, nil
	}
	if err != sql.ErrNoRows {
		err = errors.Wrapf(err, "QueryRowxContext [%s]", code)
		return false, err
	}
	return true, nil
}
