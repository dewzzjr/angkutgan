package main

import (
	"flag"
	"os"

	"github.com/dewzzjr/angkutgan/backend/delivery"
	"github.com/dewzzjr/angkutgan/backend/repository"
	"github.com/dewzzjr/angkutgan/backend/usecase"
	"github.com/dewzzjr/angkutgan/backend/view"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/heroku/x/hmetrics/onload"
)

func init() {
	os.Setenv("TZ", "Asia/Jakarta")
}
func main() {
	run := flag.String("run", "backend", "backend/fileserver")
	flag.Parse()
	switch *run {
	case "fileserver":
		fileserver()
	case "backend":
		backend()
	}
}

func backend() {
	r := repository.New()
	u := usecase.New(r)
	v := view.New()
	d := delivery.New(v, u)
	d.Start()
}

func fileserver() {
	v := view.New()
	v.Run()
}
