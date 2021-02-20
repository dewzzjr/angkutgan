package users

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// Get user information using username
func (u *Users) Get(ctx context.Context, username string) (user model.UserInfo, err error) {
	if user, err = u.database.GetUserProfile(ctx, username); err != nil {
		err = errors.Wrap(err, "GetUserProfile")
	}
	return
}

// Edit user information using user id
func (u *Users) Edit(ctx context.Context, user model.UserInfo, actionBy int64) (err error) {
	if err = u.database.EditUser(ctx, user, actionBy); err != nil {
		err = errors.Wrap(err, "EditUser")
	}
	return
}

// ChangePassword by username
func (u *Users) ChangePassword(ctx context.Context, username, newPass, oldPass string) (err error) {
	var ok bool
	if ok, err = u.database.VerifyUser(ctx, username, oldPass); err != nil {
		err = errors.Wrap(err, "VerifyUser")
		return
	}
	if !ok {
		err = errors.New("password salah")
		return
	}
	if err = u.database.ChangePassword(ctx, username, newPass); err != nil {
		err = errors.Wrap(err, "ChangePassword")
		return
	}
	return
}
