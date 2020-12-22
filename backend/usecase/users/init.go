package users

import (
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
}
