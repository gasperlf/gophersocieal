package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "my-secret-key"

var testClaims = jwt.MapClaims{
	"aud": "test-audience",
	"iss": "test-audience",
	"sub": int64(1),
	"exp": time.Now().Add(time.Hour * 24).Unix(),
}

type TestAuthenticator struct{}

func (a *TestAuthenticator) GenerateToken(tclaims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)

	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString, nil
}

func (a *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {

	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
}
