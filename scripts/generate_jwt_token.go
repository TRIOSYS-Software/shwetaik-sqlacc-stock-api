package scripts

import (
	"shwetaik-sqlacc-stock-api/internal/config"
	"shwetaik-sqlacc-stock-api/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken() {
	claims := utils.ServiceClaims{
		Service: config.Cfg.ServiceName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3600)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(utils.JWT_SECRET))
	if err != nil {
		panic(err)
	}
	println(tokenString)
}
