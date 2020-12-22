package http

import (
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/usecase"
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
		users:  u.Users,
	}
}

type iUsers interface {
	Verify(username, password string) (ok bool, err error)
	CreateSession(username string) (claim model.Claims, expire time.Time, err error)
	CreateToken(claim *model.Claims) (token string, err error)
	GetByToken(token string) (user model.Claims, status int, err error)
	RefreshSession(claim *model.Claims) (expire time.Time, ok bool)
}
