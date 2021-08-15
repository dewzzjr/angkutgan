package sales

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Sales usecase object
type Sales struct {
	database    iDatabase
	payments    iPayments
	shipment    iShipment
	transaction iTransaction
}

// New initiate usecase/sales
func New(database iDatabase, payments iPayments, shipment iShipment, transaction iTransaction) *Sales {
	return &Sales{
		database:    database,
		payments:    payments,
		shipment:    shipment,
		transaction: transaction,
	}
}

type iDatabase interface {
	GetTransaction(ctx context.Context, date time.Time, code string, txType model.TransactionType) (tx model.Transaction, err error)
	GetTransactionID(ctx context.Context, date time.Time, code string, txType model.TransactionType) (txID int64, err error)
	InsertTransaction(ctx context.Context, txType model.TransactionType, tx model.CreateTransaction, actionBy int64) (err error)
	UpdateTransaction(ctx context.Context, txID int64, txType model.TransactionType, tx model.CreateTransaction, actionBy int64) (err error)
	DeleteTransaction(ctx context.Context, txID int64) (err error)
}

type iShipment interface {
	GetByTransactionID(ctx context.Context, txID int64) (ship []model.Shipment, err error)
}

type iPayments interface {
	GetByTransactionID(ctx context.Context, txID int64) (pay []model.Payment, err error)
}

type iTransaction interface {
	GetList(ctx context.Context, txType model.TransactionType, date time.Time, page, row int) (txs []model.Transaction, err error)
	GetByCustomer(ctx context.Context, customer string, txType model.TransactionType, date time.Time, page, row int) (txs []model.Transaction, err error)
}
