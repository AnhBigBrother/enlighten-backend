package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/dto"
	"github.com/AnhBigBrother/enlighten-backend/internal/models"
	"github.com/AnhBigBrother/enlighten-backend/pkg/parser"
	"github.com/AnhBigBrother/enlighten-backend/pkg/resp"
	"github.com/google/uuid"
)

type PostsHandler struct {
	DB *database.Queries
}

func NewPostsHandler() PostsHandler {
	return PostsHandler{
		DB: cfg.DBQueries,
	}
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
	if sort == "new" {
		posts, err := postsHandler.DB.GetNewPosts(r.Context(), database.GetNewPostsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			resp.Err(w, 404, err.Error())
			return
		}
		ret := []models.Post{}
		for _, p := range posts {
			ret = append(ret, models.FormatDatabaseGetNewPostsRow(p))
		}
		resp.Json(w, 200, ret)
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
		ret := []models.Post{}
		for _, p := range posts {
			ret = append(ret, models.FormatDatabaseGetTopPostsRow(p))
		}
		resp.Json(w, 200, ret)
		return
	}
	posts, err := postsHandler.DB.GetHotPosts(r.Context(), database.GetHotPostsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	ret := []models.Post{}
	for _, p := range posts {
		ret = append(ret, models.FormatDatabaseGetHotPostsRow(p))
	}
	resp.Json(w, 200, ret)
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
	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	currentTime := time.Now()
	post, err := postsHandler.DB.CreatePost(r.Context(), database.CreatePostParams{
		ID:        uuid.New(),
		Title:     params.Title,
		Content:   params.Content,
		AuthorID:  authorUuid,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, models.FormatDatabasePost(post))
}

func (postsHandler *PostsHandler) GetPostById(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	post, err := postsHandler.DB.GetPostById(r.Context(), postUUID)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, models.FormatDatabaseGetPostByIdRow(post))
}

func (postsHandler *PostsHandler) UpVotePost(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	err = postsHandler.DB.VotePost(r.Context(), cfg.DBConnection, authorUuid, postUUID, "up")
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
	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	err = postsHandler.DB.VotePost(r.Context(), cfg.DBConnection, authorUuid, postUUID, "down")
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
	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	pv, err := postsHandler.DB.GetPostVotes(r.Context(), database.GetPostVotesParams{
		PostID:  postUUID,
		VoterID: authorUuid,
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
	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
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
	com, err := postsHandler.DB.AddComment(r.Context(), cfg.DBConnection, params.Comment, authorUuid, postUUID, uuid.NullUUID{})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, models.Comment{
		ID:              com.ID,
		Comment:         com.Comment,
		AuthorId:        com.AuthorID,
		PostID:          com.PostID,
		ParentCommentID: com.ParentCommentID,
		CreatedAt:       com.CreatedAt,
	})
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
		PostID: postUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	cms := []models.Comment{}
	for _, c := range comments {
		cms = append(cms, models.FormatDatabaseGetPostCommentsRow(c))
	}
	resp.Json(w, 200, cms)
}
