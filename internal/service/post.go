package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/dto"

	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/resp"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostsService struct {
	DB *database.Queries
}

func NewPostsService() PostsService {
	return PostsService{
		DB: cfg.DBQueries,
	}
}

func (postsService *PostsService) GetFollowedPosts(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	if !ok {
		postsService.GetAllPosts(w, r)
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
		limit = 5
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	if sort == "hot" {
		posts, err := postsService.DB.GetHotFollowedPosts(r.Context(), database.GetHotFollowedPostsParams{
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
		posts, err := postsService.DB.GetTopFollowedPosts(r.Context(), database.GetTopFollowedPostsParams{
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
	posts, err := postsService.DB.GetNewFollowedPosts(r.Context(), database.GetNewFollowedPostsParams{
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

func (postsService *PostsService) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	sort, limitStr, offsetStr := queryParams.Get("sort"), queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	if sort == "hot" {
		posts, err := postsService.DB.GetHotPosts(r.Context(), database.GetHotPostsParams{
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
		posts, err := postsService.DB.GetTopPosts(r.Context(), database.GetTopPostsParams{
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
	posts, err := postsService.DB.GetNewPosts(r.Context(), database.GetNewPostsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, posts)
}

func (postsService *PostsService) CreatePost(w http.ResponseWriter, r *http.Request) {
	params := dto.CretaePostDTO{}
	err := json.NewDecoder(r.Body).Decode(&params)
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
	post, err := postsService.DB.CreatePost(r.Context(), database.CreatePostParams{
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

func (postsService *PostsService) GetPostById(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	post, err := postsService.DB.GetPostById(r.Context(), pgtype.UUID{Bytes: postUUID, Valid: true})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, post)
}

func (postsService *PostsService) UpVotePost(w http.ResponseWriter, r *http.Request) {
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
	err = postsService.DB.VotePost(r.Context(), cfg.DBConnection, pgtype.UUID{Bytes: authorUuid, Valid: true}, pgtype.UUID{Bytes: postUUID, Valid: true}, "up")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (postsService *PostsService) DownVotePost(w http.ResponseWriter, r *http.Request) {
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
	err = postsService.DB.VotePost(r.Context(), cfg.DBConnection, pgtype.UUID{Bytes: authorUuid, Valid: true}, pgtype.UUID{Bytes: postUUID, Valid: true}, "down")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (postsService *PostsService) SavePost(w http.ResponseWriter, r *http.Request) {
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
	savedPost, err := postsService.DB.CreateSavedPost(r.Context(), database.CreateSavedPostParams{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:    pgtype.UUID{Bytes: authorUuid, Valid: true},
		PostID:    pgtype.UUID{Bytes: postUUID, Valid: true},
		CreatedAt: pgtype.Timestamp{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, savedPost)
}

func (postsService *PostsService) UnSavePost(w http.ResponseWriter, r *http.Request) {
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
	err = postsService.DB.DeleteSavedPost(r.Context(), database.DeleteSavedPostParams{
		UserID: pgtype.UUID{Bytes: authorUuid, Valid: true},
		PostID: pgtype.UUID{Bytes: postUUID, Valid: true},
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "Success"})
}

func (postsService *PostsService) GetAllSavedPost(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limitStr, offsetStr := queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	session, _ := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	userId, _ := session["jti"].(string)
	userUuid, err := uuid.Parse(userId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	savedPosts, err := postsService.DB.GetAllSavedPosts(r.Context(), database.GetAllSavedPostsParams{
		UserID: pgtype.UUID{Bytes: userUuid, Valid: true},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, savedPosts)
}

func (postsService *PostsService) CheckPost(w http.ResponseWriter, r *http.Request) {
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
	checked, err := postsService.DB.CheckPostInteracted(r.Context(), database.CheckPostInteractedParams{
		VoterID: pgtype.UUID{Bytes: authorUuid, Valid: true},
		PostID:  pgtype.UUID{Bytes: postUUID, Valid: true},
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 200, checked)
}

func (postsService *PostsService) AddPostComment(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Comment string `json:"comment"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&params)
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
	comment, err := postsService.DB.AddComment(r.Context(), cfg.DBConnection, params.Comment, pgtype.UUID{Bytes: authorUuid, Valid: true}, pgtype.UUID{Bytes: postUUID, Valid: true}, pgtype.UUID{Valid: false})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, comment)
}

func (postsService *PostsService) GetPostComments(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limitStr, offsetStr := queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
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
	comments, err := postsService.DB.GetPostComments(r.Context(), database.GetPostCommentsParams{
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
