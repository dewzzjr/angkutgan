package model

import (
	"time"

	"github.com/pkg/errors"
)

// Shipment is model for Pengiriman
type Shipment struct {
	Date  string         `json:"date"`
	Items []ShipmentItem `json:"items"`
}

// ShipmentItem is model for Barang dalam Pengiriman
type ShipmentItem struct {
	ID         int64  `json:"id" db:"id"`
	ItemID     int64  `json:"item_id" db:"i_id"`
	Code       string `json:"code" db:"code"`
	Amount     int    `json:"amount" db:"amount"`
	Deadline   string `json:"deadline" db:"deadline"`
	NeedReturn int    `json:"need_return,omitempty" db:"need_return"`
}

// Validate shipment with snapshot item
func (s *Shipment) Validate(items []SnapshotItem, old ...ShipmentItem) (err error) {
	mapItem := make(map[int64]*SnapshotItem)
	for _, item := range items {
		for _, o := range old {
			if item.ID == o.ItemID {
				item.NeedShipment += o.Amount
				break
			}
		}
		mapItem[item.ID] = &item
	}
	for i, ship := range s.Items {
		if mapItem[ship.ItemID] == nil || ship.Amount > mapItem[ship.ItemID].NeedShipment {
			err = errors.Errorf("pengiriman [%s] melebihi jumlah barang dalam transaksi", ship.Code)
			return
		}
		mapItem[ship.ItemID].NeedShipment -= ship.Amount
		if mapItem[ship.ItemID].Duration > 0 {
			s.Items[i].Deadline = calculateDuration(
				s.Date,
				mapItem[ship.ItemID].Duration,
				mapItem[ship.ItemID].TimeUnit,
			)
		}
	}
	return
}

func calculateDuration(start string, duration int, timeunit RentUnit) string {
	startTime, err := time.Parse(DateFormat, start)
	if err != nil {
		return ""
	}
	switch timeunit {
	case Month:
		return startTime.AddDate(0, duration, 0).Format(DateFormat)
	case Week:
		return startTime.AddDate(0, 0, 7*duration).Format(DateFormat)
	default:
		return ""
	}
}
