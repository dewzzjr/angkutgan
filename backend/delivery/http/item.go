package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetItems get list of items
func (h *HTTP) GetItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// GetItemByCode get item detail
func (h *HTTP) GetItemByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// PatchItemByCode update item detail
func (h *HTTP) PatchItemByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// PostItemByCode create item detail
func (h *HTTP) PostItemByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// DeleteItemByCode delete item detail
func (h *HTTP) DeleteItemByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
