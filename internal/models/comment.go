package models

import (
	"time"

	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID     `json:"id"`
	Comment         string        `json:"comment"`
	AuthorId        uuid.UUID     `json:"author_id"`
	PostID          uuid.UUID     `json:"post_id"`
	ParentCommentID uuid.NullUUID `json:"parent_comment_id"`
	CreatedAt       time.Time     `json:"created_at"`
	UpVoted         int32         `json:"up_voted"`
	DownVoted       int32         `json:"down_voted"`
	AuthorName      string        `json:"author_name"`
	AuthorEmail     string        `json:"author_email"`
	AuthorImage     string        `json:"author_image"`
}

func FormatDatabaseGetCommentsRepliesRow(dbCom database.GetCommentsRepliesRow) Comment {
	return Comment{
		ID:              dbCom.ID,
		Comment:         dbCom.Comment,
		AuthorId:        dbCom.AuthorID,
		PostID:          dbCom.PostID,
		ParentCommentID: dbCom.ParentCommentID,
		UpVoted:         dbCom.UpVoted,
		DownVoted:       dbCom.DownVoted,
		AuthorName:      dbCom.UserName.String,
		AuthorEmail:     dbCom.UserEmail.String,
		AuthorImage:     dbCom.UserImage.String,
		CreatedAt:       dbCom.CreatedAt,
	}
}

func FormatDatabaseGetPostCommentsRow(dbCom database.GetPostCommentsRow) Comment {
	return Comment{
		ID:              dbCom.ID,
		Comment:         dbCom.Comment,
		AuthorId:        dbCom.AuthorID,
		PostID:          dbCom.PostID,
		ParentCommentID: dbCom.ParentCommentID,
		UpVoted:         dbCom.UpVoted,
		DownVoted:       dbCom.DownVoted,
		AuthorName:      dbCom.AuthorName.String,
		AuthorEmail:     dbCom.AuthorEmail.String,
		AuthorImage:     dbCom.AuthorImage.String,
		CreatedAt:       dbCom.CreatedAt,
	}
}
