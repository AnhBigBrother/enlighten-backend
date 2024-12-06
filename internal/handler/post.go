package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/dto"
	"github.com/AnhBigBrother/enlighten-backend/internal/models"
	"github.com/AnhBigBrother/enlighten-backend/pkg/parser"
	"github.com/AnhBigBrother/enlighten-backend/pkg/resp"
	"github.com/google/uuid"
)

type PostApi struct {
	DB *database.Queries
}

func NewPostApi() PostApi {
	return PostApi{
		DB: cfg.DBQueries,
	}
}

func (postApi *PostApi) GetAllPost(w http.ResponseWriter, r *http.Request) {
	posts, err := postApi.DB.GetAllPosts(r.Context())
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	ret := []models.Post{}
	for _, p := range posts {
		ret = append(ret, models.FormatDatabaseGetAllPostsRow(p))
	}

	resp.Json(w, 200, ret)
}

func (postApi *PostApi) CreatePost(w http.ResponseWriter, r *http.Request) {
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
	post, err := postApi.DB.CreatePost(r.Context(), database.CreatePostParams{
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

func (postApi *PostApi) GetPostById(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	post, err := postApi.DB.GetPostById(r.Context(), postUUID)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, models.FormatDatabaseGetPostByIdRow(post))
}

func (postApi *PostApi) UpVotePost(w http.ResponseWriter, r *http.Request) {
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
	err = postApi.DB.VotePost(r.Context(), cfg.DBConnection, authorUuid, postUUID, "up")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (postApi *PostApi) DownVotePost(w http.ResponseWriter, r *http.Request) {
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
	err = postApi.DB.VotePost(r.Context(), cfg.DBConnection, authorUuid, postUUID, "down")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (postApi *PostApi) CheckVoted(w http.ResponseWriter, r *http.Request) {
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
	pv, err := postApi.DB.GetPostVotes(r.Context(), database.GetPostVotesParams{
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

func (postApi *PostApi) AddPostComment(w http.ResponseWriter, r *http.Request) {
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
	com, err := postApi.DB.AddComment(r.Context(), cfg.DBConnection, params.Comment, authorUuid, postUUID, uuid.NullUUID{})
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

func (postApi *PostApi) GetAllComments(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("post_id")
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	comments, err := postApi.DB.GetPostComments(r.Context(), postUUID)
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
