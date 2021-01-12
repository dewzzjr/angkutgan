package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	qGetTransaction = `SELECT
	tx.id,
	ss.address,
	ss.total_price,
	COALESCE(ss.deposit, 0),
	COALESCE(ss.discount, 0),
	COALESCE(ss.shipping_fee, 0),
	COALESCE(DATE_FORMAT(ss.done_date, '%d/%m/%Y'), ''),
	COALESCE(DATE_FORMAT(ss.paid_date, '%d/%m/%Y'), ''),
	COALESCE(ss.project, 0),
	COALESCE(pr.name, ''),
	cs.code, 
	cs.name, 
	cs.type, 
	cs.address,
	cs.phone, 
	COALESCE(cs.nik, ''),
	COALESCE(cs.role, ''),
	COALESCE(cs.group_name, '')
FROM
	transactions AS tx
JOIN
	snapshot AS ss 
		ON tx.id = ss.t_id
		AND tx.date = ?
		AND tx.customer = ?
		AND tx.type = ?
JOIN
	customer AS cs ON tx.customer = cs.code
LEFT JOIN
	projects AS pr ON ss.project = pr.id
`
	qGetSnapshotItems = `SELECT
	s.id,
	s.item,
	i.name,
	s.amount,
	s.price,
	s.claim,
	s.time_unit,
	s.duration
FROM
	snapshot_item AS s
JOIN
	items AS i ON s.item = i.code AND s.t_id = ?
`
)

// GetTransaction by customer code, date, and transaction type
func (d *Database) GetTransaction(ctx context.Context, date time.Time, code string, txType model.TransactionType) (tx model.Transaction, err error) {
	if err = d.DB.QueryRowxContext(ctx, qGetTransaction,
		date,
		code,
		txType,
	).Scan(
		&tx.ID,
		&tx.Address,
		&tx.TotalPrice,
		&tx.Deposit,
		&tx.Discount,
		&tx.ShippingFee,
		&tx.DoneDate,
		&tx.PaidDate,
		&tx.ProjectID,
		&tx.ProjectName,
		&tx.Customer.Code,
		&tx.Customer.Name,
		&tx.Customer.Type,
		&tx.Customer.Address,
		&tx.Customer.Phone,
		&tx.Customer.NIK,
		&tx.Customer.Role,
		&tx.Customer.GroupName,
	); err != nil {
		err = errors.Wrapf(err, "QueryRowxContext [%s, %s, %v]", txType.String(), code, date)
		return
	}
	tx.Items = make([]model.SnapshotItem, 0)
	if err = d.DB.SelectContext(ctx, &tx.Items, qGetSnapshotItems, tx.ID); err != nil {
		err = errors.Wrapf(err, "SelectContext [%s, %s, %v]", txType.String(), code, date)
	}
	return
}

const (
	qGetTxID = `SELECT
	id
FROM
	transactions
WHERE
	date = STR_TO_DATE(?, '%d/%m/%Y') AND customer = ? AND type = ?
LIMIT 1
`
	qInsertTransaction = `INSERT
INTO
	transactions (
		date,
		customer,
		type,
		modified_by
	)
VALUES ( ?, ?, ?, ? )
`
	qUpdateTransaction = `UPDATE transactions
SET
	modified_by = ?,
	update_time = CURRENT_TIMESTAMP
WHERE
	id = ?
`
	qInsertSnapshot = `INSERT
INTO
	snapshot (
		t_id,
		address,
		project,
		deposit,
		discount,
		shipping_fee,
		total_price
	)
VALUES ( ?, ?, ?, ?, ?, ?, ? )
`
	qUpdateSnapshot = `UPDATE snapshot
SET
	address = ?,
	project = ?,
	deposit = ?,
	discount = ?,
	shipping_fee = ?,
	total_price = ?
WHERE
	t_id = ?
`
	qInsertSnapshotItem = `INSERT
INTO
	snapshot_item (
		t_id,
		item,
		amount,
		price,
		claim,
		time_unit,
		duration
	)
VALUE ( ?, ?, ?, ?, ?, ?, ? )
`
)

// UpdateInsertTransaction insert transaction or update if exists
func (d *Database) UpdateInsertTransaction(ctx context.Context, txType model.TransactionType, tx model.CreateTransaction, actionBy int64) (err error) {
	var txID int64
	var txx *sqlx.Tx
	if err = d.DB.QueryRowxContext(ctx, qGetTxID,
		tx.Date,
		tx.Customer,
		txType,
	).Scan(&txID); err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "QueryRowxContext [%s, %v]", txType.String(), tx)
		return
	}
	if txx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	if txID == 0 {
		res, e := txx.ExecContext(ctx, qInsertTransaction,
			tx.Date,
			tx.Customer,
			txType,
			NullInt64(actionBy),
		)
		if err != nil {
			err = errors.Wrapf(e, "ExecContext [qInsertTransaction, %s, %v]", txType.String(), tx)
			_ = txx.Rollback()
			return
		}
		if txID, err = res.LastInsertId(); err != nil {
			err = errors.Wrapf(err, "LastInsertId [%s, %v]", txType.String(), tx)
			_ = txx.Rollback()
			return
		}
		if _, err = txx.ExecContext(ctx, qInsertSnapshot,
			txID,
			tx.Address,
			NullInt64(tx.ProjectID),
			NullInt(tx.Deposit),
			NullInt(tx.Discount),
			NullInt(tx.ShippingFee),
			tx.TotalPrice,
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [qInsertSnapshot, %s, %v]", txType.String(), tx)
			_ = txx.Rollback()
			return
		}
	} else {
		if _, err = txx.ExecContext(ctx, qUpdateTransaction,
			NullInt64(actionBy),
			txID,
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [qUpdateTransaction, %s, %d]", txType.String(), txID)
			_ = txx.Rollback()
			return
		}
		if _, err = txx.ExecContext(ctx, qUpdateSnapshot,
			tx.Address,
			NullInt64(tx.ProjectID),
			NullInt(tx.Deposit),
			NullInt(tx.Discount),
			NullInt(tx.ShippingFee),
			tx.TotalPrice,
			txID,
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [qUpdateSnapshot, %s, %d]", txType.String(), txID)
			_ = txx.Rollback()
			return
		}
		if _, err = txx.ExecContext(ctx, qDeleteSnapshotItem, txID); err != nil {
			err = errors.Wrapf(err, "ExecContext [qDeleteSnapshotItem, %s, %v]", txType.String(), tx)
			_ = txx.Rollback()
			return
		}
	}
	for _, item := range tx.Items {
		if _, err = txx.ExecContext(ctx, qInsertSnapshotItem,
			txID,
			item.Code,
			item.Amount,
			item.Price,
			NullInt(item.Claim),
			NullInt(int(item.TimeUnit)),
			NullInt(item.Duration),
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [qInsertSnapshotItem, %s, %d, %v]", txType.String(), txID, tx)
			_ = txx.Rollback()
			return
		}
	}
	if err = txx.Commit(); err != nil {
		err = errors.Wrapf(err, "Commit [%s, %v]", txType.String(), tx)
		_ = txx.Rollback()
	}
	return
}

const (
	qDeleteTransaction = `DELETE
FROM
	transactions
WHERE
	id = ?
`
	qDeleteSnapshot = `DELETE
FROM
	snapshot
WHERE
	t_id = ?
`
	qDeleteSnapshotItem = `DELETE
FROM
	snapshot_item
WHERE
	t_id = ?
`
)

// DeleteTransaction by transaction id
func (d *Database) DeleteTransaction(ctx context.Context, txID int64) (err error) {
	var tx *sqlx.Tx
	if tx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	if _, err = tx.ExecContext(ctx, qDeleteSnapshotItem, txID); err != nil {
		err = errors.Wrapf(err, "ExecContext [qDeleteSnapshotItem, %d]", txID)
		_ = tx.Rollback()
		return
	}
	if _, err = tx.ExecContext(ctx, qDeleteSnapshot, txID); err != nil {
		err = errors.Wrapf(err, "ExecContext [qDeleteSnapshot, %d]", txID)
		_ = tx.Rollback()
		return
	}
	if _, err = tx.ExecContext(ctx, qDeleteTransaction, txID); err != nil {
		err = errors.Wrapf(err, "ExecContext [qDeleteTransaction, %d]", txID)
		_ = tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		err = errors.Wrapf(err, "Commit [%d]", txID)
		_ = tx.Rollback()
	}
	return
}
