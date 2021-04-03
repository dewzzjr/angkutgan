package payments

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetByTransactionID payments from the transaction
func (p *Payments) GetByTransactionID(ctx context.Context, txID int64) (pay []model.Payment, err error) {
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

// Add payment to the transaction
func (p *Payments) Add(ctx context.Context, txID int64, pay model.Payment, actionBy int64) (err error) {
	if err = p.database.InsertPayment(ctx, txID, pay, actionBy); err != nil {
		err = errors.Wrap(err, "InsertPayment")
		return
	}
	if err = p.UpdateTxPaidDate(ctx, txID); err != nil {
		err = errors.Wrap(err, "UpdateTxPaidDate")
		return
	}
	return
}

// Edit the last payment in transaction
func (p *Payments) Edit(ctx context.Context, txID int64, pay model.Payment, actionBy int64) (err error) {
	if err = p.database.UpdatePayment(ctx, txID, pay, actionBy); err != nil {
		err = errors.Wrap(err, "UpdatePayment")
		return
	}
	if err = p.UpdateTxPaidDate(ctx, txID); err != nil {
		err = errors.Wrap(err, "UpdateTxPaidDate")
		return
	}
	return
}

// Delete the last payment in transaction
func (p *Payments) Delete(ctx context.Context, txID int64) (err error) {
	if err = p.database.DeletePayment(ctx, txID); err != nil {
		err = errors.Wrap(err, "DeletePayment")
		return
	}
	if err = p.UpdateTxPaidDate(ctx, txID); err != nil {
		err = errors.Wrap(err, "UpdateTxPaidDate")
		return
	}
	return
}

// UpdateTxPaidDate update the paid date based on payment done and total payment of transaction
func (p *Payments) UpdateTxPaidDate(ctx context.Context, txID int64) (err error) {
	var paid, bill int
	if bill, err = p.database.GetTotalPayment(ctx, txID); err != nil {
		err = errors.Wrap(err, "GetTotalPayment")
		return
	}
	var date string
	if paid, date, err = p.database.GetLastPaidAmount(ctx, txID); err != nil {
		err = errors.Wrap(err, "GetLastPaidAmount")
		return
	}
	var paidDate time.Time
	if paid >= bill {
		if paidDate, err = time.Parse(model.DateFormat, date); err != nil {
			return
		}
	}
	if err = p.database.UpdatePaidDate(ctx, txID, paidDate); err != nil {
		err = errors.Wrap(err, "UpdatePaidDate")
		return
	}
	return
}
