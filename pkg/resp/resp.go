package resp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
)

func Json(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error when Marshal json response:", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func Err(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Server error:", msg)
	}
	Json(w, code, struct {
		Error string `json:"error"`
	}{Error: msg})
}

func SetCookie(w http.ResponseWriter, key, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		MaxAge:   cfg.CookieAge,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
}

func DeleteCookie(w http.ResponseWriter, key string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}
