package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/dto"
	"github.com/AnhBigBrother/enlighten-backend/pkg/parser"
	"github.com/AnhBigBrother/enlighten-backend/pkg/resp"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostsHandler struct {
	DB *database.Queries
}

func NewPostsHandler() PostsHandler {
	return PostsHandler{
		DB: cfg.DBQueries,
	}
}

func (postsHandler *PostsHandler) GetFollowedPosts(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	if !ok {
		postsHandler.GetAllPosts(w, r)
		return
	}
	userId := user["jti"].(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		resp.Err(w, 400, "invalid userId")
		return
	}
	queryParams := r.URL.Query()
	sort, limitStr, offsetStr := queryParams.Get("sort"), queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	if sort == "hot" {
		posts, err := postsHandler.DB.GetHotFollowedPosts(r.Context(), database.GetHotFollowedPostsParams{
			FollowerID: pgtype.UUID{Bytes: userUUID, Valid: true},
			Limit:      int32(limit),
			Offset:     int32(offset),
		})
		if err != nil {
			resp.Err(w, 404, err.Error())
			return
		}
		resp.Json(w, 200, posts)
		return
	}
	if sort == "top" {
		posts, err := postsHandler.DB.GetTopFollowedPosts(r.Context(), database.GetTopFollowedPostsParams{
			FollowerID: pgtype.UUID{Bytes: userUUID, Valid: true},
			Limit:      int32(limit),
			Offset:     int32(offset),
		})
		if err != nil {
			resp.Err(w, 404, err.Error())
			return
		}
		resp.Json(w, 200, posts)
		return
	}
	posts, err := postsHandler.DB.GetNewFollowedPosts(r.Context(), database.GetNewFollowedPostsParams{
		FollowerID: pgtype.UUID{Bytes: userUUID, Valid: true},
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, posts)
}

func (postsHandler *PostsHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	sort, limitStr, offsetStr := queryParams.Get("sort"), queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	if sort == "hot" {
		posts, err := postsHandler.DB.GetHotPosts(r.Context(), database.GetHotPostsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			resp.Err(w, 404, err.Error())
			return
		}
		resp.Json(w, 200, posts)
		return
	}
	if sort == "top" {
		posts, err := postsHandler.DB.GetTopPosts(r.Context(), database.GetTopPostsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			resp.Err(w, 404, err.Error())
			return
		}
		resp.Json(w, 200, posts)
		return
	}
	posts, err := postsHandler.DB.GetNewPosts(r.Context(), database.GetNewPostsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, posts)
}

func (postsHandler *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	params := dto.CretaePostDTO{}
	err := parser.ParseBody(r.Body, &params)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	err = params.ValidateInput()
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	session, _ := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	currentTime := time.Now()
	post, err := postsHandler.DB.CreatePost(r.Context(), database.CreatePostParams{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Title:     params.Title,
		Content:   params.Content,
		AuthorID:  pgtype.UUID{Bytes: authorUuid, Valid: true},
		CreatedAt: pgtype.Timestamp{Time: currentTime, InfinityModifier: pgtype.Finite, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: currentTime, InfinityModifier: pgtype.Finite, Valid: true},
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, post)
}

func (postsHandler *PostsHandler) GetPostById(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	post, err := postsHandler.DB.GetPostById(r.Context(), pgtype.UUID{Bytes: postUUID, Valid: true})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, post)
}

func (postsHandler *PostsHandler) UpVotePost(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	session, _ := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	err = postsHandler.DB.VotePost(r.Context(), cfg.DBConnection, pgtype.UUID{Bytes: authorUuid, Valid: true}, pgtype.UUID{Bytes: postUUID, Valid: true}, "up")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (postsHandler *PostsHandler) DownVotePost(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	session, _ := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	err = postsHandler.DB.VotePost(r.Context(), cfg.DBConnection, pgtype.UUID{Bytes: authorUuid, Valid: true}, pgtype.UUID{Bytes: postUUID, Valid: true}, "down")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (postsHandler *PostsHandler) CheckVoted(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	session, _ := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	pv, err := postsHandler.DB.GetPostVotes(r.Context(), database.GetPostVotesParams{
		PostID:  pgtype.UUID{Bytes: postUUID, Valid: true},
		VoterID: pgtype.UUID{Bytes: authorUuid, Valid: true},
	})
	if err != nil {
		resp.Json(w, 200, struct {
			Voted string `json:"voted"`
		}{Voted: "none"})
		return
	}
	resp.Json(w, 200, struct {
		Voted string `json:"voted"`
	}{Voted: string(pv.Voted)})
}

func (postsHandler *PostsHandler) AddPostComment(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Comment string `json:"comment"`
	}{}
	err := parser.ParseBody(r.Body, &params)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	if params.Comment == "" {
		resp.Err(w, 400, "comment is required")
		return
	}
	session, _ := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	comment, err := postsHandler.DB.AddComment(r.Context(), cfg.DBConnection, params.Comment, pgtype.UUID{Bytes: authorUuid, Valid: true}, pgtype.UUID{Bytes: postUUID, Valid: true}, pgtype.UUID{Valid: false})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, comment)
}

func (postsHandler *PostsHandler) GetPostComments(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limitStr, offsetStr := queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	comments, err := postsHandler.DB.GetPostComments(r.Context(), database.GetPostCommentsParams{
		PostID: pgtype.UUID{Bytes: postUUID, Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, comments)
}
