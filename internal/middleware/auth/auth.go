package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sinfirst/GophKeeper/internal/config"
)

type Claims struct {
	jwt.RegisteredClaims
	Username string
}

func BuildJWTString(user string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExp)),
		},
		Username: user,
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
