package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/package/response"
	"github.com/julienschmidt/httprouter"
)

// GetCustomers get list of customers
func (h *HTTP) GetCustomers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	page, _ := strconv.Atoi(r.FormValue("page"))
	row, _ := strconv.Atoi(r.FormValue("row"))
	keyword := r.FormValue("keyword")
	var err error
	var result []model.Customer
	if len(keyword) > model.MinLengthKeyword {
		result, err = h.customers.GetByKeyword(ctx, page, row, keyword)
	} else {
		result, err = h.customers.GetList(ctx, page, row)
	}
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": result,
	})
}

// GetCustomerByCode get customer detail
func (h *HTTP) GetCustomerByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	result, err := h.customers.Get(ctx, p.ByName("code"))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": result,
	})
}

// PatchCustomerByCode update customer detail
func (h *HTTP) PatchCustomerByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.Customer{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload.Code = p.ByName("code")
	if err := h.customers.Update(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// PostCustomerByCode create customer detail
func (h *HTTP) PostCustomerByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.Customer{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.customers.Create(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// DeleteCustomerByCode delete customer detail
func (h *HTTP) DeleteCustomerByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	code := p.ByName("code")
	if err := h.customers.Remove(ctx, code); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
	log.Printf("TRACE %s %s %s", r.Method, r.RequestURI, claims.Username)
}
