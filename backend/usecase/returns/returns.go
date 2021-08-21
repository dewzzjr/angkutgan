package returns

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetByTransactionID return from the transaction
func (r *Returns) GetByTransactionID(ctx context.Context, txID int64) (returns []model.Return, err error) {
	if returns, err = r.database.GetReturns(ctx, txID); err != nil {
		err = errors.Wrap(err, "GetReturns")
		return
	}
	for i, r := range returns {
		for j, item := range r.Items {
			returns[i].Items[j].StatusDesc = item.Status.String()
		}
	}
	return
}

// Add return to the transaction by date
func (r *Returns) Add(ctx context.Context, txID int64, returns model.Return, actionBy int64) (err error) {
	var items []model.Shipment
	if items, err = r.shipment.GetByTransactionID(ctx, txID); err != nil {
		err = errors.Wrap(err, "GetSnapshotItems")
		return
	}
	if err = (&returns).Validate(items); err != nil {
		return
	}
	if err = r.database.DeleteInsertReturn(ctx, txID, returns, false, actionBy); err != nil {
		err = errors.Wrap(err, "DeleteInsertReturn")
		return
	}
	return
}

// Edit the return in transaction by date
func (r *Returns) Edit(ctx context.Context, txID int64, returns model.Return, actionBy int64) (err error) {
	var items []model.Shipment
	if items, err = r.shipment.GetByTransactionID(ctx, txID); err != nil {
		err = errors.Wrap(err, "GetSnapshotItems")
		return
	}
	date, _ := time.Parse(model.DateFormat, returns.Date)
	var get model.Return
	if get, err = r.database.GetReturnByDate(ctx, txID, date); err != nil {
		err = errors.Wrap(err, "GetReturnByDate")
		return
	}
	if err = (&returns).Validate(items, get.Items...); err != nil {
		return
	}
	if err = r.database.DeleteInsertReturn(ctx, txID, returns, true, actionBy); err != nil {
		err = errors.Wrap(err, "DeleteInsertReturn")
		return
	}
	return
}

// Delete the return in transaction by date
func (r *Returns) Delete(ctx context.Context, txID int64, date time.Time, actionBy int64) (err error) {
	if err = r.database.DeleteInsertReturn(ctx, txID, model.Return{
		Date: date.Format(model.DateFormat),
	}, true, actionBy); err != nil {
		err = errors.Wrap(err, "DeleteInsertReturn")
		return
	}
	return
}
