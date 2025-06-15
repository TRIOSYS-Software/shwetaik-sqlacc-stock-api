package utils

import (
	"errors"
	"shwetaik-sqlacc-stock-api/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = config.Cfg.JWT_SECRET

type ServiceClaims struct {
	Service string `json:"service"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenString string) (*ServiceClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ServiceClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*ServiceClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
