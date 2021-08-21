package ajax

import (
	"context"
	"fmt"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// IsValidItemCode check if code either useable or not
func (u *Ajax) IsValidItemCode(ctx context.Context, code string) (ok bool, err error) {
	return u.database.IsValidItemCode(ctx, code)
}

// GetItems get autocomplete for item
func (u *Ajax) GetItems(ctx context.Context, keyword string) (values []model.AutoComplete, err error) {
	var items []model.Item
	if items, err = u.database.GetListItemsByKeyword(ctx, keyword, model.MaxAutoComplete, 0); err != nil {
		err = errors.Wrap(err, "GetListItemsByKeyword")
		return
	}
	values = make([]model.AutoComplete, 0)
	for _, i := range items {
		values = append(values, model.AutoComplete{
			Value: i.Code,
			Text:  fmt.Sprintf("%s - %s", i.Code, i.Name),
		})
	}
	return
}
