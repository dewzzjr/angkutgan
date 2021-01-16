package model

// Shipment is model for Pengiriman
type Shipment struct {
	Date  string         `json:"date"`
	Items []ShipmentItem `json:"items"`
}

// ShipmentItem is model for Barang dalam Pengiriman
type ShipmentItem struct {
	ID       int64  `json:"id" db:"id"`
	ItemID   int64  `json:"item_id" db:"i_id"`
	Code     string `json:"code" db:"code"`
	Amount   int    `json:"amount" db:"amount"`
	Deadline string `json:"daadline" db:"deadline"`
}
