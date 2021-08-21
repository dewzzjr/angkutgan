package view

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/dewzzjr/angkutgan/backend/package/middleware"
	"github.com/dewzzjr/angkutgan/backend/package/response"
	"github.com/julienschmidt/httprouter"
)

// Run html server
func (v *View) Run() {
	v.Load()
	v.Routing()
	port := fmt.Sprintf(":%d", v.Config.ViewPort)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, middleware.NewLogger(v.Router)))
}

// Routing add routing pattern
func (v *View) Routing() {
	for name, url := range URLs {
		v.Router.GET(url, v.HTML(name))
	}
	v.Static.ServeFiles("/assets/*filepath", http.Dir(v.Config.Path+"/assets"))
	v.Files.ServeFiles("/*filepath", http.Dir(v.Config.Path+"/html"))
	v.Router.NotFound = v.Static
	v.Static.NotFound = v.Files
}

// Load template
func (v *View) Load() {
	v.files(Templates...)
}

func (v *View) files(path ...string) {
	var err error
	for i, p := range path {
		path[i] = v.root(p)
	}
	v.Template, err = template.ParseFiles(path...)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *View) root(path string) string {
	return v.Config.Path + path
}

// HTML create router from template
func (v *View) HTML(name string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := v.Template.ExecuteTemplate(w, name, nil); err != nil {
			response.Error(w, err)
			return
		}
	}
}
