package database

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

const qGetPayments = `SELECT
	p.id,
	p.name,
	p.amount,
	p.method,
	p.account,
	DATE_FORMAT(p.date, '%d/%m/%Y') AS date,
	COALESCE(u.fullname, '') AS accept_by
FROM
	payments AS p
LEFT JOIN
	user_info AS u ON p.accept_by = u.user_id
WHERE
	t_id = ?
ORDER BY date DESC
`

// GetPayments by transaction id
func (d *Database) GetPayments(ctx context.Context, txID int64) (payments []model.Payment, err error) {
	payments = make([]model.Payment, 0)
	if err = d.DB.SelectContext(ctx, &payments, qGetPayments, txID); err != nil {
		err = errors.Wrapf(err, "SelectContext [%d]", txID)
		return
	}
	return
}

const qGetLastPayment = `SELECT
	p.id,
	p.name,
	p.amount,
	p.method,
	p.account,
	DATE_FORMAT(p.date, '%d/%m/%Y') AS date
FROM
	payments AS p
WHERE
	t_id = ?
ORDER BY p.date DESC
LIMIT 1
`

// GetLastPayment last payment in a transaction
func (d *Database) GetLastPayment(ctx context.Context, txID int64) (payment model.Payment, err error) {
	if err = d.DB.QueryRowxContext(ctx, qGetLastPayment, txID).StructScan(&payment); err != nil {
		err = errors.Wrapf(err, "QueryRowxContext [%d]", txID)
		return
	}
	return
}

const qInsertPayment = `INSERT
INTO 
	payments (
		t_id,
		date,
		name,
		amount,
		method,
		account,
		accept_by
	)
VALUES ( ?, ?, ?, ?, ?, ?, ? )
`

// InsertPayment new payment in a transaction
func (d *Database) InsertPayment(ctx context.Context, txID int64, payment model.Payment, actionBy int64) (err error) {
	var date time.Time
	if date, err = time.Parse(model.DateFormat, payment.Date); err != nil {
		return
	}
	if _, err = d.DB.ExecContext(ctx, qInsertPayment,
		txID,
		date,
		payment.Name,
		payment.Amount,
		payment.Method,
		payment.Account,
		NullInt64(actionBy),
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [%d]", txID)
		return
	}
	return
}

const qUpdatePayment = `UPDATE payments 
SET 
	t_id = ?,
	date = ?,
	name = ?,
	amount = ?,
	method = ?,
	account = ?,
	accept_by = ?
WHERE
	id = ?
`

// UpdatePayment edit last payment in a transaction
func (d *Database) UpdatePayment(ctx context.Context, txID int64, payment model.Payment, actionBy int64) (err error) {
	var get model.Payment
	if get, err = d.GetLastPayment(ctx, txID); err != nil {
		err = errors.Wrapf(err, "GetLastPayment [%d]", txID)
		return
	}
	var date time.Time
	if date, err = time.Parse(model.DateFormat, payment.Date); err != nil {
		return
	}
	if _, err = d.DB.ExecContext(ctx, qUpdatePayment,
		txID,
		date,
		payment.Name,
		payment.Amount,
		payment.Method,
		payment.Account,
		NullInt64(actionBy),
		get.ID,
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [qUpdatePayment, %d]", txID)
		return
	}
	return
}

const qDeletePayment = `DELETE
FROM
	payments 
WHERE
	id = ?
`

// DeletePayment remove last payment in a transaction
func (d *Database) DeletePayment(ctx context.Context, txID int64) (err error) {
	var get model.Payment
	if get, err = d.GetLastPayment(ctx, txID); err != nil {
		err = errors.Wrapf(err, "GetLastPayment [%d]", txID)
		return
	}
	if _, err = d.DB.ExecContext(ctx, qDeletePayment,
		get.ID,
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [qDeletePayment, %d]", txID)
		return
	}
	return
}

const qPaidAmount = `SELECT 
	COALESCE(SUM(amount), 0) AS paid_amount,
	COALESCE(DATE_FORMAT(MAX(date), '%d/%m/%Y'), '') AS last_payment
FROM
	payments
WHERE
	t_id = ? AND account = 100
`

// GetLastPaidAmount get last date payment and total amount paid
func (d *Database) GetLastPaidAmount(ctx context.Context, txID int64) (amount int, date string, err error) {
	if err = d.DB.QueryRowxContext(ctx, qPaidAmount, txID).Scan(
		&amount,
		&date,
	); err != nil {
		err = errors.Wrapf(err, "QueryRowxContext [GetLastPaidAmount, %d]", txID)
		return
	}
	return
}
