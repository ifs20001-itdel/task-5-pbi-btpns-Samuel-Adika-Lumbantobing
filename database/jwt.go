package config

import "github.com/golang-jwt/jwt/v5"

var JWT_KEY = []byte("asdksamdkmskmk12m3kmksdmaks")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
