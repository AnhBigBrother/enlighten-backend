package server

import (
	"enlighten-backend/cfg"
	"enlighten-backend/internal/handler"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func RegisterRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	userApi := handler.NewUserApi()
	userRouter := chi.NewRouter()
	userRouter.Post("/signup", userApi.SignUp)

	router.Mount("/user", userRouter)

	return router
}

func NewServer() *http.Server {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.Envs.Port),
		Handler:      RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
