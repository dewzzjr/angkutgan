package model

import (
	"time"

	"github.com/pkg/errors"
)

// Transaction is model for Transaksi Penjualan/Persewaan
type Transaction struct {
	ID       int64    `json:"id"`
	Date     string   `json:"tx_date"`
	Customer Customer `json:"customer"`
	Snapshot
	PaidDate    string     `json:"paid_date"`
	DoneDate    string     `json:"done_date"`
	Payment     []Payment  `json:"payment"`
	Shipment    []Shipment `json:"shipment"`
	Return      []Return   `json:"return,omitempty"`
	Status      TxStatus   `json:"status"`
	DateSummary TxDate     `json:"date"`
}

// Snapshot is model for extend Transaction
type Snapshot struct {
	Address     string         `json:"address"`
	ProjectID   int64          `json:"project_id,omitempty"`
	ProjectName string         `json:"project_name,omitempty"`
	TotalPrice  int            `json:"total_price"`
	Deposit     int            `json:"deposit"`
	Discount    int            `json:"discount"`
	ShippingFee int            `json:"shipping_fee"`
	Items       []SnapshotItem `json:"items"`
}

// SnapshotItem is model for Barang dalam Transaksi
type SnapshotItem struct {
	// GET and POST
	Code         string   `json:"code" db:"item"`
	Amount       int      `json:"amount" db:"amount"`
	Price        int      `json:"price" db:"price"`
	TimeUnit     RentUnit `json:"time_unit,omitempty" db:"time_unit"`
	TimeUnitDesc string   `json:"time_unit_desc,omitempty" db:"-"`
	Duration     int      `json:"duration,omitempty" db:"duration"`
	PreviousID   int64    `json:"previous_id,omitempty" db:"previous_id"`
	// GET only
	ID           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Unit         string `json:"item_unit" db:"item_unit"`
	NeedShipment int    `json:"need_shipment,omitempty" db:"need_shipment"`
	ExtendAmount int    `json:"extend_amount,omitempty" db:"extend_amount"`
	// POST on specific action
	Claim int `json:"claim,omitempty" db:"claim"`
}

// CreateTransaction payload to create transaction
type CreateTransaction struct {
	Date        string         `json:"date"`
	Customer    string         `json:"customer"`
	ProjectID   int64          `json:"project_id"`
	Address     string         `json:"address"`
	Deposit     int            `json:"deposit"`
	Discount    int            `json:"discount"`
	ShippingFee int            `json:"shipping_fee"`
	Items       []SnapshotItem `json:"items"`
	TotalPrice  int            `json:"-"`
}

// Calculate total price transaction by type
func (tx *CreateTransaction) Calculate(txType TransactionType) (err error) {
	var total int
	if tx.Customer == "" {
		err = errors.New("kode pelanggan kosong")
		return
	}
	switch txType {
	case Sales:
		tx.Deposit = 0
		for i, item := range tx.Items {
			if item.Code == "" {
				err = errors.New("kode barang kosong")
				return
			}
			tx.Items[i].Claim = 0
			tx.Items[i].TimeUnit = 0
			tx.Items[i].Duration = 0

			price := (100 - tx.Discount) / 100 * item.Price
			subTotal := item.Amount * price

			tx.Items[i].Price = price
			total += subTotal
		}
	case Rental:
		for i, item := range tx.Items {
			if item.Code == "" {
				err = errors.New("kode barang kosong")
				return
			}
			if !item.TimeUnit.Valid() {
				err = errors.New("satuan waktu tidak valid")
				return
			}
			price := (100 - tx.Discount) / 100 * item.Price
			subTotal := item.Amount * item.Duration * price

			tx.Items[i].Price = price
			total += subTotal
		}
	default:
		err = errors.New("jenis transaksi tidak valid")
		return
	}
	tx.TotalPrice = total
	return
}

// TxStatus is part of TxSummary
type TxStatus struct {
	Description  string `json:"desc"`
	Done         bool   `json:"done"`
	InPayment    bool   `json:"in_payment"`
	PaymentDone  bool   `json:"payment_done"`
	InShipping   bool   `json:"in_shipping"`
	ShippingDone bool   `json:"shipping_done"`
	IsDeadline   bool   `json:"is_deadline"`
	IsReturn     bool   `json:"is_return"`
	IsExtend     bool   `json:"is_extend"`
}

// Desc status description
func (t *TxStatus) Desc(txType TransactionType) {
	switch {
	case t.Done && t.ShippingDone && t.PaymentDone:
		t.Description = "SELESAI"
	case t.Done && t.InPayment && !t.PaymentDone:
		t.Description = "BELUM LUNAS"
	case t.Done && t.IsExtend:
		t.Description = "DIPERPANJANG"
	case !t.Done && t.ShippingDone && t.IsDeadline && !t.IsReturn && txType == Rental:
		t.Description = "BATAS WAKTU"
	case !t.Done && t.ShippingDone && !t.IsDeadline && !t.IsReturn && txType == Rental:
		t.Description = "SEDANG DIPINJAM"
	case !t.Done && t.ShippingDone:
		t.Description = "SELESAI DIKIRIM"
	case !t.Done && t.InShipping && !t.ShippingDone:
		t.Description = "SEDANG DIKIRIM"
	case !t.Done && t.InPayment && !t.InShipping:
		t.Description = "DIBAYAR"
	case !t.Done && !t.InPayment:
		t.Description = "BARU"
	default:
	}
}

// TxDate is part of TxSummary
type TxDate struct {
	Transaction      string `json:"transaction"`
	LastShipmentDate string `json:"last_shipment"`
	LastPaymentDate  string `json:"last_payment"`
	RecentDateline   string `json:"recent_deadline"`
}

// Summary date and status transaction summary
func (t *Transaction) Summary(txType TransactionType) {
	t.DateSummary.Transaction = t.Date
	if len(t.Payment) > 0 {
		t.DateSummary.LastPaymentDate = t.Payment[0].Date
		t.Status.InPayment = true
	}
	if len(t.Shipment) > 0 {
		t.DateSummary.LastShipmentDate = t.Shipment[0].Date
		t.Status.InShipping = true
	}
	if t.PaidDate != "" {
		t.Status.PaymentDone = true
	}
	if t.DoneDate != "" {
		t.Status.ShippingDone = true
	}
	if txType == Rental {
		var needReturn int
		var deadline time.Time
		for _, s := range t.Shipment {
			for _, i := range s.Items {
				needReturn += i.NeedReturn
				if i.NeedReturn > 0 {
					newDeadline, _ := time.Parse(DateFormat, i.Deadline)
					if deadline.IsZero() || deadline.After(newDeadline) {
						deadline = newDeadline
						continue
					}
				}
			}
		}
		if needReturn == 0 {
			t.Status.IsReturn = true
		}
		if !deadline.IsZero() {
			t.DateSummary.RecentDateline = deadline.Format(DateFormat)
			if deadline.After(time.Now()) {
				t.Status.IsDeadline = true
			}
		}
		for _, i := range t.Snapshot.Items {
			if i.ExtendAmount != 0 {
				t.Status.IsExtend = true
				break
			}
		}
	}
	(&t.Status).Desc(txType)
}
