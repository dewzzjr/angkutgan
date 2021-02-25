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

// GetRentalByCustomerDate get rental transaction detail by customer code and date
func (h *HTTP) GetRentalByCustomerDate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	date, err := time.Parse(model.ParamFormat, p.ByName("date"))
	if err != nil {
		response.Error(w, err)
		return
	}
	result, err := h.rental.GetDetail(ctx, p.ByName("customer"), date)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": result,
	})
}

// GetRental get list of rental tx
func (h *HTTP) GetRental(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		result, err = h.rental.GetByCustomer(ctx, page, row, customer, date)
	} else {
		result, err = h.rental.GetList(ctx, page, row, date)
	}
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": result,
	})
}

// PostRentalTransaction create rental transaction
func (h *HTTP) PostRentalTransaction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	if err := h.rental.CreateTransaction(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// PatchRentalTransaction edit rental transaction
func (h *HTTP) PatchRentalTransaction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	if err := h.rental.EditTransaction(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}
