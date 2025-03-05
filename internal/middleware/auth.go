package middleware

import (
	"context"
	"net/http"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	token "github.com/AnhBigBrother/enlighten-backend/internal/pkg/jwt-token"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access_token := r.Header.Get("Authorization")
		if access_token == "" {
			cookie, err := r.Cookie("access_token")
			if err == nil {
				access_token = cookie.Value
			}
		}
		if access_token == "" {
			next.ServeHTTP(w, r)
			return
		}
		userClaim, err := token.ParseAndVerify(access_token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), cfg.CtxKeys.User, userClaim)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
