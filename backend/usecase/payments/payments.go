package payments

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// GetPayments from the transaction
func (i *Payments) GetPayments(ctx context.Context, txID int64) (pay []model.Payment, err error) {
	return
}

// AddPayment to the transaction
func (i *Payments) AddPayment(ctx context.Context, txID int64, pay model.Payment, actionBy int64) (err error) {
	return
}

// EditPayment from the last payment in transaction
func (i *Payments) EditPayment(ctx context.Context, txID int64, pay model.Payment, actionBy int64) (err error) {
	return
}

// DeletePayment from the last payment in transaction
func (i *Payments) DeletePayment(ctx context.Context, txID int64) (err error) {
	return
}
