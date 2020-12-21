package http

import (
	"fmt"
	"log"
	"net/http"
)

// Run http server
func (h *HTTP) Run() {
	h.Routing()
	port := fmt.Sprintf(":%d", h.Config.Port)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, h.Router))
}

// Routing add routing pattern
func (h *HTTP) Routing() {
	h.Router.GET("/items", h.GetItems)
	h.Router.GET("/item/:code", h.GetItemByCode)
	h.Router.POST("/item/:code", h.PostItemByCode)
	h.Router.PATCH("/item/:code", h.PatchItemByCode)
	h.Router.DELETE("/item/:code", h.PatchItemByCode)
}
