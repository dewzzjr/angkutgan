package admins

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// GetList sequence of users
func (a *Admins) GetList(ctx context.Context, page int, row int) (users []model.UserInfo, err error) {
	// TODO repository
	return
}

// GetByKeyword sequence of users by keyword
func (a *Admins) GetByKeyword(ctx context.Context, page int, row int, key string) (users []model.UserInfo, err error) {
	// TODO repository
	return
}

// Create new user information
func (a *Admins) Create(ctx context.Context, user model.UserInfo, actionBy int64) (err error) {
	if err = a.database.CreateUser(ctx, user, actionBy); err != nil {
		err = errors.Wrap(err, "CreateUser")
		return
	}
	return
}
