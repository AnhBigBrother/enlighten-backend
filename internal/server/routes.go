package server

import (
	"net/http"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/guard"
	"github.com/AnhBigBrother/enlighten-backend/internal/middleware"
	"github.com/AnhBigBrother/enlighten-backend/internal/service"
	"github.com/rs/cors"
)

func RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	userService := service.NewUserService()
	oauthService := service.NewOauthService()
	postService := service.NewPostsService()
	commentService := service.NewCommentsService()
	sudokuService := service.NewSudokuService()

	authRouter := http.NewServeMux()
	authRouter.HandleFunc("POST /signup", userService.SignUp)
	authRouter.HandleFunc("POST /signin", userService.SignIn)
	authRouter.HandleFunc("POST /signout", guard.Auth(userService.SignOut))
	authRouter.HandleFunc("GET /google", oauthService.HandleGoogleOauth)
	authRouter.HandleFunc("GET /github", oauthService.HandleGithubOauth)
	authRouter.HandleFunc("GET /microsoft", oauthService.HandleMicrosoftOauth)
	authRouter.HandleFunc("GET /discord", oauthService.HandleDiscordOauth)

	signedUsersRouter := http.NewServeMux()
	signedUsersRouter.HandleFunc("GET /", guard.Auth(userService.GetMe))
	signedUsersRouter.HandleFunc("PATCH /", guard.Auth(userService.UpdateMe))
	signedUsersRouter.HandleFunc("DELETE /", guard.Auth(userService.DeleteMe))
	signedUsersRouter.HandleFunc("GET /session", guard.Auth(userService.GetSesion))
	signedUsersRouter.HandleFunc("GET /access_token", userService.GetAccessToken)
	signedUsersRouter.HandleFunc("GET /followed", guard.Auth(userService.GetFollowedAuthor))
	signedUsersRouter.HandleFunc("GET /overview", guard.Auth(userService.GetMyOverview))
	signedUsersRouter.HandleFunc("GET /posts", guard.Auth(userService.GetMyPosts))

	usersRouter := http.NewServeMux()
	usersRouter.HandleFunc("GET /all", userService.GetTopAuthor)
	usersRouter.HandleFunc("GET /{user_id}/overview", userService.GetOverview)
	usersRouter.HandleFunc("GET /{user_id}/posts", userService.GetPosts)
	usersRouter.HandleFunc("POST /{user_id}/follows", guard.Auth(userService.Follow))
	usersRouter.HandleFunc("DELETE /{user_id}/follows", guard.Auth(userService.UnFollow))
	usersRouter.HandleFunc("GET /{user_id}/follows/check", guard.Auth(userService.CheckFollowed))

	postsRouter := http.NewServeMux()
	postsRouter.HandleFunc("GET /", postService.GetFollowedPosts)
	postsRouter.HandleFunc("GET /all", postService.GetAllPosts)
	postsRouter.HandleFunc("GET /saved", guard.Auth(postService.GetAllSavedPost))
	postsRouter.HandleFunc("POST /create", guard.Auth(postService.CreatePost))
	postsRouter.HandleFunc("GET /{post_id}", postService.GetPostById)
	postsRouter.HandleFunc("GET /{post_id}/check", guard.Auth(postService.CheckPost))
	postsRouter.HandleFunc("POST /{post_id}/save", guard.Auth(postService.SavePost))
	postsRouter.HandleFunc("DELETE /{post_id}/save", guard.Auth(postService.UnSavePost))
	postsRouter.HandleFunc("POST /{post_id}/vote/up", guard.Auth(postService.UpVotePost))
	postsRouter.HandleFunc("POST /{post_id}/vote/down", guard.Auth(postService.DownVotePost))
	postsRouter.HandleFunc("GET /{post_id}/comments", postService.GetPostComments)
	postsRouter.HandleFunc("POST /{post_id}/comments", guard.Auth(postService.AddPostComment))
	postsRouter.HandleFunc("GET /{post_id}/comments/{comment_id}", commentService.GetCommentReplies)
	postsRouter.HandleFunc("GET /{post_id}/comments/{comment_id}/check", guard.Auth(commentService.CheckComment))
	postsRouter.HandleFunc("POST /{post_id}/comments/{comment_id}/vote/up", guard.Auth(commentService.UpVoteComment))
	postsRouter.HandleFunc("POST /{post_id}/comments/{comment_id}/vote/down", guard.Auth(commentService.DownVoteComment))
	postsRouter.HandleFunc("POST /{post_id}/comments/{comment_id}/reply", guard.Auth(commentService.ReplyComment))

	gameRouter := http.NewServeMux()
	gameRouter.HandleFunc("GET /sudoku", sudokuService.GenerateSudoku)
	gameRouter.HandleFunc("POST /sudoku", sudokuService.SolveSudoku)
	gameRouter.HandleFunc("POST /sudoku/check", sudokuService.CheckSolvable)

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

	// stack := middleware.CreateStack(corsMiddleware.Handler, middleware.Auth)
	stack := middleware.CreateStack(corsMiddleware.Handler, middleware.Logging, middleware.Auth)

	return stack(router)
}
