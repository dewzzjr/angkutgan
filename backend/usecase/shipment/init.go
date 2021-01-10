package shipment

// Shipment usecase object
type Shipment struct {
	database iDatabase
}

// New initiate usecase/shipment
func New(database iDatabase) *Shipment {
	return &Shipment{
		database: database,
	}
}

type iDatabase interface {
}
