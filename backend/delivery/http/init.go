package http

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/usecase"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

// HTTP delivery object
type HTTP struct {
	Router *httprouter.Router
	Static *httprouter.Router
	Config model.Delivery
	users  iUsers
}

// New initiate delivery/http
func New(cfg model.Delivery, u *usecase.Usecase) *HTTP {
	return &HTTP{
		Router: httprouter.New(),
		Static: httprouter.New(),
		Config: cfg,
		// users: u.Users,
		users: bypass(u.Users, cfg.ByPass),
	}
}

type iUsers interface {
	Verify(ctx context.Context, username, password string) (ok bool, err error)
	CreateSession(ctx context.Context, username string) (claim model.Claims, expire time.Time, err error)
	CreateToken(ctx context.Context, claim *model.Claims) (token string, err error)
	GetByToken(ctx context.Context, token string) (user model.Claims, tkn *jwt.Token, err error)
	RefreshSession(ctx context.Context, claim *model.Claims) (expire time.Time, ok bool)
	Create(ctx context.Context, data model.UserInfo, actionBy int64) (err error)
}
