package rental

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetDetail by customer code and transaction date
func (i *Rental) completeTx(ctx context.Context, in model.Transaction) (out model.Transaction, err error) {
	defer func() {
		out = in
	}()
	if in.Payment, err = i.payments.GetByTransactionID(ctx, in.ID); err != nil {
		err = errors.Wrap(err, "GetByTransactionID")
		return
	}
	if in.Shipment, err = i.shipment.GetByTransactionID(ctx, in.ID); err != nil {
		err = errors.Wrap(err, "GetByTransactionID")
		return
	}
	if in.Return, err = i.returns.GetByTransactionID(ctx, in.ID); err != nil {
		err = errors.Wrap(err, "GetByTransactionID")
		return
	}
	(&in).Summary(model.Rental)
	return
}

type goStruct struct {
	Index int
	Tx    model.Transaction
	Error error
}

func (i *Rental) worker(ctx context.Context, jobs <-chan goStruct, results chan<- goStruct) {
	for j := range jobs {
		tx, e := i.completeTx(ctx, j.Tx)
		results <- goStruct{j.Index, tx, e}
	}
}

func (i *Rental) bulkTx(ctx context.Context, in []model.Transaction) (out []model.Transaction, err error) {
	defer func() {
		out = in
	}()
	numJobs := len(in)
	maxWorkers := model.MaxWorkers
	jobs := make(chan goStruct, numJobs)
	results := make(chan goStruct, numJobs)
	for n := 0; n < maxWorkers; n++ {
		go i.worker(ctx, jobs, results)
	}
	for index, tx := range in {
		jobs <- goStruct{Index: index, Tx: tx}
	}
	close(jobs)
	for range in {
		r := <-results
		if r.Error != nil {
			err = r.Error
			return
		}
		in[r.Index] = r.Tx
	}
	return
}
