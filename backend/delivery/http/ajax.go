package http

import (
	"log"
	"net/http"

	"github.com/dewzzjr/angkutgan/backend/package/response"
	"github.com/julienschmidt/httprouter"
)

// AJAX used to serve ajax call usecase
func (h *HTTP) AJAX(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	action := r.FormValue("action")
	switch r.Method {
	case http.MethodGet:
		if action == "username" {
			var message string
			ok, err := h.ajax.IsValidUsername(ctx, r.FormValue("old"), r.FormValue("new"))
			if err != nil {
				message = "gagal validasi"
				log.Print(err)
			}
			if !ok {
				message = "username sudah dipakai"
			}
			if message != "" {
				response.JSON(w, map[string]interface{}{
					"valid":   false,
					"message": message,
				})
				return
			}
			response.JSON(w, map[string]interface{}{
				"valid": true,
			})
			return
		}
		if action == "validate_code_item" {
			var message string
			ok, err := h.ajax.IsValidItemCode(ctx, r.FormValue("new"))
			if err != nil {
				message = "gagal validasi"
				log.Print(err)
			}
			if !ok {
				message = "kode sudah dipakai"
			}
			if message != "" {
				response.JSON(w, map[string]interface{}{
					"valid":   false,
					"message": message,
				})
				return
			}
			response.JSON(w, map[string]interface{}{
				"valid": true,
			})
			return
		}
		if action == "validate_code_customer" {
			var message string
			ok, err := h.ajax.IsValidCustomerCode(ctx, r.FormValue("new"))
			if err != nil {
				message = "gagal validasi"
				log.Print(err)
			}
			if !ok {
				message = "kode sudah dipakai"
			}
			if message != "" {
				response.JSON(w, map[string]interface{}{
					"valid":   false,
					"message": message,
				})
				return
			}
			response.JSON(w, map[string]interface{}{
				"valid": true,
			})
			return
		}
		if action == "items" {
			result, err := h.ajax.GetItems(ctx, r.FormValue("q"))
			if err != nil {
				log.Print(err)
			}
			response.JSON(w, result)
			return
		}
		if action == "customers" {
			result, err := h.ajax.GetCustomers(ctx, r.FormValue("q"))
			if err != nil {
				log.Print(err)
			}
			response.JSON(w, result)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
