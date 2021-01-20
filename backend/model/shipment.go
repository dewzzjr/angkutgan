package model

import (
	"github.com/pkg/errors"
)

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
	Deadline string `json:"deadline" db:"deadline"`
}

// Validate shipment with snapshot item
func (s *Shipment) Validate(items []SnapshotItem, old ...ShipmentItem) (err error) {
	mapItem := make(map[int64]*SnapshotItem)
	for _, item := range items {
		for _, o := range old {
			if item.ID == o.ItemID {
				item.NeedShipment += o.Amount
			}
		}
		mapItem[item.ID] = &item
	}
	for _, ship := range s.Items {
		if mapItem[ship.ItemID] == nil || ship.Amount > mapItem[ship.ItemID].NeedShipment {
			err = errors.Errorf("pengiriman [%s] melebihi jumlah barang dalam transaksi", ship.Code)
			return
		}
		mapItem[ship.ItemID].NeedShipment -= ship.Amount
	}
	return
}
