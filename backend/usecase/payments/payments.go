package payments

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetPayments from the transaction
func (p *Payments) GetPayments(ctx context.Context, txID int64) (pay []model.Payment, err error) {
	if pay, err = p.database.GetPayments(ctx, txID); err != nil {
		err = errors.Wrap(err, "GetPayments")
		return
	}
	for i, p := range pay {
		pay[i].AccountDesc = p.Account.String()
		pay[i].MethodDesc = p.Method.String()
	}
	return
}

// AddPayment to the transaction
func (p *Payments) AddPayment(ctx context.Context, txID int64, pay model.Payment, actionBy int64) (err error) {
	if err = p.database.InsertPayment(ctx, txID, pay, actionBy); err != nil {
		err = errors.Wrap(err, "InsertPayment")
		return
	}
	// TODO update transaction payment date
	return
}

// EditPayment from the last payment in transaction
func (p *Payments) EditPayment(ctx context.Context, txID int64, pay model.Payment, actionBy int64) (err error) {
	if err = p.database.UpdatePayment(ctx, txID, pay, actionBy); err != nil {
		err = errors.Wrap(err, "UpdatePayment")
		return
	}
	// TODO update transaction payment date
	return
}

// DeletePayment from the last payment in transaction
func (p *Payments) DeletePayment(ctx context.Context, txID int64) (err error) {
	if err = p.database.DeletePayment(ctx, txID); err != nil {
		err = errors.Wrap(err, "DeletePayment")
		return
	}
	// TODO update transaction payment date
	return
}

// UpdateTxPaidDate update the paid date based on payment done and total payment of transaction
func (p *Payments) UpdateTxPaidDate(ctx context.Context, txID int64) (err error) {
	return
}
