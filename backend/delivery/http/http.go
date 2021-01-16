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
	h.Router.POST("/user/logout", h.Logout)
	h.Router.POST("/user/session", h.Refresh)
	h.Router.GET("/user/info", h.GetUserInfo)

	h.Router.POST("/user/create", h.CreateUser)
	h.Router.GET("/ajax", h.AJAX)

	h.Router.GET("/items", h.GetItems)
	h.Router.POST("/item", h.PostItemByCode)
	h.Router.GET("/item/:code", h.GetItemByCode)
	h.Router.PATCH("/item/:code", h.PatchItemByCode)
	h.Router.DELETE("/item/:code", h.DeleteItemByCode)

	h.Router.GET("/customers", h.GetCustomers)
	h.Router.POST("/customer", h.PostCustomerByCode)
	h.Router.GET("/customer/:code", h.GetCustomerByCode)
	h.Router.PATCH("/customer/:code", h.PatchCustomerByCode)
	h.Router.DELETE("/customer/:code", h.DeleteCustomerByCode)

	h.Router.GET("/sales/:customer/:date", h.GetSalesByCustomerDate)
	h.Router.POST("/sales", h.PostSalesTransaction)
	h.Router.PATCH("/sales", h.PatchSalesTransaction)

	h.Router.POST("/payment/:txid", h.PostPaymentByTxID)
	h.Router.PATCH("/payment/:txid", h.PatchPaymentByTxID)
	h.Router.DELETE("/payment/:txid", h.DeletePaymentByTxID)
	h.Router.NotFound = h.View.Router
}
