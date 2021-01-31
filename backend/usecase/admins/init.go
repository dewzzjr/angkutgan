package admins

import (
	"context"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// Admins usecase object
type Admins struct {
	database iDatabase
}

// New initiate usecase/admins
func New(database iDatabase) *Admins {
	return &Admins{
		database: database,
	}
}

type iDatabase interface {
	CreateUser(ctx context.Context, data model.UserInfo, actionBy int64) (err error)
}
