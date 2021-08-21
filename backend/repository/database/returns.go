package database

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const qGetReturnByDate = `SELECT
	DATE_FORMAT(s.date, '%d/%m/%Y') AS s_date,
	r.id AS id,
	r.s_id AS s_id,
	i.item AS code,
	r.amount AS amount,
	r.status AS status,
	r.claim AS claim
FROM
	returns r 
JOIN 
	shipment s ON r.s_id = s.id
JOIN 
	snapshot_item i ON s.i_id = i.id
WHERE r.t_id ? AND r.date = ?
`

// GetReturnByDate return by date in a transaction
func (d *Database) GetReturnByDate(ctx context.Context, txID int64, date time.Time) (returns model.Return, err error) {
	returns = model.Return{
		Date:  date.Format(model.DateFormat),
		Items: make([]model.ReturnItem, 0),
	}
	if err = d.DB.SelectContext(ctx, &returns.Items, qGetReturnByDate, txID, date); err != nil {
		err = errors.Wrapf(err, "SelectContext [%d, %s]", txID, date)
		return
	}
	return
}

const (
	qGetReturns = `SELECT
	DATE_FORMAT(r.date, '%d/%m/%Y') AS date,
	DATE_FORMAT(s.date, '%d/%m/%Y') AS s_date,
	DATE_FORMAT(s.deadline, '%d/%m/%Y') AS deadline,
	r.id AS id,
	r.s_id AS s_id,
	i.item AS code,
	r.amount AS amount,
	r.status AS status,
	r.claim AS claim
FROM
	returns r 
JOIN 
	shipment s ON r.s_id = s.id
JOIN 
	snapshot_item i ON s.i_id = i.id
WHERE
	r.t_id = ?
ORDER BY r.date DESC
`
	qGetDateReturns = `SELECT DISTINCT
	date
FROM
	returns
WHERE
	t_id = ?
ORDER BY date DESC
`
)

// GetReturns by transaction id
func (d *Database) GetReturns(ctx context.Context, txID int64) (returns []model.Return, err error) {
	var rows *sqlx.Rows
	if rows, err = d.DB.QueryxContext(ctx, qGetDateReturns, txID); err != nil {
		err = errors.Wrapf(err, "QueryxContext [qGetDateReturns, %d]", txID)
		return
	}
	defer rows.Close()

	mapIndex := make(map[string]int)
	returns = make([]model.Return, 0)
	for rows.Next() {
		var date time.Time
		if err = rows.Scan(&date); err != nil {
			err = errors.Wrapf(err, "Scan [qGetDateReturns, %d]", txID)
			return
		}
		r := model.Return{
			Date:  date.Format(model.DateFormat),
			Items: make([]model.ReturnItem, 0),
		}
		returns = append(returns, r)
		mapIndex[r.Date] = len(returns) - 1
	}
	rows.Close()

	if rows, err = d.DB.QueryxContext(ctx, qGetReturns, txID); err != nil {
		err = errors.Wrapf(err, "QueryxContext [qGetReturns, %d]", txID)
		return
	}

	for rows.Next() {
		var date string
		var item model.ReturnItem
		if err = rows.Scan(
			&date,
			&item.ShipmentDate,
			&item.ID,
			&item.ShipmentID,
			&item.Code,
			&item.Amount,
			&item.Status,
			&item.Claim,
		); err != nil {
			err = errors.Wrapf(err, "Scan [qGetReturns, %d]", txID)
			continue
		}
		index := mapIndex[date]
		returns[index].Items = append(returns[index].Items, item)
	}
	return
}

const (
	qInsertReturn = `INSERT
INTO 
	returns (
		t_id,
		date,
		s_id,
		amount,
		status,
		claim,
		modified_by
	)
VALUES ( ?, ?, ?, ?, ?, ? )
`
	qDeleteReturnByDate = `DELETE
FROM
	returns
WHERE t_id ? AND date = ?
`
)

// DeleteInsertReturn if exist delete then add return in a transaction
func (d *Database) DeleteInsertReturn(ctx context.Context, txID int64, returns model.Return, isDelete bool, actionBy int64) (err error) {
	var date time.Time
	if date, err = time.Parse(model.DateFormat, returns.Date); err != nil {
		return
	}
	var tx *sqlx.Tx
	if tx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	if isDelete {
		if _, err = tx.ExecContext(ctx, qDeleteReturnByDate,
			txID,
			date,
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [qDeleteReturnByDate, %d]", txID)
			_ = tx.Rollback()
			return
		}
	}
	for _, item := range returns.Items {
		if _, err = tx.ExecContext(ctx, qInsertReturn,
			txID,
			date,
			item.ShipmentID,
			item.Amount,
			item.Status,
			item.Claim,
			NullInt64(actionBy),
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [qInsertReturn, %d]", txID)
			_ = tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		err = errors.Wrapf(err, "Commit [%d]", txID)
		_ = tx.Rollback()
	}
	return
}
