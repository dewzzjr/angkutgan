package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dewzzjr/angkutgan/backend/package/middleware"
)

// Run http server
func (h *HTTP) Run() {
	h.Routing()
	port := fmt.Sprintf(":%d", h.Config.Port)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, middleware.NewLogger(h.Router)))
}

// Routing add routing pattern
func (h *HTTP) Routing() {
	h.Router.POST("/user/login", h.Login)
	h.Router.POST("/user/logout", h.Login)
	h.Router.POST("/user/session", h.Refresh)
	h.Router.GET("/user/info", h.GetUserInfo)

	h.Router.POST("/user/create", h.CreateUser)

	h.Router.GET("/items", h.GetItems)
	h.Router.POST("/item", h.PostItemByCode)
	h.Router.GET("/item/:code", h.GetItemByCode)
	h.Router.PATCH("/item/:code", h.PatchItemByCode)
	h.Router.DELETE("/item/:code", h.DeleteItemByCode)
	h.Static.ServeFiles("/*filepath", http.Dir(h.Config.StaticPath))
	h.Router.NotFound = h.Static
}
