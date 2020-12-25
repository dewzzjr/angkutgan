package users

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Users usecase object
type Users struct {
	database iDatabase
	Config   model.Users
	Key      []byte
}

// New initiate usecase/users
func New(database iDatabase, cfg model.Users) *Users {
	return &Users{
		database: database,
		Config:   cfg,
		Key:      []byte(cfg.JWTKey),
	}
}

type iDatabase interface {
	GetUserLogin(ctx context.Context, claim *model.Claims) (err error)
	VerifyUser(ctx context.Context, username, password string) (ok bool, err error)
	CreateUser(ctx context.Context, data model.UserInfo, actionBy int64) (err error)
}