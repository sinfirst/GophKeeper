package auth

import (
	"fmt"
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

func CheckToken(tokenFromReq string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenFromReq, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.SecretKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", err
	}
	return claims.Username, nil
}
