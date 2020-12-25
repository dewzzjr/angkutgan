package model

// Item is model for Barang
type Item struct {
	Code string `json:"code" db:"code"`
	Name string `json:"name" db:"name"`
	Unit string `json:"unit" db:"unit"`

	Price struct {
		Sell int         `json:"sell"`
		Rent []PriceRent `json:"rent"`
	} `json:"price,omitempty" db:"-"`
}

// PriceRent is model for Harga Sewa Barang
type PriceRent struct {
	ID          int64    `json:"id" db:"id"`
	Description string   `json:"desc" db:"description"`
	Duration    int      `json:"duration" db:"duration"`
	TimeUnit    RentUnit `json:"unit" db:"time_unit"`
	Value       int      `json:"value" db:"value"`
}

// A RentUnit specified unit used for rent (Week = 1, Month = 2).
type RentUnit int

// Rent unit type
const (
	Week RentUnit = 1 + iota
	Month
)

func (r RentUnit) String() string {
	return [...]string{
		"Minggu",
		"Bulan",
	}[r-1]
}
