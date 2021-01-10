package shipment

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// GetShipments from the transaction
func (i *Shipment) GetShipments(ctx context.Context, txID int64) (ship []model.Shipment, err error) {
	return
}

// AddShipment to the transaction
func (i *Shipment) AddShipment(ctx context.Context, txID int64, ship model.Shipment, actionBy int64) (err error) {
	return
}

// EditShipment from the last shipment in transaction
func (i *Shipment) EditShipment(ctx context.Context, txID int64, ship model.Shipment, actionBy int64) (err error) {
	return
}

// DeleteShipment from the last shipment in transaction
func (i *Shipment) DeleteShipment(ctx context.Context, txID int64) (err error) {
	return
}
