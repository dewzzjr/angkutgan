package transaction

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/package/pagination"
	"github.com/pkg/errors"
)

// GetList sequence of tx using pagination
func (i *Transaction) GetList(ctx context.Context, txType model.TransactionType, date time.Time, page, row int) (txs []model.Transaction, err error) {
	if txs, err = i.database.GetListTransactions(ctx,
		txType,
		date,
		row,
		pagination.Offset(page, row),
	); err != nil {
		err = errors.Wrap(err, "GetListTransactions")
		return
	}
	return
}

// GetByCustomer and date using pagination
func (i *Transaction) GetByCustomer(ctx context.Context, customer string, txType model.TransactionType, date time.Time, page, row int) (txs []model.Transaction, err error) {
	if txs, err = i.database.GetListTransactionsByCustomer(ctx,
		customer,
		txType,
		date,
		row,
		pagination.Offset(page, row),
	); err != nil {
		err = errors.Wrap(err, "GetListTransactions")
		return
	}
	return
}
