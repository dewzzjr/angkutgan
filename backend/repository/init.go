package repository

import (
	"github.com/dewzzjr/angkutgan/backend/package/config"
	"github.com/dewzzjr/angkutgan/backend/repository/database"
)

// Repository object
type Repository struct {
	Database *database.Database
}

// New initiate repository
func New() *Repository {
	cfg := config.Get()
	return &Repository{
		Database: database.New(cfg.Repository),
	}
}
