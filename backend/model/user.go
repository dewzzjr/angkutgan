package model

import "github.com/dgrijalva/jwt-go"

// Credentials user sign in payload
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Claims struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
