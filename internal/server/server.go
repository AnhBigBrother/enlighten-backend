package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/handler"
	"github.com/AnhBigBrother/enlighten-backend/internal/middleware"
	"github.com/rs/cors"
)

func RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	userApi := handler.NewUserApi()
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("POST /signup", userApi.SignUp)
	userRouter.HandleFunc("POST /signin", userApi.SignIn)
	userRouter.HandleFunc("POST /signout", middleware.Auth(userApi.SignOut))
	userRouter.HandleFunc("GET /me", middleware.Auth(userApi.GetMe))
	userRouter.HandleFunc("POST /me", middleware.Auth(userApi.UpdateMe))
	userRouter.HandleFunc("DELETE /me", middleware.Auth(userApi.DeleteMe))
	userRouter.HandleFunc("GET /me/session", middleware.Auth(userApi.GetSesion))
	userRouter.HandleFunc("GET /me/access_token", userApi.GetAccessToken)

	oauthApi := handler.NewOauthApi()
	oauthRouter := http.NewServeMux()
	oauthRouter.HandleFunc("POST /google", oauthApi.OauthGoogle)
	oauthRouter.HandleFunc("POST /github", oauthApi.OauthGithub)
	oauthRouter.HandleFunc("POST /microsoft", oauthApi.OauthMicrosoft)
	oauthRouter.HandleFunc("POST /discord", oauthApi.OauthDiscord)

	router.Handle("/user/", http.StripPrefix("/user", userRouter))
	router.Handle("/oauth/", http.StripPrefix("/oauth", oauthRouter))

	return cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(router)
}

func NewServer() *http.Server {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
