package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/resp"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CommentsService struct {
	DB *database.Queries
}

func NewCommentsService() CommentsService {
	return CommentsService{
		DB: cfg.DBQueries,
	}
}

func (commentsService *CommentsService) GetCommentReplies(w http.ResponseWriter, r *http.Request) {
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
	replies, err := commentsService.DB.GetCommentsReplies(r.Context(), database.GetCommentsRepliesParams{
		PostID: pgtype.UUID{
			Bytes: postUUID,
			Valid: true,
		},
		ParentCommentID: pgtype.UUID{
			Bytes: parentCommentUUID,
			Valid: true,
		},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	resp.Json(w, 200, replies)
}

func (commentsService *CommentsService) UpVoteComment(w http.ResponseWriter, r *http.Request) {
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
	err = commentsService.DB.VoteComment(r.Context(), cfg.DBConnection, pgtype.UUID{Bytes: authorUuid, Valid: true}, pgtype.UUID{Bytes: commentUUID, Valid: true}, "up")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (commentsService *CommentsService) DownVoteComment(w http.ResponseWriter, r *http.Request) {
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
	err = commentsService.DB.VoteComment(r.Context(), cfg.DBConnection, pgtype.UUID{Bytes: authorUuid, Valid: true}, pgtype.UUID{Bytes: commentUUID, Valid: true}, "down")
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (commentsService *CommentsService) CheckComment(w http.ResponseWriter, r *http.Request) {
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
	cv, err := commentsService.DB.GetCommentVotes(r.Context(), database.GetCommentVotesParams{
		CommentID: pgtype.UUID{Bytes: commentUUID, Valid: true},
		VoterID:   pgtype.UUID{Bytes: authorUuid, Valid: true},
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

func (commentsService *CommentsService) ReplyComment(w http.ResponseWriter, r *http.Request) {
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
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	if params.Reply == "" {
		resp.Err(w, 400, "reply is required")
		return
	}
	comment, err := commentsService.DB.AddComment(r.Context(), cfg.DBConnection, params.Reply, pgtype.UUID{Bytes: authorUuid, Valid: true}, pgtype.UUID{Bytes: postUUID, Valid: true}, pgtype.UUID{Bytes: parentCommentUUID, Valid: true})
	if err != nil {
		resp.Err(w, 400, "reply is required")
		return
	}
	resp.Json(w, 201, comment)
}
