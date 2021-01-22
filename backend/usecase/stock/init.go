package stock

// Stock usecase object
type Stock struct {
	database iDatabase
}

// New initiate usecase/shipment
func New(database iDatabase) *Stock {
	return &Stock{
		database: database,
	}
}

type iDatabase interface {
}
