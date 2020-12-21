package http

import (
	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/julienschmidt/httprouter"
)

// HTTP delivery object
type HTTP struct {
	Router *httprouter.Router
	Config model.Delivery
}

// New initiate delivery/http
func New(cfg model.Delivery) *HTTP {
	return &HTTP{
		Router: httprouter.New(),
		Config: cfg,
	}
}
