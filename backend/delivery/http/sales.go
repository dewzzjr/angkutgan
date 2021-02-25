package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/package/response"
	"github.com/julienschmidt/httprouter"
)

// GetSalesByCustomerDate get sales transaction detail by customer code and date
func (h *HTTP) GetSalesByCustomerDate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	date, err := time.Parse(model.ParamFormat, p.ByName("date"))
	if err != nil {
		response.Error(w, err)
		return
	}
	result, err := h.sales.GetDetail(ctx, p.ByName("customer"), date)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": result,
	})
}

// PostSalesTransaction create sales transaction
func (h *HTTP) PostSalesTransaction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.CreateTransaction{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.sales.CreateTransaction(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// PatchSalesTransaction edit sales transaction
func (h *HTTP) PatchSalesTransaction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.CreateTransaction{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.sales.EditTransaction(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// GetSales get list of sales tx
func (h *HTTP) GetSales(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	page, _ := strconv.Atoi(r.FormValue("page"))
	row, _ := strconv.Atoi(r.FormValue("row"))
	customer := r.FormValue("customer")
	date, err := time.Parse(model.DateFormat, r.FormValue("date"))
	if err != nil {
		date = time.Now()
	}
	var result []model.Transaction
	if customer != "" {
		result, err = h.sales.GetByCustomer(ctx, page, row, customer, date)
	} else {
		result, err = h.sales.GetList(ctx, page, row, date)
	}
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": result,
	})
}
