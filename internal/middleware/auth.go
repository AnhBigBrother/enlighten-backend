package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AnhBigBrother/enlighten-backend/pkg/token"
)

func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		access_token := ""
		cookie, err := r.Cookie("access_token")
		if err == nil {
			access_token = cookie.Value
		} else {
			reqToken := r.Header.Get("Authorization")
			splitToken := strings.Split(reqToken, " ")
			if len(splitToken) == 2 && splitToken[0] == "Bearer" {
				access_token = splitToken[1]
			}
		}
		if access_token == "" {
			handler(w, r)
			return
		}
		userClaim, err := token.Parse(access_token)
		if err != nil {
			handler(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "user", userClaim)
		handler(w, r.WithContext(ctx))
	}
}
