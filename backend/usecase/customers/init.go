package customers

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Customers usecase object
type Customers struct {
	database iDatabase
}

// New initiate usecase/customer
func New(database iDatabase) *Customers {
	return &Customers{
		database: database,
	}
}

type iDatabase interface {
	GetListCustomers(ctx context.Context, limit, offset int) (customers []model.Customer, err error)
	GetListCustomersByKeyword(ctx context.Context, keyword string, limit, offset int, column ...string) (customers []model.Customer, err error)
	GetCustomerDetail(ctx context.Context, code string) (customer model.Customer, err error)
	UpdateInsertCustomer(ctx context.Context, customer model.Customer, actionBy int64) (err error)
	DeleteCustomer(ctx context.Context, code string) (err error)
	InsertDeleteProject(ctx context.Context, code string, insert []model.Project, delete []int64, actionBy int64) (err error)
}
