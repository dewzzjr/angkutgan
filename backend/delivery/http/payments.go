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

// PatchPaymentByTxID update last payment in the transaction
func (h *HTTP) PatchPaymentByTxID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.Payment{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	txID, err := strconv.ParseInt(p.ByName("txid"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.payments.Edit(ctx, txID, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// PostPaymentByTxID create new payment in the transaction
func (h *HTTP) PostPaymentByTxID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.Payment{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	txID, err := strconv.ParseInt(p.ByName("txid"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.payments.Add(ctx, txID, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// DeletePaymentByTxID delete last payment in the transaction
func (h *HTTP) DeletePaymentByTxID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	txID, err := strconv.ParseInt(p.ByName("txid"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.payments.Delete(ctx, txID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
	log.Printf("TRACE %s %s %s", r.Method, r.RequestURI, claims.Username)
}
