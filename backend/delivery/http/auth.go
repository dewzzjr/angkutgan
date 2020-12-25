package http

import (
	"context"
	"log"
	"net/http"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dgrijalva/jwt-go"
)

// Auth function to authenticate
func (h *HTTP) Auth(ctx context.Context, w http.ResponseWriter, r *http.Request) (claim model.Claims, ok bool) {
	var c *http.Cookie
	var err error
	c, err = r.Cookie(h.Config.CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var token *jwt.Token
	claim, token, err = h.users.GetByToken(ctx, c.Value)
	if err == jwt.ErrSignatureInvalid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok = true
	return
}
