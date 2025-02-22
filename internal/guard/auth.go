package guard

import (
	"net/http"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/resp"
)

func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
		if !ok {
			resp.Err(w, 401, "unauthorized: access_token failed")
			return
		}
		handler(w, r)
	}
}
