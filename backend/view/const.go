package view

// Templates files
var Templates = []string{
	"/index.html",
	"/login.html",
	"/layout/header.html",
	"/layout/script.html",
	"/layout/sidebar.html",
	"/layout/pagination.html",
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
}

// URLs by template name
var URLs = map[string]string{
	"index":     "/",
	"login":     "/login",
	"barang":    "/barang",
	"pelanggan": "/pelanggan",
	"penjualan": "/penjualan",
	"persewaan": "/persewaan",
}
