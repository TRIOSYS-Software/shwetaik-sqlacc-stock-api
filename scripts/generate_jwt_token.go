package main

import (
	"shwetaik-sqlacc-stock-api/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	claims := utils.ServiceClaims{
		Service: "stock-api",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(utils.JWT_SECRET))
	if err != nil {
		panic(err)
	}
	println(tokenString)
}
