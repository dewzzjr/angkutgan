package model

import (
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Credentials user sign in payload
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Claims struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserID   int64    `json:"uid"`
	Username string   `json:"username"`
	Fullname string   `json:"fullname"`
	UAM      []string `json:"uam"`
	jwt.StandardClaims
}

// UserInfo user information
// to create and get biodata
type UserInfo struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username"`
	Type     string `json:"type"`
	Biodata
}

// Biodata user detail information
type Biodata struct {
	Fullname     string        `json:"fullname" db:"fullname"`
	Address      string        `json:"address" db:"address"`
	Phone        string        `json:"phone" db:"phone"`
	NIK          string        `json:"nik" db:"nik"`
	KTP          string        `json:"ktp" db:"ktp"`
	Religion     string        `json:"religion" db:"religion"`
	BirthdateStr string        `json:"birthdate" db:"-"`
	Birthdate    time.Time     `json:"-" db:"birthdate"`
	ModifiedBy   sql.NullInt64 `json:"-" db:"modified_by"`
	UpdateTime   sql.NullTime  `json:"-" db:"update_time"`
	InfoID       int64         `json:"info_id,omitempty" db:"id"`
}
