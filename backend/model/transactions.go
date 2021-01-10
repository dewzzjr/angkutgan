package model

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

// Payment is model for Pembayaran
type Payment struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Date        string        `json:"date"`
	Amount      int           `json:"amount"`
	Method      PaymentMethod `json:"method"`
	MethodDesc  string        `json:"method_desc"`
	Account     AccountType   `json:"account"`
	AccountDesc string        `json:"account_desc"`
	AcceptBy    string        `json:"accept_by"`
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
