package server

import (
	"net/http"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/handler"
	"github.com/AnhBigBrother/enlighten-backend/internal/middleware"
	"github.com/rs/cors"
)

func RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	usersHandler := handler.NewUsersHandler()
	oauthHandler := handler.NewOauthHandler()
	postsHandler := handler.NewPostsHandler()
	commentsHandler := handler.NewCommentsHandler()
	sudokuHandler := handler.Sudoku{}

	authRouter := http.NewServeMux()
	authRouter.HandleFunc("POST /signup", usersHandler.SignUp)
	authRouter.HandleFunc("POST /signin", usersHandler.SignIn)
	authRouter.HandleFunc("/signout", middleware.Auth(usersHandler.SignOut))
	authRouter.HandleFunc("GET /google", oauthHandler.HandleGoogleOauth)
	authRouter.HandleFunc("GET /github", oauthHandler.HandleGithubOauth)
	authRouter.HandleFunc("GET /microsoft", oauthHandler.HandleMicrosoftOauth)
	authRouter.HandleFunc("GET /discord", oauthHandler.HandleDiscordOauth)

	signedUsersRouter := http.NewServeMux()
	signedUsersRouter.HandleFunc("GET /", middleware.Auth(usersHandler.GetMe))
	signedUsersRouter.HandleFunc("POST /", middleware.Auth(usersHandler.UpdateMe))
	signedUsersRouter.HandleFunc("DELETE /", middleware.Auth(usersHandler.DeleteMe))
	signedUsersRouter.HandleFunc("GET /session", middleware.Auth(usersHandler.GetSesion))
	signedUsersRouter.HandleFunc("GET /access_token", usersHandler.GetAccessToken)

	usersRouter := http.NewServeMux()
	usersRouter.HandleFunc("GET /{user_id}/overview", usersHandler.GetOverview)
	usersRouter.HandleFunc("GET /{user_id}/posts", usersHandler.GetPosts)

	postsRouter := http.NewServeMux()
	postsRouter.HandleFunc("GET /", postsHandler.GetAllPosts)
	postsRouter.HandleFunc("POST /create", middleware.Auth(postsHandler.CreatePost))
	postsRouter.HandleFunc("GET /{post_id}", postsHandler.GetPostById)
	postsRouter.HandleFunc("GET /{post_id}/checkvoted", middleware.Auth(postsHandler.CheckVoted))
	postsRouter.HandleFunc("POST /{post_id}/upvote", middleware.Auth(postsHandler.UpVotePost))
	postsRouter.HandleFunc("POST /{post_id}/downvote", middleware.Auth(postsHandler.DownVotePost))
	postsRouter.HandleFunc("GET /{post_id}/comments", postsHandler.GetPostComments)
	postsRouter.HandleFunc("POST /{post_id}/comments", middleware.Auth(postsHandler.AddPostComment))
	postsRouter.HandleFunc("GET /{post_id}/comments/{comment_id}", commentsHandler.GetCommentReplies)
	postsRouter.HandleFunc("GET /{post_id}/comments/{comment_id}/checkvoted", middleware.Auth(commentsHandler.CheckVoted))
	postsRouter.HandleFunc("POST /{post_id}/comments/{comment_id}/upvote", middleware.Auth(commentsHandler.UpVoteComment))
	postsRouter.HandleFunc("POST /{post_id}/comments/{comment_id}/downvote", middleware.Auth(commentsHandler.DownVoteComment))
	postsRouter.HandleFunc("POST /{post_id}/comments/{comment_id}/reply", middleware.Auth(commentsHandler.ReplyComment))

	gameRouter := http.NewServeMux()
	gameRouter.HandleFunc("GET /sudoku", sudokuHandler.GenerateSudoku)
	gameRouter.HandleFunc("POST /sudoku", sudokuHandler.SolveSudoku)
	gameRouter.HandleFunc("POST /sudoku/check", sudokuHandler.CheckSolvable)

	router.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", authRouter))
	router.Handle("/api/v1/me/", http.StripPrefix("/api/v1/me", signedUsersRouter))
	router.Handle("/api/v1/users/", http.StripPrefix("/api/v1/users", usersRouter))
	router.Handle("/api/v1/posts/", http.StripPrefix("/api/v1/posts", postsRouter))
	router.Handle("/api/v1/games/", http.StripPrefix("/api/v1/games", gameRouter))

	return cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendUrl},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(router)
}
