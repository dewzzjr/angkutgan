package model

import "github.com/pkg/errors"

// DateFormat date format standarization
const DateFormat = "02/01/2006"

// CustomerType specified type for customer (Individu = 1, Group = 2).
type CustomerType int

// Customer type
const (
	Individu CustomerType = 1 + iota
	Group
)

func (c CustomerType) String() string {
	return [...]string{
		"Perorangan",
		"Perusahaan",
	}[c-1]
}

// Validate customer
func (c *Customer) Validate() (err error) {
	switch c.Type {
	case Individu:
		c.Role = ""
		c.GroupName = ""
		c.Projects = nil
	case Group:
	default:
		err = errors.Errorf("tipe pelanggan salah")
	}
	return
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

// TransactionType specified type for transaction (Sales = 1, Rental = 2).
type TransactionType int

// Transaction type
const (
	Sales CustomerType = 1 + iota
	Rental
)

func (c TransactionType) String() string {
	return [...]string{
		"Penjualan",
		"Persewaan",
	}[c-1]
}

// PaymentMethod type
type PaymentMethod int

// Payment method type
const (
	Cash     PaymentMethod = 100
	Transfer PaymentMethod = 200
)

func (i PaymentMethod) String() string {
	return map[PaymentMethod]string{
		Cash:     "TUNAI",
		Transfer: "TRANSFER",
	}[i]
}

// AccountType type
type AccountType int

// Account type
const (
	Debit  AccountType = 100
	Credit AccountType = 200
)

func (i AccountType) String() string {
	return map[AccountType]string{
		Debit:  "DEBIT",
		Credit: "KREDIT",
	}[i]
}

// ItemStatus condition
type ItemStatus int

// Item status type
const (
	Restock   ItemStatus = 100
	Sold      ItemStatus = 200
	Good      ItemStatus = 300
	LowRepair ItemStatus = 302
	MidRepair ItemStatus = 305
	LowBroken ItemStatus = 402
	MidBroken ItemStatus = 405
	Lost      ItemStatus = 410
)

func (i ItemStatus) String() string {
	return map[ItemStatus]string{
		Restock:   "Stok Ulang",
		Sold:      "Terjual",
		Good:      "Baik",
		LowRepair: "Perbaikan Ringan",
		MidRepair: "Perbaikan Sedang",
		LowBroken: "Rusak Ringan",
		MidBroken: "Rusak Sedang",
		Lost:      "Hilang/Rusak Parah",
	}[i]
}
