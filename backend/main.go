package main

import (
	"github.com/dewzzjr/angkutgan/backend/delivery"
	"github.com/dewzzjr/angkutgan/backend/repository"
	"github.com/dewzzjr/angkutgan/backend/usecase"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := repository.New()
	u := usecase.New(r)
	d := delivery.New(u)
	d.Start()
}
