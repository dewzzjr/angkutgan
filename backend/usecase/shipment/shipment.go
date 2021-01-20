package shipment

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetByTransactionID shipment from the transaction
func (s *Shipment) GetByTransactionID(ctx context.Context, txID int64) (ship []model.Shipment, err error) {
	if ship, err = s.database.GetShipments(ctx, txID); err != nil {
		err = errors.Wrap(err, "GetShipments")
		return
	}
	return
}

// Add shipment to the transaction by date
func (s *Shipment) Add(ctx context.Context, txID int64, ship model.Shipment, actionBy int64) (err error) {
	var items []model.SnapshotItem
	if items, err = s.database.GetSnapshotItems(ctx, txID); err != nil {
		err = errors.Wrap(err, "GetSnapshotItems")
		return
	}
	if err = (&ship).Validate(items); err != nil {
		return
	}
	if err = s.database.DeleteInsertShipment(ctx, txID, ship, false, actionBy); err != nil {
		err = errors.Wrap(err, "DeleteInsertShipment")
		return
	}
	return
}

// Edit the shipment in transaction by date
func (s *Shipment) Edit(ctx context.Context, txID int64, ship model.Shipment, actionBy int64) (err error) {
	var items []model.SnapshotItem
	if items, err = s.database.GetSnapshotItems(ctx, txID); err != nil {
		err = errors.Wrap(err, "GetSnapshotItems")
		return
	}
	date, _ := time.Parse(model.DateFormat, ship.Date)
	var get model.Shipment
	if get, err = s.database.GetShipmentByDate(ctx, txID, date); err != nil {
		err = errors.Wrap(err, "GetShipmentByDate")
		return
	}
	if err = (&ship).Validate(items, get.Items...); err != nil {
		return
	}
	if err = s.database.DeleteInsertShipment(ctx, txID, ship, true, actionBy); err != nil {
		err = errors.Wrap(err, "DeleteInsertShipment")
		return
	}
	return
}

// Delete the shipment in transaction by date
func (s *Shipment) Delete(ctx context.Context, txID int64, date time.Time, actionBy int64) (err error) {
	if err = s.database.DeleteInsertShipment(ctx, txID, model.Shipment{
		Date: date.Format(model.DateFormat),
	}, true, actionBy); err != nil {
		err = errors.Wrap(err, "DeleteInsertShipment")
		return
	}
	return
}
