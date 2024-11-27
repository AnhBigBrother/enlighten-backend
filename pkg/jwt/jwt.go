package jwt

import (
	"enlighten-backend/cfg"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateJwt(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Envs.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	if _, ok := claims["typ"]; !ok {
		return nil, fmt.Errorf("invalid token")
	}
	if _, ok := claims["iat"]; !ok {
		return nil, fmt.Errorf("invalid token")
	}
	if exp, ok := claims["exp"]; !ok || exp.(int) < time.Now().Second() {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func SignJwt(payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, payload)
	tokenString, err := token.SignedString(cfg.Envs.JwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
