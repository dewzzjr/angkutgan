package sales

import (
	"context"
	"strings"
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
	tx, err = i.completeTx(ctx, tx)
	return
}

// GetByCustomer by customer code and transaction date
func (i *Sales) GetByCustomer(ctx context.Context, page, row int, customer string, date time.Time) (txs []model.Transaction, err error) {
	if txs, err = i.transaction.GetByCustomer(ctx,
		strings.TrimSpace(customer),
		model.Sales,
		date,
		page,
		row,
	); err != nil {
		err = errors.Wrap(err, "GetByCustomer")
		return
	}

	txs, err = i.bulkTx(ctx, txs)
	return
}

// GetList sales using pagination
func (i *Sales) GetList(ctx context.Context, page, row int, date time.Time) (txs []model.Transaction, err error) {
	if txs, err = i.transaction.GetList(ctx,
		model.Sales,
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
