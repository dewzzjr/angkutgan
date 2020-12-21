package items

// Items usecase object
type Items struct {
	database iDatabase
}

// New initiate usecase/items
func New(database iDatabase) *Items {
	return &Items{
		database: database,
	}
}

type iDatabase interface {
}
