package database

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const qGetListTransactions = `SELECT
	DATE_FORMAT(tx.date, '%d/%m/%Y') as date,
	tx.id as id,
	ss.address as address,
	ss.total_price as total_price,
	COALESCE(ss.deposit, 0) as deposit,
	COALESCE(ss.discount, 0) as discount,
	COALESCE(ss.shipping_fee, 0) as shipping_fee,
	COALESCE(DATE_FORMAT(tx.done_date, '%d/%m/%Y'), '') as done_date,
	COALESCE(DATE_FORMAT(tx.paid_date, '%d/%m/%Y'), '') as paid_date,
	COALESCE(ss.project, 0) as project_id,
	COALESCE(pr.name, '') as project_name,
	cs.code as customer_code, 
	cs.name as customer_name, 
	cs.type as customer_type,
	cs.address as customer_address,
	cs.phone as customer_phone, 
	COALESCE(cs.nik, '') as nik,
	COALESCE(cs.role, '') as role,
	COALESCE(cs.group_name, '') as group_name
FROM
	transactions AS tx
JOIN
	snapshot AS ss ON tx.id = ss.t_id
JOIN
	customers AS cs ON tx.customer = cs.code
LEFT JOIN
	projects AS pr ON ss.project = pr.id
WHERE tx.type = ? AND tx.date <= ?
ORDER BY tx.date DESC
LIMIT ? OFFSET ?
`

// GetListTransactions with pagination
func (d *Database) GetListTransactions(ctx context.Context, txType model.TransactionType, date time.Time, limit, offset int) (txs []model.Transaction, err error) {
	var rows *sqlx.Rows
	if rows, err = d.DB.QueryxContext(ctx, qGetListTransactions,
		txType,
		date,
		limit,
		offset,
	); err != nil {
		err = errors.Wrapf(err, "QueryxContext [%s, %d, %d]", txType, limit, offset)
		return
	}
	if txs, err = d.GetSnapshotTx(ctx, rows); err != nil {
		err = errors.Wrap(err, "GetSnapshotTx")
		return
	}
	return
}

const qGetListTransactionsByCustomer = `SELECT
	DATE_FORMAT(tx.date, '%d/%m/%Y') as date,
	tx.id as id,
	ss.address as address,
	ss.total_price as total_price,
	COALESCE(ss.deposit, 0) as deposit,
	COALESCE(ss.discount, 0) as discount,
	COALESCE(ss.shipping_fee, 0) as shipping_fee,
	COALESCE(DATE_FORMAT(tx.done_date, '%d/%m/%Y'), '') as done_date,
	COALESCE(DATE_FORMAT(tx.paid_date, '%d/%m/%Y'), '') as paid_date,
	COALESCE(ss.project, 0) as project_id,
	COALESCE(pr.name, '') as project_name,
	cs.code as customer_code, 
	cs.name as customer_name, 
	cs.type as customer_type,
	cs.address as customer_address,
	cs.phone as customer_phone, 
	COALESCE(cs.nik, '') as nik,
	COALESCE(cs.role, '') as role,
	COALESCE(cs.group_name, '') as group_name
FROM
	transactions AS tx
JOIN
	snapshot AS ss 
		ON tx.id = ss.t_id
JOIN
	customers AS cs ON tx.customer = cs.code
LEFT JOIN
	projects AS pr ON ss.project = pr.id
WHERE tx.type = ? AND tx.customer = ? AND tx.date <= ?
ORDER BY tx.date DESC
LIMIT ? OFFSET ?
`

// GetListTransactionsByCustomer by customer code using pagination
func (d *Database) GetListTransactionsByCustomer(ctx context.Context, customer string, txType model.TransactionType, date time.Time, limit, offset int) (txs []model.Transaction, err error) {
	var rows *sqlx.Rows
	if rows, err = d.DB.QueryxContext(ctx, qGetListTransactionsByCustomer,
		txType,
		strings.ToUpper(customer),
		date,
		limit,
		offset,
	); err != nil {
		err = errors.Wrapf(err, "QueryxContext [%s, %s, %d, %d]", txType, customer, limit, offset)
		return
	}
	if txs, err = d.GetSnapshotTx(ctx, rows); err != nil {
		err = errors.Wrap(err, "GetSnapshotTx")
		return
	}
	return
}

// GetSnapshotTx bulk snapshot
func (d *Database) GetSnapshotTx(ctx context.Context, rows *sqlx.Rows) (txs []model.Transaction, err error) {
	defer rows.Close()
	for rows.Next() {
		var tx model.Transaction
		if err = rows.Scan(
			&tx.Date,
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
			err = errors.Wrapf(err, "Scan")
			continue
		}
		if tx.Snapshot.Items, err = d.GetSnapshotItems(ctx, tx.ID); err != nil {
			err = errors.Wrapf(err, "GetSnapshotItems")
			continue
		}
		txs = append(txs, tx)
	}
	return
}

const qGetTransaction = `SELECT
	tx.id,
	ss.address,
	ss.total_price,
	COALESCE(ss.deposit, 0),
	COALESCE(ss.discount, 0),
	COALESCE(ss.shipping_fee, 0),
	COALESCE(DATE_FORMAT(tx.done_date, '%d/%m/%Y'), ''),
	COALESCE(DATE_FORMAT(tx.paid_date, '%d/%m/%Y'), ''),
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
	customers AS cs ON tx.customer = cs.code
LEFT JOIN
	projects AS pr ON ss.project = pr.id
`

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
		err = errors.Wrapf(err, "QueryRowxContext [%s, %v]", code, date)
		return
	}
	tx.Date = date.Format(model.DateFormat)
	if tx.Items, err = d.GetSnapshotItems(ctx, tx.ID); err != nil {
		err = errors.Wrapf(err, "GetSnapshotItems [%s, %v]", code, date)
	}
	return
}

const qGetSnapshotItems = `SELECT
	t.id,
	t.item,
	i.name,
	t.amount,
	t.price,
	COALESCE(t.claim, 0) AS claim,
	COALESCE(t.time_unit, 0) AS time_unit,
	COALESCE(t.duration, 0) AS duration,
	COALESCE(t.amount, 0) - SUM(COALESCE(p.amount, 0)) - SUM(COALESCE(s.amount, 0)) AS need_shipment,
	SUM(COALESCE(n.amount, 0)) AS extend_amount,
	COALESCE(p.prev_snapshot, 0) AS previous_id
FROM
	snapshot_item AS t
JOIN
	items AS i ON t.item = i.code AND t.t_id = ?
LEFT JOIN
	shipment AS s ON s.i_id = t.id
LEFT JOIN
	extends AS n ON t.id = n.next_snapshot
LEFT JOIN
	extends AS p ON t.id = p.prev_snapshot
WHERE n.next_snapshot IS NULL
GROUP BY t.id
`

// GetSnapshotItems by transaction id
func (d *Database) GetSnapshotItems(ctx context.Context, txID int64) (items []model.SnapshotItem, err error) {
	items = make([]model.SnapshotItem, 0)
	if err = d.DB.SelectContext(ctx, &items, qGetSnapshotItems, txID); err != nil {
		err = errors.Wrapf(err, "SelectContext [%d]", txID)
	}
	return
}

const qGetTxID = `SELECT
	id
FROM
	transactions
WHERE
	date = ? AND customer = ? AND type = ?
LIMIT 1
`

// GetTransactionID by date, customer code, type
func (d *Database) GetTransactionID(ctx context.Context, date time.Time, code string, txType model.TransactionType) (txID int64, err error) {
	if err = d.DB.QueryRowxContext(ctx, qGetTxID,
		date,
		code,
		txType,
	).Scan(&txID); err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "QueryRowxContext [%s, %s, %v]", txType.String(), code, date)
		return
	}
	return txID, nil
}

const (
	qUpdateTransaction = `UPDATE transactions
SET
	modified_by = ?,
	update_time = CURRENT_TIMESTAMP
WHERE
	id = ?
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
)

// UpdateTransaction change transaction and snapshot
func (d *Database) UpdateTransaction(ctx context.Context, txID int64, txType model.TransactionType, tx model.CreateTransaction, actionBy int64) (err error) {
	var txx *sqlx.Tx
	if txx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
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

// InsertTransaction new transaction and snapshot
func (d *Database) InsertTransaction(ctx context.Context, txType model.TransactionType, tx model.CreateTransaction, actionBy int64) (err error) {
	var txx *sqlx.Tx
	if txx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	date, _ := time.Parse(model.DateFormat, tx.Date)
	res, e := txx.ExecContext(ctx, qInsertTransaction,
		date,
		tx.Customer,
		txType,
		NullInt64(actionBy),
	)
	if e != nil {
		err = errors.Wrapf(e, "ExecContext [qInsertTransaction, %s, %v]", txType.String(), tx)
		_ = txx.Rollback()
		return
	}
	var txID int64
	if txID, err = res.LastInsertId(); err != nil {
		err = errors.Wrapf(err, "LastInsertId [%s, %v]", txType.String(), tx)
		_ = txx.Rollback()
		return
	}
	if _, err = txx.ExecContext(ctx, qInsertSnapshot,
		txID,
		tx.Address,
		NullInt64(tx.ProjectID),
		tx.Deposit,
		tx.Discount,
		tx.ShippingFee,
		tx.TotalPrice,
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [qInsertSnapshot, %s, %v]", txType.String(), tx)
		_ = txx.Rollback()
		return
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

const qUpdatePaidDate = `UPDATE transactions
SET
	paid_date = ?
WHERE
	t_id = ?
`

// UpdatePaidDate update payment date when already paid off
func (d *Database) UpdatePaidDate(ctx context.Context, txID int64, date time.Time) (err error) {
	if _, err = d.DB.ExecContext(ctx, qUpdatePaidDate,
		date,
		txID,
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [%d]", txID)
		return
	}
	return
}

const qGetTotalPayment = `SELECT
	(total_price + deposit + shipping_fee) AS amount
FROM
	snapshot
WHERE
	t_id = ?
`

// GetTotalPayment get total payment by transaction id
func (d *Database) GetTotalPayment(ctx context.Context, txID int64) (amount int, err error) {
	if err = d.DB.QueryRowxContext(ctx, qGetTotalPayment, txID).Scan(&amount); err != nil {
		err = errors.Wrapf(err, "QueryRowxContext [%d]", txID)
		return
	}
	return
}
