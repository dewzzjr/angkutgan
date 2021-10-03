package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/package/response"
	"github.com/julienschmidt/httprouter"
)

// Login sign in using jwt
func (h *HTTP) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	creds := model.Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ok, err := h.users.Verify(ctx, creds.Username, creds.Password)
	if err != nil {
		response.Error(w, err)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims, expirationTime, err := h.users.CreateSession(ctx, creds.Username)
	if err != nil {
		response.Error(w, err)
		return
	}
	tokenString, err := h.users.CreateToken(ctx, &claims)
	if err != nil {
		response.Error(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    h.Config.CookieName,
		Value:   tokenString,
		Expires: expirationTime,
		Domain:  r.Host,
		Path:    "/",
	})
	response.JSON(w, map[string]interface{}{
		"result": claims,
	})
}

// Logout remove cookie
func (h *HTTP) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.SetCookie(w, &http.Cookie{
		Name:    h.Config.CookieName,
		Domain:  r.Host,
		Expires: time.Now(),
	})
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// Refresh jwt token
func (h *HTTP) Refresh(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	expirationTime, ok := h.users.RefreshSession(ctx, &claims)
	if !ok {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	tokenString, err := h.users.CreateToken(ctx, &claims)
	if err != nil {
		response.Error(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    h.Config.CookieName,
		Value:   tokenString,
		Expires: expirationTime,
		Domain:  r.Host,
		Path:    "/",
	})
	response.JSON(w, map[string]interface{}{
		"result": claims,
	})
}

// GetUserInfo get user information
func (h *HTTP) GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": claims,
	})
}

// CreateUser add new user information
func (h *HTTP) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.UserInfo{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload.Birthdate, _ = time.Parse(model.DateFormat, payload.BirthdateStr)
	if err := h.users.Create(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// GetUserProfile by token session cookies
func (h *HTTP) GetUserProfile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	result, err := h.users.Get(ctx, claims.Username)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": result,
	})
}

// PatchUserProfile by token session cookies
func (h *HTTP) PatchUserProfile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	payload := model.UserInfo{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload.ID = claims.UserID
	payload.Birthdate, _ = time.Parse(model.DateFormat, payload.BirthdateStr)
	if err := h.users.Edit(ctx, payload, claims.UserID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}

// ChangePassword by old password new password and username from token session cookies
func (h *HTTP) ChangePassword(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	claims, ok := h.Auth(ctx, w, r)
	if !ok {
		return
	}
	var payload struct {
		New string `json:"new"`
		Old string `json:"old"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.users.ChangePassword(ctx, claims.Username, payload.New, payload.Old); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, map[string]interface{}{
		"result": "OK",
	})
}
