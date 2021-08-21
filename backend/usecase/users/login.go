package users

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// GetByToken convert token to user info
func (u *Users) GetByToken(ctx context.Context, token string) (user model.Claims, tkn *jwt.Token, err error) {
	tkn, err = jwt.ParseWithClaims(token, &user, func(token *jwt.Token) (interface{}, error) {
		return u.Key, nil
	})
	if err != nil {
		err = errors.Wrap(err, "ParseWithClaims")
		return
	}
	return
}

// CreateToken from claim
func (u *Users) CreateToken(ctx context.Context, claim *model.Claims) (token string, err error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err = tkn.SignedString(u.Key)
	err = errors.Wrap(err, "SignedString")
	return
}

// CreateSession for username
func (u *Users) CreateSession(ctx context.Context, username string) (claim model.Claims, expire time.Time, err error) {
	expire = time.Now().Add(time.Duration(u.Config.TokenExpiry) * time.Minute)
	claim = model.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}
	if err = u.database.GetUserLogin(ctx, &claim); err != nil {
		err = errors.Wrap(err, "GetUserLogin")
		return
	}
	return
}

// RefreshSession based on previous claim or session
func (u *Users) RefreshSession(ctx context.Context, claim *model.Claims) (expire time.Time, ok bool) {
	if ok = time.Until(time.Unix(claim.ExpiresAt, 0)) < (time.Duration(u.Config.RefreshToken) * time.Second); !ok {
		return
	}
	expire = time.Now().Add(time.Duration(u.Config.TokenExpiry) * time.Minute)
	claim.ExpiresAt = expire.Unix()
	return
}

// Verify username and password
func (u *Users) Verify(ctx context.Context, username, password string) (ok bool, err error) {
	if ok, err = u.database.VerifyUser(ctx, username, password); err != nil {
		err = errors.Wrap(err, "VerifyUser")
		return
	}
	return
}

// Create new user
func (u *Users) Create(ctx context.Context, data model.UserInfo, actionBy int64) (err error) {
	if err = u.database.CreateUser(ctx, data, actionBy); err != nil {
		err = errors.Wrap(err, "CreateUser")
		return
	}
	return
}
