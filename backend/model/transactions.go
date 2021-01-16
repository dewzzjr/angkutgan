package model

import "errors"

// Transaction is model for Transaksi Penjualan/Persewaan
type Transaction struct {
	ID       int64    `json:"id"`
	Date     string   `json:"date"`
	Customer Customer `json:"customer"`
	Snapshot
	Payment  []Payment  `json:"payment"`
	PaidDate string     `json:"paid_date"`
	Shipment []Shipment `json:"shipment"`
	Return   []Return   `json:"return,omitempty"`
	DoneDate string     `json:"done_date"`
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
	ID       int64    `json:"id" db:"id"`
	Code     string   `json:"code" db:"item"`
	Name     string   `json:"name" db:"name"`
	Amount   int      `json:"amount" db:"amount"`
	Price    int      `json:"price" db:"price"`
	Claim    int      `json:"claim,omitempty" db:"claim"`
	TimeUnit RentUnit `json:"time_unit,omitempty" db:"time_unit"`
	Duration int      `json:"duration,omitempty" db:"duration"`
}

// Shipment is model for Pengiriman
type Shipment struct {
	Date  string         `json:"date"`
	Items []ShipmentItem `json:"items"`
}

// Return is model for Pengembalian
type Return struct {
	Date  string       `json:"date"`
	Items []ReturnItem `json:"items"`
}

// ShipmentItem is model for Barang dalam Pengiriman
type ShipmentItem struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Amount   int    `json:"amount"`
	Deadline string `json:"daadline"`
}

// ReturnItem is model for Barang dalam Pengembalian
type ReturnItem struct {
	ID         int64      `json:"id"`
	Code       string     `json:"code"`
	Amount     int        `json:"amount"`
	Status     ItemStatus `json:"status"`
	StatusDesc string     `json:"status_desc"`
	Claim      int        `json:"claim,omitempty"`
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
	switch txType {
	case Sales:
		for i, item := range tx.Items {
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
