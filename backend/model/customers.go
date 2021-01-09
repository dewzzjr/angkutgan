package model

import "github.com/pkg/errors"

// Customer is model for Pelanggan
type Customer struct {
	Code      string       `json:"code" db:"code"`
	Name      string       `json:"name" db:"name"`
	Type      CustomerType `json:"type" db:"type"`
	Address   string       `json:"address" db:"address"`
	Phone     string       `json:"phone" db:"phone"`
	NIK       string       `json:"nik" db:"nik"`
	GroupName string       `json:"group_name" db:"group_name"`
	Role      string       `json:"role" db:"role"`
	Projects  []Project    `json:"project,omitempty" db:"-"`
}

// Project is model for Proyek Perusahaan
type Project struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Location string `json:"location" db:"location"`
}

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
