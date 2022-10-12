package jwt

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Authenticator interface {
	Issue(phone string, isVerified bool) (string, error)
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
	return token.SignedString([]byte(a.secret))
}
