package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/pkg/token"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			next.ServeHTTP(w, r)
			return
		}
		userClaim, err := token.Parse(access_token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), cfg.CtxKeys.User, userClaim)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
