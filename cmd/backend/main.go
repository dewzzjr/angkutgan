package main

import (
	"os"

	"github.com/dewzzjr/angkutgan/backend/delivery"
	"github.com/dewzzjr/angkutgan/backend/repository"
	"github.com/dewzzjr/angkutgan/backend/usecase"
	"github.com/dewzzjr/angkutgan/backend/view"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	os.Setenv("TZ", "Asia/Jakarta")
}

func main() {
	r := repository.New()
	u := usecase.New(r)
	v := view.New()
	d := delivery.New(v, u)
	d.Start()
}
