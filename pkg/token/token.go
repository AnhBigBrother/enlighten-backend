package token

import (
	"fmt"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Image string `json:"image"`
	jwt.RegisteredClaims
}

func Parse(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}
	if exp.Compare(time.Now()) < 0 {
		return nil, fmt.Errorf("token expired")
	}
	return claims, nil
}

func Sign(payload Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
