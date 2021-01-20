package model

import "github.com/pkg/errors"

// Return is model for Pengembalian
type Return struct {
	Date  string       `json:"date"`
	Items []ReturnItem `json:"item"`
}

// ReturnItem is model for Barang dalam Pengembalian
type ReturnItem struct {
	ID           int64      `json:"id" db:"id"`
	Code         string     `json:"code" db:"code"`
	Status       ItemStatus `json:"status" db:"status"`
	StatusDesc   string     `json:"status_desc" db:"-"`
	Amount       int        `json:"amount" db:"amount"`
	Claim        int        `json:"claim,omitempty" db:"claim"`
	ShipmentID   int64      `json:"shipment_id" db:"s_id"`
	ShipmentDate string     `json:"shipment_date" db:"s_date"`
}

// Validate return with snapshot item
func (r *Return) Validate(ship []Shipment, old ...ReturnItem) (err error) {
	mapItem := make(map[int64]*ShipmentItem)
	for _, s := range ship {
		for _, item := range s.Items {
			for _, o := range old {
				if item.ID == o.ShipmentID {
					item.NeedReturn += o.Amount
				}
			}
			mapItem[item.ID] = &item
		}
	}
	for _, ret := range r.Items {
		if mapItem[ret.ShipmentID] == nil || ret.Amount > mapItem[ret.ShipmentID].NeedReturn {
			err = errors.Errorf("pengembalian [%s] melebihi jumlah barang dalam pengiriman", ret.Code)
			return
		}
		mapItem[ret.ShipmentID].NeedReturn -= ret.Amount
	}
	return
}
