package delivery

import (
	"flag"

	"github.com/dewzzjr/angkutgan/backend/delivery/http"
	"github.com/dewzzjr/angkutgan/backend/package/config"
	"github.com/dewzzjr/angkutgan/backend/usecase"
)

// Delivery object
type Delivery struct {
	http    *http.HTTP
	usecase *usecase.Usecase
}

// New initiate delivery
func New(u *usecase.Usecase) *Delivery {
	cfg := config.Get()
	return &Delivery{
		http:    http.New(cfg.Delivery),
		usecase: u,
	}
}

var service string

// Start delivery using service type
func (d *Delivery) Start() {
	flag.StringVar(&service, "service", "http", "type of service [http]")
	flag.Parse()

	switch service {
	case "http":
		d.http.Run()
	default:
	}
}
