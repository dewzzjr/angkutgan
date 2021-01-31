package ajax

import (
	"context"
	"fmt"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// IsValidCustomerCode check if code either useable or not
func (u *Ajax) IsValidCustomerCode(ctx context.Context, code string) (ok bool, err error) {
	return u.database.IsValidCustomerCode(ctx, code)
}

// GetCustomers get autocomplete for customer
func (u *Ajax) GetCustomers(ctx context.Context, keyword string) (values []model.AutoComplete, err error) {
	var customers []model.Customer
	if customers, err = u.database.GetListCustomersByKeyword(ctx, keyword, model.MaxAutoComplete, 0); err != nil {
		err = errors.Wrap(err, "GetListCustomersByKeyword")
		return
	}
	values = make([]model.AutoComplete, 0)
	for _, i := range customers {
		var name string
		if i.GroupName != "" {
			name = i.GroupName
		} else {
			name = i.Name
		}
		values = append(values, model.AutoComplete{
			Value: i.Code,
			Text:  fmt.Sprintf("%s - %s", i.Code, name),
		})
	}
	return
}
