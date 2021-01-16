package payments

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Payments usecase object
type Payments struct {
	database iDatabase
}

// New initiate usecase/payments
func New(database iDatabase) *Payments {
	return &Payments{
		database: database,
	}
}

type iDatabase interface {
	GetPayments(ctx context.Context, txID int64) (payments []model.Payment, err error)
	InsertPayment(ctx context.Context, txID int64, payment model.Payment, actionBy int64) (err error)
	UpdatePayment(ctx context.Context, txID int64, payment model.Payment, actionBy int64) (err error)
	DeletePayment(ctx context.Context, txID int64) (err error)

	UpdatePaidDate(ctx context.Context, txID int64, date time.Time) (err error)
}
