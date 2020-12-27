package view

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/dewzzjr/angkutgan/backend/package/response"
	"github.com/julienschmidt/httprouter"
)

// Run html server
func (v *View) Run() {
	v.Load()
	v.Routing()
}

// Routing add routing pattern
func (v *View) Routing() {
	v.Router.GET("/", v.HTML("index"))
	v.Router.GET("/login", v.HTML("login"))
	v.Router.GET("/barang", v.HTML("barang"))
	v.Static.ServeFiles("/assets/*filepath", http.Dir(v.Config.Path+"/assets"))
	v.Router.NotFound = v.Static
}

// Load template
func (v *View) Load() {
	v.files(
		"/index.html",
		"/login.html",
		"/layout/header.html",
		"/layout/script.html",
		"/layout/sidebar.html",
		"/barang/index.html",
		"/barang/daftar.html",
		"/barang/tambah.html",
		"/barang/ubah.html",
		"/barang/jual.html",
		"/barang/sewa.html",
	)
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

func (v View) dir(path ...string) {
	var err error
	for _, p := range path {
		pattern := v.root(p) + "/*.html"
		fmt.Print(pattern)
		v.Template, err = template.ParseGlob(pattern)
		if err != nil {
			log.Fatal(err)
		}
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