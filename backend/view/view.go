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
	port := fmt.Sprintf(":%d", v.Config.Port)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, middleware.NewLogger(v.Router)))
}

// Routing add routing pattern
func (v *View) Routing() {
	v.Router.GET("/", v.HTML("index"))
	v.Router.GET("/login", v.HTML("login"))
	v.Router.GET("/barang", v.HTML("barang"))
	v.Router.GET("/pelanggan", v.HTML("pelanggan"))
	v.Router.GET("/penjualan", v.HTML("penjualan"))
	v.Router.GET("/persewaan", v.HTML("persewaan"))
	v.Static.ServeFiles("/assets/*filepath", http.Dir(v.Config.Path+"/assets"))
	v.Files.ServeFiles("/*filepath", http.Dir(v.Config.Path+"/html"))
	v.Router.NotFound = v.Static
	v.Static.NotFound = v.Files
}

// Load template
func (v *View) Load() {
	v.files(
		"/index.html",
		"/login.html",
		"/layout/header.html",
		"/layout/script.html",
		"/layout/sidebar.html",
		"/layout/modal.html",
		"/barang/index.html",
		"/barang/daftar.html",
		"/barang/tambah.html",
		"/barang/ubah.html",
		"/barang/harga.html",
		"/pelanggan/index.html",
		"/pelanggan/daftar.html",
		"/pelanggan/tambah.html",
		"/pelanggan/ubah.html",
		"/persewaan/index.html",
		"/persewaan/daftar.html",
		"/persewaan/buat.html",
		"/penjualan/index.html",
		"/penjualan/daftar.html",
		"/penjualan/buat.html",
		"/pengiriman/keluar.html",
		"/pengiriman/masuk.html",
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
