package database

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	qGetShipments = `SELECT
	DATE_FORMAT(s.date, '%d/%m/%Y') AS date,
	s.id AS id,
	s.i_id AS i_id,
	i.item AS code,
	s.amount AS amount,
	COALESCE(DATE_FORMAT(s.deadline, '%d/%m/%Y'), '') AS deadline
FROM
	shipment s JOIN snapshot_item i ON s.i_id = i.id
WHERE
	s.t_id = ?
ORDER BY s.date DESC
`
	qGetDateShipments = `SELECT DISTINCT
	date
FROM
	shipment
WHERE
	t_id = ?
ORDER BY date DESC
`
)

// GetShipments by transaction id
func (d *Database) GetShipments(ctx context.Context, txID int64) (shipment []model.Shipment, err error) {
	var rows *sqlx.Rows
	if rows, err = d.DB.QueryxContext(ctx, qGetDateShipments, txID); err != nil {
		err = errors.Wrapf(err, "QueryxContext [qGetDateShipments, %d]", txID)
		return
	}

	mapIndex := make(map[string]int)
	shipment = make([]model.Shipment, 0)
	for rows.Next() {
		var date time.Time
		if err = rows.Scan(&date); err != nil {
			err = errors.Wrapf(err, "Scan [qGetDateShipments, %d]", txID)
			continue
		}
		s := model.Shipment{
			Date:  date.Format(model.DateFormat),
			Items: make([]model.ShipmentItem, 0),
		}
		shipment = append(shipment, s)
		mapIndex[s.Date] = len(shipment) - 1
	}
	rows.Close()

	if rows, err = d.DB.QueryxContext(ctx, qGetShipments, txID); err != nil {
		err = errors.Wrapf(err, "QueryxContext [qGetShipments, %d]", txID)
		return
	}

	for rows.Next() {
		var date string
		var item model.ShipmentItem
		if err = rows.Scan(
			&date,
			&item.ID,
			&item.ItemID,
			&item.Code,
			&item.Amount,
			&item.Deadline,
		); err != nil {
			err = errors.Wrapf(err, "Scan [qGetShipments, %d]", txID)
			continue
		}
		index := mapIndex[date]
		shipment[index].Items = append(shipment[index].Items, item)
	}
	rows.Close()
	return
}

const (
	qGetShipmentByDate = `SELECT
	s.id AS id,
	s.i_id AS i_id,
	i.item AS code,
	s.amount AS amount,
	COALESCE(DATE_FORMAT(s.deadline, '%d/%m/%Y'), '') AS deadline
FROM
	shipment s JOIN snapshot_item i ON s.i_id = i.id
WHERE s.t_id ? AND s.date = ?
`
	qGetLastShipmentDate = `SELECT DISTINCT
	date
FROM
	shipment
WHERE
	t_id = ?
ORDER BY date DESC
LIMIT 1
`
)

// GetLastShipment last shipment in a transaction
func (d *Database) GetLastShipment(ctx context.Context, txID int64) (shipment model.Shipment, err error) {
	var date time.Time
	if err = d.DB.QueryRowxContext(ctx, qGetLastShipmentDate, txID).Scan(&date); err != nil {
		err = errors.Wrapf(err, "QueryRowxContext [%d]", txID)
		return
	}
	shipment = model.Shipment{
		Date:  date.Format(model.DateFormat),
		Items: make([]model.ShipmentItem, 0),
	}
	if err = d.DB.SelectContext(ctx, &shipment.Items, qGetShipmentByDate, txID, date); err != nil {
		err = errors.Wrapf(err, "SelectContext [%d, %s]", txID, date)
		return
	}
	return
}

const (
	qInsertShipment = `INSERT
INTO 
	shipment (
		t_id,
		date,
		i_id,
		amount,
		deadline,
		modified_by
	)
VALUES ( ?, ?, ?, ?, ?, ? )
`
	qDeleteShipmentByDate = `DELETE
FROM
	shipment
WHERE t_id ? AND date = ?
`
)

// DeleteInsertShipment if exist delete then add shipment in a transaction
func (d *Database) DeleteInsertShipment(ctx context.Context, txID int64, shipment model.Shipment, isDelete bool, actionBy int64) (err error) {
	var date time.Time
	if date, err = time.Parse(model.DateFormat, shipment.Date); err != nil {
		return
	}
	var tx *sqlx.Tx
	if tx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	if isDelete {
		if _, err = tx.ExecContext(ctx, qDeleteShipmentByDate,
			txID,
			date,
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [qDeleteShipmentByDate, %d]", txID)
			_ = tx.Rollback()
			return
		}
	}
	for _, item := range shipment.Items {
		var deadline time.Time
		if deadline, err = time.Parse(model.DateFormat, item.Deadline); err != nil {
			_ = tx.Rollback()
			return
		}
		if _, err = tx.ExecContext(ctx, qInsertPayment,
			txID,
			date,
			item.ItemID,
			item.Amount,
			deadline,
			NullInt64(actionBy),
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [qInsertPayment, %d]", txID)
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
