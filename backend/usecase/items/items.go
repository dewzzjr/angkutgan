package items

import (
	"github.com/dewzzjr/angkutgan/backend/model"
)

// GetList sequence of items
func (g *Items) GetList(page int, row int) (items []model.Item, err error) {
	return
}

// GetByKeyword sequence of items by keyword
func (g *Items) GetByKeyword(page int, row int, key string) (items []model.Item, err error) {
	return
}

// Get item by code
func (g *Items) Get(code string) (item model.Item, err error) {
	return
}

// Create new item
func (g *Items) Create(item model.Item) (err error) {
	// TODO check eligible to create
	return
}

// Update item by code
func (g *Items) Update(item model.Item) (err error) {
	return
}

// Remove item by code
func (g *Items) Remove(code string) (err error) {
	// TODO check eligible to delete
	return
}
