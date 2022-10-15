package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

type Authenticator interface {
	Issue(phone string, isVerified bool) (string, error)
	GetMiddlewareConfig() middleware.JWTConfig
}

type Claim struct {
	PhoneNumber string `json:"phone_number"`
	IsVerified  bool   `json:"is_verified"`
	jwt.StandardClaims
}

type auth struct {
	secret             string
	expirationDuration time.Duration
}

func NewAuth(secret string, exp time.Duration) Authenticator {
	return &auth{
		secret:             secret,
		expirationDuration: exp,
	}
}

func (a *auth) Issue(phone string, isVerified bool) (string, error) {
	expirationTime := time.Now().Add(a.expirationDuration)
	claims := &Claim{
		PhoneNumber: phone,
		IsVerified:  isVerified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("secret:", a.secret)
	return token.SignedString([]byte(a.secret))
}

func (a *auth) GetMiddlewareConfig() middleware.JWTConfig {
	fmt.Println("----secret:", a.secret)
	return middleware.JWTConfig{
		Claims:     &Claim{},
		SigningKey: []byte(a.secret),
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/api/v1/accounts" {
				return true
			}
			if c.Request().URL.Path == "/api/v1/accounts/auth" {
				return true
			}
			return false
		},
	}
}
