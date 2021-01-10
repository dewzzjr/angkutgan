package sales

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Sales usecase object
type Sales struct {
	database iDatabase
	payments iPayments
	shipment iShipment
}

// New initiate usecase/sales
func New(database iDatabase, payments iPayments, shipment iShipment) *Sales {
	return &Sales{
		database: database,
		payments: payments,
		shipment: shipment,
	}
}

type iDatabase interface {
}

type iShipment interface {
	GetShipments(ctx context.Context, txID int64) (ship []model.Shipment, err error)
}

type iPayments interface {
	GetPayments(ctx context.Context, txID int64) (pay []model.Payment, err error)
}
