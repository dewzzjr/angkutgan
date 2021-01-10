package sales

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetSalesDetail by customer code and transaction date
func (i *Sales) GetSalesDetail(ctx context.Context, code string, date time.Time) (tx model.Transaction, err error) {
	// TODO get transaction, snapshot, snapshot_item
	if tx.Payment, err = i.payments.GetPayments(ctx, tx.ID); err != nil {
		err = errors.Wrap(err, "GetPayments")
		return
	}
	if tx.Shipment, err = i.shipment.GetShipments(ctx, tx.ID); err != nil {
		err = errors.Wrap(err, "GetShipments")
		return
	}
	return
}

// CreateTransaction sales transaction
func (i *Sales) CreateTransaction(ctx context.Context, tx model.CreateTransaction, actionBy int64) (err error) {
	if _, err = time.Parse(model.DateFormat, tx.Date); err != nil {
		return
	}
	var items []string
	var total int
	for i, item := range tx.Items {
		discountedPrice := (100 - tx.Discount) / 100 * item.Price
		tx.Items[i].Claim = 0
		tx.Items[i].TimeUnit = 0
		tx.Items[i].Duration = 0
		tx.Items[i].Price = discountedPrice
		items = append(items, item.Code)
		total += item.Amount * discountedPrice
	}
	// TODO insert transaction, snapshot, snapshot_item
	_, _, _ = model.Sales, tx, total
	return
}
