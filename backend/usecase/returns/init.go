package returns

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Returns usecase object
type Returns struct {
	database iDatabase
	shipment iShipment
}

// New initiate usecase/return
func New(database iDatabase, shipment iShipment) *Returns {
	return &Returns{
		database: database,
		shipment: shipment,
	}
}

type iDatabase interface {
	GetReturns(ctx context.Context, txID int64) (returns []model.Return, err error)
	GetReturnByDate(ctx context.Context, txID int64, date time.Time) (returns model.Return, err error)
	DeleteInsertReturn(ctx context.Context, txID int64, returns model.Return, isDelete bool, actionBy int64) (err error)
}

type iShipment interface {
	GetByTransactionID(ctx context.Context, txID int64) (ship []model.Shipment, err error)
}
