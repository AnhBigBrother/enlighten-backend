package server

import (
	"net/http"

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
	router.Handle("/user/", http.StripPrefix("/user", userRouter))

	oauthApi := handler.NewOauthApi()
	oauthRouter := http.NewServeMux()
	oauthRouter.HandleFunc("POST /google", oauthApi.OauthGoogle)
	oauthRouter.HandleFunc("POST /github", oauthApi.OauthGithub)
	oauthRouter.HandleFunc("POST /microsoft", oauthApi.OauthMicrosoft)
	oauthRouter.HandleFunc("POST /discord", oauthApi.OauthDiscord)
	router.Handle("/oauth/", http.StripPrefix("/oauth", oauthRouter))

	postApi := handler.NewPostApi()
	postRouter := http.NewServeMux()
	postRouter.HandleFunc("GET /", postApi.GetAllPost)
	postRouter.HandleFunc("POST /", middleware.Auth(postApi.CreatePost))
	postRouter.HandleFunc("GET /{id}", postApi.GetPostById)
	postRouter.HandleFunc("POST /{id}/upvote", middleware.Auth(postApi.UpVotePost))
	postRouter.HandleFunc("POST /{id}/downvote", middleware.Auth(postApi.DownVotePost))
	postRouter.HandleFunc("GET /{id}/comment", postApi.GetAllComments)
	postRouter.HandleFunc("POST /{id}/comment", middleware.Auth(postApi.AddPostComment))
	router.Handle("/post/", http.StripPrefix("/post", postRouter))

	return cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(router)
}
