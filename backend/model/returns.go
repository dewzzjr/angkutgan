package model

// Return is model for Pengembalian
type Return struct {
	Date  string       `json:"date"`
	Items []ReturnItem `json:"item"`
}

// ReturnItem is model for Barang dalam Pengembalian
type ReturnItem struct {
	ID           int64      `json:"id" db:"id"`
	ShipmentID   int64      `json:"shipment_id" db:"s_id"`
	ShipmentDate string     `json:"shipment_date" db:"s_date"`
	Code         string     `json:"code" db:"code"`
	Status       ItemStatus `json:"status" db:"status"`
	StatusDesc   string     `json:"status_desc,omitempty" db:"-"`
	Amount       int        `json:"amount" db:"amount"`
	Claim        int        `json:"claim,omitempty" db:"claim"`
}
