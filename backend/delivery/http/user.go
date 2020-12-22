package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/package/response"
	"github.com/julienschmidt/httprouter"
)

// Login sign in using jwt
func (h *HTTP) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var creds model.Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := h.users.Verify(creds.Username, creds.Password)
	if err != nil {
		response.Error(w, err)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, expirationTime, err := h.users.CreateSession(creds.Username)
	if err != nil {
		response.Error(w, err)
		return
	}

	tokenString, err := h.users.CreateToken(&claims)
	if err != nil {
		response.Error(w, err)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    h.Config.CookieName,
		Value:   tokenString,
		Expires: expirationTime,
	})
	response.JSON(w, map[string]interface{}{
		"result": claims,
	})
	return
}

// Refresh jwt token
func (h *HTTP) Refresh(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, err := r.Cookie(h.Config.CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims, status, err := h.users.GetByToken(c.Value)
	if status != 0 {
		w.WriteHeader(status)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	expirationTime, ok := h.users.RefreshSession(&claims)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenString, err := h.users.CreateToken(&claims)
	if err != nil {
		response.Error(w, err)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    h.Config.CookieName,
		Value:   tokenString,
		Expires: expirationTime,
	})
	response.JSON(w, map[string]interface{}{
		"result": claims,
	})
	return
}

// GetUserInfo get user information
func (h *HTTP) GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie(h.Config.CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims, status, err := h.users.GetByToken(c.Value)
	if status != 0 {
		w.WriteHeader(status)
		return
	}
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, map[string]interface{}{
		"result": claims,
	})
	return
}
