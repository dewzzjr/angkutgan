package sales

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetDetail by customer code and transaction date
func (i *Sales) GetDetail(ctx context.Context, code string, date time.Time) (tx model.Transaction, err error) {
	if tx, err = i.database.GetTransaction(ctx, date, code, model.Sales); err != nil {
		err = errors.Wrap(err, "GetTransaction")
		return
	}
	if tx.Payment, err = i.payments.GetByTransactionID(ctx, tx.ID); err != nil {
		err = errors.Wrap(err, "GetByTransactionID")
		return
	}
	if tx.Shipment, err = i.shipment.GetByTransactionID(ctx, tx.ID); err != nil {
		err = errors.Wrap(err, "GetByTransactionID")
		return
	}
	return
}

// CreateTransaction sales transaction
func (i *Sales) CreateTransaction(ctx context.Context, tx model.CreateTransaction, actionBy int64) (err error) {
	var date time.Time
	if date, err = time.Parse(model.DateFormat, tx.Date); err != nil {
		return
	}
	var txID int64
	if txID, err = i.database.GetTransactionID(ctx, date, tx.Customer, model.Sales); err != nil {
		err = errors.Wrap(err, "GetTransactionID")
		return
	}
	if txID != 0 {
		err = errors.New("transaksi sudah dibuat")
		return
	}
	if err = (&tx).Calculate(model.Sales); err != nil {
		return
	}
	if err = i.database.InsertTransaction(ctx, model.Sales, tx, actionBy); err != nil {
		err = errors.Wrap(err, "InsertTransaction")
		return
	}
	return
}

// EditTransaction sales transaction
func (i *Sales) EditTransaction(ctx context.Context, tx model.CreateTransaction, actionBy int64) (err error) {
	var date time.Time
	if date, err = time.Parse(model.DateFormat, tx.Date); err != nil {
		return
	}
	var txID int64
	if txID, err = i.database.GetTransactionID(ctx, date, tx.Customer, model.Sales); err != nil {
		err = errors.Wrap(err, "GetTransactionID")
		return
	}
	if txID == 0 {
		err = errors.New("transaksi belum dibuat")
		return
	}
	if err = (&tx).Calculate(model.Sales); err != nil {
		return
	}
	if err = i.database.UpdateTransaction(ctx, txID, model.Sales, tx, actionBy); err != nil {
		err = errors.Wrap(err, "UpdateTransaction")
		return
	}
	return
}
