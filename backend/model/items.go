package model

// Item is model for Barang
type Item struct {
	Code  string `json:"code" db:"code"`
	Name  string `json:"name" db:"name"`
	Unit  string `json:"unit" db:"unit"`
	Stock int    `json:"stock" db:"-"`

	Price struct {
		Sell int         `json:"sell"`
		Rent []PriceRent `json:"rent"`
	} `json:"price,omitempty" db:"-"`

	Available Stock `json:"avail,omitempty" db:"-"`
}

// PriceRent is model for Harga Sewa Barang
type PriceRent struct {
	ID           int64    `json:"id" db:"id"`
	Description  string   `json:"desc" db:"description"`
	Duration     int      `json:"duration" db:"duration"`
	TimeUnit     RentUnit `json:"unit" db:"time_unit"`
	TimeUnitDesc string   `json:"unit_desc" db:"-"`
	Value        int      `json:"value" db:"value"`
}

// Stock availability
type Stock struct {
	Inventory int `json:"inventory" db:"inventory"`
	Asset     int `json:"asset" db:"asset"`
}
