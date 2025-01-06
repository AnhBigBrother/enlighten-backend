package server

import (
	"net/http"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/guard"
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
	sudokuHandler := handler.NewSudokuHandler()

	authRouter := http.NewServeMux()
	authRouter.HandleFunc("POST /signup", usersHandler.SignUp)
	authRouter.HandleFunc("POST /signin", usersHandler.SignIn)
	authRouter.HandleFunc("/signout", guard.Auth(usersHandler.SignOut))
	authRouter.HandleFunc("GET /google", oauthHandler.HandleGoogleOauth)
	authRouter.HandleFunc("GET /github", oauthHandler.HandleGithubOauth)
	authRouter.HandleFunc("GET /microsoft", oauthHandler.HandleMicrosoftOauth)
	authRouter.HandleFunc("GET /discord", oauthHandler.HandleDiscordOauth)

	signedUsersRouter := http.NewServeMux()
	signedUsersRouter.HandleFunc("GET /", guard.Auth(usersHandler.GetMe))
	signedUsersRouter.HandleFunc("POST /", guard.Auth(usersHandler.UpdateMe))
	signedUsersRouter.HandleFunc("DELETE /", guard.Auth(usersHandler.DeleteMe))
	signedUsersRouter.HandleFunc("GET /session", guard.Auth(usersHandler.GetSesion))
	signedUsersRouter.HandleFunc("GET /access_token", usersHandler.GetAccessToken)
	signedUsersRouter.HandleFunc("GET /followed", guard.Auth(usersHandler.GetFollowedAuthor))
	signedUsersRouter.HandleFunc("GET /overview", guard.Auth(usersHandler.GetMyOverview))
	signedUsersRouter.HandleFunc("GET /posts", guard.Auth(usersHandler.GetMyPosts))

	usersRouter := http.NewServeMux()
	usersRouter.HandleFunc("GET /all", usersHandler.GetTopAuthor)
	usersRouter.HandleFunc("GET /{user_id}/overview", usersHandler.GetOverview)
	usersRouter.HandleFunc("GET /{user_id}/posts", usersHandler.GetPosts)
	usersRouter.HandleFunc("POST /{user_id}/follows", guard.Auth(usersHandler.Follow))
	usersRouter.HandleFunc("DELETE /{user_id}/follows", guard.Auth(usersHandler.UnFollow))
	usersRouter.HandleFunc("GET /{user_id}/follows/check", guard.Auth(usersHandler.CheckFollowed))

	postsRouter := http.NewServeMux()
	postsRouter.HandleFunc("GET /", postsHandler.GetFollowedPosts)
	postsRouter.HandleFunc("GET /all", postsHandler.GetAllPosts)
	postsRouter.HandleFunc("POST /create", guard.Auth(postsHandler.CreatePost))
	postsRouter.HandleFunc("GET /{post_id}", postsHandler.GetPostById)
	postsRouter.HandleFunc("GET /{post_id}/vote/check", guard.Auth(postsHandler.CheckVoted))
	postsRouter.HandleFunc("POST /{post_id}/vote/up", guard.Auth(postsHandler.UpVotePost))
	postsRouter.HandleFunc("POST /{post_id}/vote/down", guard.Auth(postsHandler.DownVotePost))
	postsRouter.HandleFunc("GET /{post_id}/comments", postsHandler.GetPostComments)
	postsRouter.HandleFunc("POST /{post_id}/comments", guard.Auth(postsHandler.AddPostComment))
	postsRouter.HandleFunc("GET /{post_id}/comments/{comment_id}", commentsHandler.GetCommentReplies)
	postsRouter.HandleFunc("GET /{post_id}/comments/{comment_id}/vote/check", guard.Auth(commentsHandler.CheckVoted))
	postsRouter.HandleFunc("POST /{post_id}/comments/{comment_id}/vote/up", guard.Auth(commentsHandler.UpVoteComment))
	postsRouter.HandleFunc("POST /{post_id}/comments/{comment_id}/vote/down", guard.Auth(commentsHandler.DownVoteComment))
	postsRouter.HandleFunc("POST /{post_id}/comments/{comment_id}/reply", guard.Auth(commentsHandler.ReplyComment))

	gameRouter := http.NewServeMux()
	gameRouter.HandleFunc("GET /sudoku", sudokuHandler.GenerateSudoku)
	gameRouter.HandleFunc("POST /sudoku", sudokuHandler.SolveSudoku)
	gameRouter.HandleFunc("POST /sudoku/check", sudokuHandler.CheckSolvable)

	router.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", authRouter))
	router.Handle("/api/v1/me/", http.StripPrefix("/api/v1/me", signedUsersRouter))
	router.Handle("/api/v1/users/", http.StripPrefix("/api/v1/users", usersRouter))
	router.Handle("/api/v1/posts/", http.StripPrefix("/api/v1/posts", postsRouter))
	router.Handle("/api/v1/games/", http.StripPrefix("/api/v1/games", gameRouter))

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendUrl},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	stack := middleware.CreateStack(corsMiddleware.Handler, middleware.Auth)
	// stack := middleware.CreateStack(corsMiddleware.Handler, middleware.Logging, middleware.Auth)

	return stack(router)
}
