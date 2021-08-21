package transaction

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Transaction usecase object
type Transaction struct {
	database iDatabase
}

// New initiate usecase/transaction
func New(database iDatabase) *Transaction {
	return &Transaction{
		database: database,
	}
}

type iDatabase interface {
	GetListTransactions(ctx context.Context, txType model.TransactionType, date time.Time, limit, offset int) (txs []model.Transaction, err error)
	GetListTransactionsByCustomer(ctx context.Context, customer string, txType model.TransactionType, date time.Time, limit, offset int) (txs []model.Transaction, err error)
}

type iItems interface {
}
