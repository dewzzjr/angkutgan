package customers

import (
	"context"
	"strings"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/package/pagination"
	"github.com/pkg/errors"
)

// GetList sequence of customers
func (i *Customers) GetList(ctx context.Context, page int, row int) (customers []model.Customer, err error) {
	if customers, err = i.database.GetListCustomers(ctx, row, pagination.Offset(page, row)); err != nil {
		err = errors.Wrap(err, "GetListCustomers")
	}
	return
}

// GetByKeyword sequence of customers by keyword
func (i *Customers) GetByKeyword(ctx context.Context, page int, row int, key string) (customers []model.Customer, err error) {
	if customers, err = i.database.GetListCustomersByKeyword(ctx,
		strings.TrimSpace(key),
		row,
		pagination.Offset(page, row),
	); err != nil {
		err = errors.Wrap(err, "GetListCustomersByKeyword")
	}
	return
}

// Get customer by code
func (i *Customers) Get(ctx context.Context, code string) (customer model.Customer, err error) {
	if customer, err = i.database.GetCustomerDetail(ctx, code); err != nil {
		err = errors.Wrap(err, "GetCustomerDetail")
	}
	return
}

// Create new customer
func (i *Customers) Create(ctx context.Context, customer model.Customer, actionBy int64) (err error) {
	if err = (&customer).Validate(); err != nil {
		err = errors.Wrap(err, "Validate")
		return
	}
	if err = i.database.UpdateInsertCustomer(ctx, customer, actionBy); err != nil {
		err = errors.Wrap(err, "UpdateInsertCustomer")
		return
	}
	if err = i.database.InsertDeleteProject(ctx, customer.Code, customer.Projects, nil, actionBy); err != nil {
		err = errors.Wrap(err, "InsertDeleteProject")
		return
	}
	return
}

// Update customer by code
func (i *Customers) Update(ctx context.Context, customer model.Customer, actionBy int64) (err error) {
	if err = (&customer).Validate(); err != nil {
		err = errors.Wrap(err, "Validate")
		return
	}
	var get model.Customer
	if get, err = i.Get(ctx, customer.Code); err != nil {
		err = errors.Wrapf(err, "Get")
		return
	}
	if get.Type == model.Group && customer.Type == model.Individu {
		// TODO check is remove project eligible
	}
	if customer.Name != get.Name ||
		customer.Type != get.Type ||
		customer.Address != get.Address ||
		customer.Phone != get.Phone ||
		customer.Role != get.Role ||
		customer.GroupName != get.GroupName ||
		customer.NIK != get.NIK {
		if err = i.database.UpdateInsertCustomer(ctx, customer, actionBy); err != nil {
			err = errors.Wrap(err, "UpdateInsertCustomer")
			return
		}
	}
	if err = i.UpdateProject(ctx, customer.Code, get.Projects, customer.Projects, actionBy); err != nil {
		err = errors.Wrap(err, "UpdateProject")
	}
	return
}

// Remove customer by code
func (i *Customers) Remove(ctx context.Context, code string) (err error) {
	// TODO check is delete eligible
	if err = i.database.DeleteCustomer(ctx, code); err != nil {
		err = errors.Wrap(err, "UpdateProject")
	}
	return
}
