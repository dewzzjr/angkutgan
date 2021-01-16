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

// GetItems get list of items
func (h *HTTP) GetItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	page, _ := strconv.Atoi(r.FormValue("page"))
	row, _ := strconv.Atoi(r.FormValue("row"))
	keyword := r.FormValue("keyword")
	var err error
	var result []model.Item
	if len(keyword) > 3 {
		result, err = h.items.GetByKeyword(ctx, page, row, keyword)
	} else {
		result, err = h.items.GetList(ctx, page, row)
	}
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": result,
	})
}

// GetItemByCode get item detail
func (h *HTTP) GetItemByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	result, err := h.items.Get(ctx, p.ByName("code"))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": result,
	})
}

// PatchItemByCode update item detail
func (h *HTTP) PatchItemByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.Item{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload.Code = p.ByName("code")
	if err := h.items.Update(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// PostItemByCode create item detail
func (h *HTTP) PostItemByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.Item{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.items.Create(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// DeleteItemByCode delete item detail
func (h *HTTP) DeleteItemByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	code := p.ByName("code")
	if err := h.items.Remove(ctx, code); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
	log.Printf("TRACE %s %s %s", r.Method, r.RequestURI, claims.Username)
}
