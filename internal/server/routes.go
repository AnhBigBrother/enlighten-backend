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
	userRouter.HandleFunc("/signout", middleware.Auth(userApi.SignOut))
	userRouter.HandleFunc("GET /me", middleware.Auth(userApi.GetMe))
	userRouter.HandleFunc("POST /me", middleware.Auth(userApi.UpdateMe))
	userRouter.HandleFunc("DELETE /me", middleware.Auth(userApi.DeleteMe))
	userRouter.HandleFunc("GET /me/session", middleware.Auth(userApi.GetSesion))
	userRouter.HandleFunc("GET /me/access_token", userApi.GetAccessToken)

	oauthApi := handler.NewOauthApi()
	oauthRouter := http.NewServeMux()
	oauthRouter.HandleFunc("GET /google", oauthApi.HandleGoogleOauth)
	oauthRouter.HandleFunc("GET /github", oauthApi.HandleGithubOauth)
	oauthRouter.HandleFunc("GET /microsoft", oauthApi.HandleMicrosoftOauth)
	oauthRouter.HandleFunc("GET /discord", oauthApi.HandleDiscordOauth)

	postApi := handler.NewPostApi()
	commentApi := handler.NewCommentApi()
	postRouter := http.NewServeMux()
	postRouter.HandleFunc("GET /", postApi.GetAllPosts)
	postRouter.HandleFunc("POST /create", middleware.Auth(postApi.CreatePost))
	postRouter.HandleFunc("GET /{post_id}", postApi.GetPostById)
	postRouter.HandleFunc("GET /{post_id}/checkvoted", middleware.Auth(postApi.CheckVoted))
	postRouter.HandleFunc("POST /{post_id}/upvote", middleware.Auth(postApi.UpVotePost))
	postRouter.HandleFunc("POST /{post_id}/downvote", middleware.Auth(postApi.DownVotePost))
	postRouter.HandleFunc("GET /{post_id}/comment", postApi.GetPostComments)
	postRouter.HandleFunc("POST /{post_id}/comment", middleware.Auth(postApi.AddPostComment))
	postRouter.HandleFunc("GET /{post_id}/comment/{comment_id}", commentApi.GetCommentReplies)
	postRouter.HandleFunc("GET /{post_id}/comment/{comment_id}/checkvoted", middleware.Auth(commentApi.CheckVoted))
	postRouter.HandleFunc("POST /{post_id}/comment/{comment_id}/upvote", middleware.Auth(commentApi.UpVoteComment))
	postRouter.HandleFunc("POST /{post_id}/comment/{comment_id}/downvote", middleware.Auth(commentApi.DownVoteComment))
	postRouter.HandleFunc("POST /{post_id}/comment/{comment_id}/reply", middleware.Auth(commentApi.ReplyComment))

	router.Handle("/api/v1/user/", http.StripPrefix("/api/v1/user", userRouter))
	router.Handle("/api/v1/oauth/", http.StripPrefix("/api/v1/oauth", oauthRouter))
	router.Handle("/api/v1/post/", http.StripPrefix("/api/v1/post", postRouter))

	return cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(router)
}
