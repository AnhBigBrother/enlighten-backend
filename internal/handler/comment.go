package handler

import (
	"net/http"
	"strconv"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/models"
	"github.com/AnhBigBrother/enlighten-backend/pkg/parser"
	"github.com/AnhBigBrother/enlighten-backend/pkg/resp"
	"github.com/google/uuid"
)

type CommentsHandler struct {
	DB *database.Queries
}

func NewCommentsHandler() CommentsHandler {
	return CommentsHandler{
		DB: cfg.DBQueries,
	}
}

func (commentsHandler *CommentsHandler) GetCommentReplies(w http.ResponseWriter, r *http.Request) {
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
	parentCommentId := r.PathValue("comment_id")
	parentCommentUUID, err := uuid.Parse(parentCommentId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	replies, err := commentsHandler.DB.GetCommentsReplies(r.Context(), database.GetCommentsRepliesParams{
		PostID:          postUUID,
		ParentCommentID: uuid.NullUUID{UUID: parentCommentUUID, Valid: true},
		Limit:           int32(limit),
		Offset:          int32(offset),
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	ret := []models.Comment{}
	for _, rep := range replies {
		ret = append(ret, models.FormatDatabaseGetCommentsRepliesRow(rep))
	}
	resp.Json(w, 200, ret)
}

func (commentsHandler *CommentsHandler) UpVoteComment(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	commentId := r.PathValue("comment_id")
	commentUUID, err := uuid.Parse(commentId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	err = commentsHandler.DB.VoteComment(r.Context(), cfg.DBConnection, authorUuid, commentUUID, "up")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (commentsHandler *CommentsHandler) DownVoteComment(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	commentId := r.PathValue("comment_id")
	commentUUID, err := uuid.Parse(commentId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	err = commentsHandler.DB.VoteComment(r.Context(), cfg.DBConnection, authorUuid, commentUUID, "down")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (commentsHandler *CommentsHandler) CheckVoted(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	authorId, _ := session["jti"].(string)
	authorUuid, err := uuid.Parse(authorId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	commentId := r.PathValue("comment_id")
	commentUUID, err := uuid.Parse(commentId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	cv, err := commentsHandler.DB.GetCommentVotes(r.Context(), database.GetCommentVotesParams{
		CommentID: commentUUID,
		VoterID:   authorUuid,
	})
	if err != nil {
		resp.Json(w, 200, struct {
			Voted string `json:"voted"`
		}{Voted: "none"})
		return
	}
	resp.Json(w, 200, struct {
		Voted string `json:"voted"`
	}{Voted: string(cv.Voted)})
}

func (commentsHandler *CommentsHandler) ReplyComment(w http.ResponseWriter, r *http.Request) {
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
	parentCommentId := r.PathValue("comment_id")
	parentCommentUUID, err := uuid.Parse(parentCommentId)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	params := struct {
		Reply string `json:"reply"`
	}{}
	err = parser.ParseBody(r.Body, &params)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	if params.Reply == "" {
		resp.Err(w, 400, "reply is required")
		return
	}
	com, err := commentsHandler.DB.AddComment(r.Context(), cfg.DBConnection, params.Reply, authorUuid, postUUID, uuid.NullUUID{UUID: parentCommentUUID, Valid: true})
	if err != nil {
		resp.Err(w, 400, "reply is required")
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
