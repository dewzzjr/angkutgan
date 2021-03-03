package rental

import (
	"context"
	"strings"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetDetail by customer code and transaction date
func (i *Rental) GetDetail(ctx context.Context, code string, date time.Time) (tx model.Transaction, err error) {
	if tx, err = i.database.GetTransaction(ctx, date, code, model.Rental); err != nil {
		err = errors.Wrap(err, "GetTransaction")
		return
	}
	// empty struct when not found
	if tx.ID == 0 {
		return
	}
	tx, err = i.completeTx(ctx, tx)
	return
}

// GetByCustomer by customer code and transaction date
func (i *Rental) GetByCustomer(ctx context.Context, page, row int, customer string, date time.Time) (txs []model.Transaction, err error) {
	if txs, err = i.transaction.GetByCustomer(ctx,
		strings.TrimSpace(customer),
		model.Rental,
		date,
		page,
		row,
	); err != nil {
		err = errors.Wrap(err, "GetByCustomer")
		return
	}

	for t, tx := range txs {
		if txs[t].Payment, err = i.payments.GetByTransactionID(ctx, tx.ID); err != nil {
			err = errors.Wrap(err, "GetByTransactionID")
			return
		}
		if txs[t].Shipment, err = i.shipment.GetByTransactionID(ctx, tx.ID); err != nil {
			err = errors.Wrap(err, "GetByTransactionID")
			return
		}
		if txs[t].Return, err = i.returns.GetByTransactionID(ctx, tx.ID); err != nil {
			err = errors.Wrap(err, "GetByTransactionID")
			return
		}
		(&(txs[t])).Summary(model.Rental)
	}
	return
}

// GetList sales using pagination
func (i *Rental) GetList(ctx context.Context, page, row int, date time.Time) (txs []model.Transaction, err error) {
	if txs, err = i.transaction.GetList(ctx,
		model.Rental,
		date,
		page,
		row,
	); err != nil {
		err = errors.Wrap(err, "GetList")
		return
	}

	txs, err = i.bulkTx(ctx, txs)
	return
}

// CreateTransaction rental transaction
func (i *Rental) CreateTransaction(ctx context.Context, tx model.CreateTransaction, actionBy int64) (err error) {
	var date time.Time
	if date, err = time.Parse(model.DateFormat, tx.Date); err != nil {
		return
	}
	var txID int64
	if txID, err = i.database.GetTransactionID(ctx, date, tx.Customer, model.Rental); err != nil {
		err = errors.Wrap(err, "GetTransactionID")
		return
	}
	if txID != 0 {
		err = errors.New("transaksi sudah dibuat")
		return
	}
	if err = (&tx).Calculate(model.Rental); err != nil {
		return
	}
	if err = i.database.InsertTransaction(ctx, model.Rental, tx, actionBy); err != nil {
		err = errors.Wrap(err, "InsertTransaction")
		return
	}
	return
}

// ExtendTransaction previous transaction
func (i *Rental) ExtendTransaction(ctx context.Context, tx model.CreateTransaction, actionBy int64) (err error) {
	var date time.Time
	if date, err = time.Parse(model.DateFormat, tx.Date); err != nil {
		return
	}
	var txID int64
	if txID, err = i.database.GetTransactionID(ctx, date, tx.Customer, model.Rental); err != nil {
		err = errors.Wrap(err, "GetTransactionID")
		return
	}
	if txID != 0 {
		err = errors.New("transaksi sudah dibuat")
		return
	}
	if err = (&tx).Calculate(model.Rental); err != nil {
		return
	}
	if err = i.database.InsertTransaction(ctx, model.Rental, tx, actionBy); err != nil {
		err = errors.Wrap(err, "InsertTransaction")
		return
	}
	return
}

// EditTransaction sales transaction
func (i *Rental) EditTransaction(ctx context.Context, tx model.CreateTransaction, actionBy int64) (err error) {
	var date time.Time
	if date, err = time.Parse(model.DateFormat, tx.Date); err != nil {
		return
	}
	var txID int64
	if txID, err = i.database.GetTransactionID(ctx, date, tx.Customer, model.Rental); err != nil {
		err = errors.Wrap(err, "GetTransactionID")
		return
	}
	if txID == 0 {
		err = errors.New("transaksi belum dibuat")
		return
	}
	if err = (&tx).Calculate(model.Rental); err != nil {
		return
	}
	if err = i.database.UpdateTransaction(ctx, txID, model.Rental, tx, actionBy); err != nil {
		err = errors.Wrap(err, "UpdateTransaction")
		return
	}
	return
}
