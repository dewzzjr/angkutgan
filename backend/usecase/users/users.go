package users

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Get user information using username
func (u *Users) Get(ctx context.Context, username string) (user model.UserInfo, err error) {
	// TODO repository
	return
}

// Edit user information using user id
func (u *Users) Edit(ctx context.Context, user model.UserInfo, actionBy int) (err error) {
	// TODO repository
	return
}
