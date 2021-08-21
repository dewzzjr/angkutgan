package view

import (
	"html/template"

	"github.com/dewzzjr/angkutgan/backend/package/config"
	"github.com/julienschmidt/httprouter"
)

// View object
type View struct {
	Static   *httprouter.Router
	Files    *httprouter.Router
	Router   *httprouter.Router
	Config   config.View
	Template *template.Template
}

// New initiate view
func New() *View {
	cfg := config.Get()
	return &View{
		Router: httprouter.New(),
		Static: httprouter.New(),
		Files:  httprouter.New(),
		Config: cfg.View,
	}
}
