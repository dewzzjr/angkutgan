package model

// Item is model for Barang
type Item struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Unit string `json:"unit"`

	Price struct {
		Sell int `json:"sell"`
		Rent []struct {
			Description string   `json:"desc"`
			Duration    int      `json:"duration"`
			TimeUnit    RentUnit `json:"unit"`
			Value       int      `json:"value"`
		} `json:"rent"`
	} `json:"price,omitempty"`
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
