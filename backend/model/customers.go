package model

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
