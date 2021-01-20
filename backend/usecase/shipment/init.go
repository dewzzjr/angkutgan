package shipment

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Shipment usecase object
type Shipment struct {
	database iDatabase
}

// New initiate usecase/shipment
func New(database iDatabase) *Shipment {
	return &Shipment{
		database: database,
	}
}

type iDatabase interface {
	GetShipments(ctx context.Context, txID int64) (shipment []model.Shipment, err error)
	GetShipmentByDate(ctx context.Context, txID int64, date time.Time) (shipment model.Shipment, err error)
	DeleteInsertShipment(ctx context.Context, txID int64, shipment model.Shipment, isDelete bool, actionBy int64) (err error)

	GetSnapshotItems(ctx context.Context, txID int64) (items []model.SnapshotItem, err error)
}
