package usecase

import (
	"github.com/dewzzjr/angkutgan/backend/package/config"
	"github.com/dewzzjr/angkutgan/backend/repository"
	"github.com/dewzzjr/angkutgan/backend/usecase/ajax"
	"github.com/dewzzjr/angkutgan/backend/usecase/customers"
	"github.com/dewzzjr/angkutgan/backend/usecase/items"
	"github.com/dewzzjr/angkutgan/backend/usecase/payments"
	"github.com/dewzzjr/angkutgan/backend/usecase/rental"
	"github.com/dewzzjr/angkutgan/backend/usecase/sales"
	"github.com/dewzzjr/angkutgan/backend/usecase/shipment"
	"github.com/dewzzjr/angkutgan/backend/usecase/users"
)

// Usecase object
type Usecase struct {
	Customers *customers.Customers
	Items     *items.Items
	Payments  *payments.Payments
	Shipment  *shipment.Shipment
	Sales     *sales.Sales
	Rental    *rental.Rental
	Users     *users.Users
	Ajax      *ajax.Ajax
}

// New initiate usecase
func New(r *repository.Repository) *Usecase {
	cfg := config.Get()
	u := &Usecase{}
	u.Ajax = ajax.New(r.Database)
	u.Items = items.New(r.Database)
	u.Payments = payments.New(r.Database)
	u.Shipment = shipment.New(r.Database)
	u.Customers = customers.New(r.Database)
	u.Users = users.New(r.Database, cfg.Users)
	u.Sales = sales.New(r.Database, u.Payments, u.Shipment)
	u.Rental = rental.New(r.Database, u.Payments, u.Shipment)
	return u
}
