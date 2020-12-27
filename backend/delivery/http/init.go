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
	View   *httprouter.Router
	Config model.Delivery
	users  iUsers
	items  iItems
	ajax   iAjax
}

// New initiate delivery/http
func New(cfg model.Delivery, v *httprouter.Router, u *usecase.Usecase) *HTTP {
	return &HTTP{
		Router: httprouter.New(),
		View:   v,
		Config: cfg,
		// users: u.Users,
		users: bypass(u.Users, cfg.ByPass),
		items: u.Items,
		ajax:  u.Ajax,
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

type iItems interface {
	GetList(ctx context.Context, page int, row int) (items []model.Item, err error)
	GetByKeyword(ctx context.Context, page int, row int, key string) (items []model.Item, err error)
	Get(ctx context.Context, code string) (item model.Item, err error)
	Create(ctx context.Context, item model.Item, actionBy int64) (err error)
	Update(ctx context.Context, item model.Item, actionBy int64) (err error)
	Remove(ctx context.Context, code string) (err error)
}

type iAjax interface {
	IsValidUsername(ctx context.Context, newUsername, oldUsername string) (ok bool, err error)
}
