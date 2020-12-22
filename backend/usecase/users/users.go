package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dgrijalva/jwt-go"
)

// GetByToken convert token to user info
func (u *Users) GetByToken(token string) (user model.Claims, status int, err error) {
	var tkn *jwt.Token
	tkn, err = jwt.ParseWithClaims(token, &user, func(token *jwt.Token) (interface{}, error) {
		return u.Key, nil
	})
	if err == jwt.ErrSignatureInvalid {
		status = http.StatusUnauthorized
		return
	}
	if err != nil {
		status = http.StatusBadRequest
		return
	}
	if !tkn.Valid {
		status = http.StatusUnauthorized
		return
	}
	return
}

// CreateToken from claim
func (u *Users) CreateToken(claim *model.Claims) (token string, err error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err = tkn.SignedString(u.Key)
	return
}

// CreateSession for username
func (u *Users) CreateSession(username string) (claim model.Claims, expire time.Time, err error) {
	// TODO load user information

	expire = time.Now().Add(time.Duration(u.Config.TokenExpiry) * time.Minute)
	claim = model.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}
	return
}

// RefreshSession based on previous claim or session
func (u *Users) RefreshSession(claim *model.Claims) (expire time.Time, ok bool) {
	ok = time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) < time.Duration(u.Config.RefreshToken)*time.Second
	if !ok {
		return
	}
	expire = time.Now().Add(time.Duration(u.Config.TokenExpiry) * time.Minute)
	claim.ExpiresAt = expire.Unix()
	return
}

// Verify username and password
// TODO verify from database
func (u *Users) Verify(username, password string) (ok bool, err error) {
	var users = map[string]string{
		"user1": "password1",
		"user2": "password2",
	}
	expectedPassword, okDB := users[username]
	if !okDB {
		err = fmt.Errorf("not found")
		return
	}
	if expectedPassword != password {
		return
	}
	ok = true
	return
}
