package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/package/response"
	"github.com/julienschmidt/httprouter"
)

// PostShipmentByTxID create new shipment in the transaction
func (h *HTTP) PostShipmentByTxID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.Shipment{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(payload.Items) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := time.Parse(model.DateFormat, payload.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	txID, err := strconv.ParseInt(p.ByName("txid"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.shipment.Add(ctx, txID, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// PatchShipmentByTxID update shipment by date in the transaction
func (h *HTTP) PatchShipmentByTxID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.Shipment{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(payload.Items) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := time.Parse(model.DateFormat, payload.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	txID, err := strconv.ParseInt(p.ByName("txid"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = time.Parse(model.ParamFormat, p.ByName("date"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.shipment.Edit(ctx, txID, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
	log.Printf("TRACE %s %s %s", r.Method, r.RequestURI, claims.Username)
}

// DeleteShipmentByTxID delete shipment by date in the transaction
func (h *HTTP) DeleteShipmentByTxID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.Shipment{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	date, err := time.Parse(model.DateFormat, payload.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	txID, err := strconv.ParseInt(p.ByName("txid"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.shipment.Delete(ctx, txID, date, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
	log.Printf("TRACE %s %s %s", r.Method, r.RequestURI, claims.Username)
}
