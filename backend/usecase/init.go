package usecase

import (
	"github.com/dewzzjr/angkutgan/backend/package/config"
	"github.com/dewzzjr/angkutgan/backend/repository"
	"github.com/dewzzjr/angkutgan/backend/usecase/ajax"
	"github.com/dewzzjr/angkutgan/backend/usecase/customers"
	"github.com/dewzzjr/angkutgan/backend/usecase/items"
	"github.com/dewzzjr/angkutgan/backend/usecase/users"
)

// Usecase object
type Usecase struct {
	Items     *items.Items
	Users     *users.Users
	Customers *customers.Customers
	Ajax      *ajax.Ajax
}

// New initiate usecase
func New(r *repository.Repository) *Usecase {
	cfg := config.Get()
	return &Usecase{
		Items:     items.New(r.Database),
		Users:     users.New(r.Database, cfg.Users),
		Customers: customers.New(r.Database),
		Ajax:      ajax.New(r.Database),
	}
}
