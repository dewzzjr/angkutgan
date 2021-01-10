package payments

// Payments usecase object
type Payments struct {
	database iDatabase
}

// New initiate usecase/payments
func New(database iDatabase) *Payments {
	return &Payments{
		database: database,
	}
}

type iDatabase interface {
}
