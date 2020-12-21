package usecase

import (
	"github.com/dewzzjr/angkutgan/backend/repository"
	"github.com/dewzzjr/angkutgan/backend/usecase/items"
)

// Usecase object
type Usecase struct {
	Items *items.Items
}

// New initiate usecase
func New(r *repository.Repository) *Usecase {
	return &Usecase{
		Items: items.New(r.Database),
	}
}
