package http

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dgrijalva/jwt-go"
)

type iuser struct {
	i      iUsers
	bypass bool
}

func bypass(i iUsers, ok bool) *iuser {
	return &iuser{i, ok}
}
func (i *iuser) Verify(ctx context.Context, username, password string) (ok bool, err error) {
	if i.bypass {
		return true, nil
	}
	return i.i.Verify(ctx, username, password)
}
func (i *iuser) CreateSession(ctx context.Context, username string) (claim model.Claims, expire time.Time, err error) {
	if i.bypass {
		expire = time.Now().Add(5 * time.Minute)
		claim = model.Claims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expire.Unix(),
			},
		}
		return
	}
	return i.i.CreateSession(ctx, username)
}
func (i *iuser) CreateToken(ctx context.Context, claim *model.Claims) (token string, err error) {
	return i.i.CreateToken(ctx, claim)
}
func (i *iuser) GetByToken(ctx context.Context, token string) (user model.Claims, tkn *jwt.Token, err error) {
	return i.i.GetByToken(ctx, token)
}
func (i *iuser) RefreshSession(ctx context.Context, claim *model.Claims) (expire time.Time, ok bool) {
	return i.i.RefreshSession(ctx, claim)
}
func (i *iuser) Create(ctx context.Context, data model.UserInfo, actionBy int64) (err error) {
	return i.i.Create(ctx, data, actionBy)
}
