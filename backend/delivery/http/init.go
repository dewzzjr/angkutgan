package http

import (
	"context"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/dewzzjr/angkutgan/backend/package/config"
	"github.com/dewzzjr/angkutgan/backend/usecase"
	"github.com/dewzzjr/angkutgan/backend/view"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

// HTTP delivery object
type HTTP struct {
	Router    *httprouter.Router
	View      *view.View
	Config    config.Delivery
	users     iUsers
	items     iItems
	customers iCustomers
	sales     iSales
	rental    iRental
	payments  iPayments
	shipment  iShipment
	ajax      iAjax
}

// New initiate delivery/http
func New(cfg config.Delivery, v *view.View, u *usecase.Usecase) *HTTP {
	return &HTTP{
		Router:    httprouter.New(),
		Config:    cfg,
		View:      v,
		users:     u.Users,
		items:     u.Items,
		customers: u.Customers,
		sales:     u.Sales,
		rental:    u.Rental,
		payments:  u.Payments,
		shipment:  u.Shipment,
		ajax:      u.Ajax,
	}
}

type iUsers interface {
	Verify(ctx context.Context, username, password string) (ok bool, err error)
	CreateSession(ctx context.Context, username string) (claim model.Claims, expire time.Time, err error)
	CreateToken(ctx context.Context, claim *model.Claims) (token string, err error)
	GetByToken(ctx context.Context, token string) (user model.Claims, tkn *jwt.Token, err error)
	RefreshSession(ctx context.Context, claim *model.Claims) (expire time.Time, ok bool)
	Create(ctx context.Context, data model.UserInfo, actionBy int64) (err error)
	Get(ctx context.Context, username string) (user model.UserInfo, err error)
	Edit(ctx context.Context, user model.UserInfo, actionBy int64) (err error)
	ChangePassword(ctx context.Context, username, newPass, oldPass string) (err error)
}
type iItems interface {
	GetList(ctx context.Context, page int, row int) (items []model.Item, err error)
	GetByKeyword(ctx context.Context, page int, row int, key string) (items []model.Item, err error)
	Get(ctx context.Context, code string) (item model.Item, err error)
	Create(ctx context.Context, item model.Item, actionBy int64) (err error)
	Update(ctx context.Context, item model.Item, actionBy int64) (err error)
	Remove(ctx context.Context, code string) (err error)
}
type iCustomers interface {
	GetList(ctx context.Context, page int, row int) (customers []model.Customer, err error)
	GetByKeyword(ctx context.Context, page int, row int, key string) (customers []model.Customer, err error)
	Get(ctx context.Context, code string) (customer model.Customer, err error)
	Create(ctx context.Context, customer model.Customer, actionBy int64) (err error)
	Update(ctx context.Context, customer model.Customer, actionBy int64) (err error)
	Remove(ctx context.Context, code string) (err error)
}
type iAjax interface {
	IsValidUsername(ctx context.Context, newUsername, oldUsername string) (ok bool, err error)
	IsValidItemCode(ctx context.Context, code string) (ok bool, err error)
	GetItems(ctx context.Context, keyword string) (values []model.AutoComplete, err error)
	IsValidCustomerCode(ctx context.Context, code string) (ok bool, err error)
	GetCustomers(ctx context.Context, keyword string) (values []model.AutoComplete, err error)
}
type iSales interface {
	GetDetail(ctx context.Context, code string, date time.Time) (tx model.Transaction, err error)
	GetByCustomer(ctx context.Context, page, row int, customer string, date time.Time) (txs []model.Transaction, err error)
	GetList(ctx context.Context, page, row int, date time.Time) (txs []model.Transaction, err error)
	CreateTransaction(ctx context.Context, tx model.CreateTransaction, actionBy int64) (err error)
	EditTransaction(ctx context.Context, tx model.CreateTransaction, actionBy int64) (err error)
	Cancel(ctx context.Context, code string, date time.Time) (err error)
}
type iRental interface {
	iSales
}
type iPayments interface {
	Add(ctx context.Context, txID int64, pay model.Payment, actionBy int64) (err error)
	Edit(ctx context.Context, txID int64, pay model.Payment, actionBy int64) (err error)
	Delete(ctx context.Context, txID int64) (err error)
}
type iShipment interface {
	Add(ctx context.Context, txID int64, pay model.Shipment, actionBy int64) (err error)
	Edit(ctx context.Context, txID int64, pay model.Shipment, actionBy int64) (err error)
	Delete(ctx context.Context, txID int64, date time.Time, actionBy int64) (err error)
}
