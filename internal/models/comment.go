package models

import (
	"time"

	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID     `json:"id"`
	Comment   string        `json:"comment"`
	UserID    uuid.UUID     `json:"user_id"`
	PostID    uuid.UUID     `json:"post_id"`
	CommentID uuid.NullUUID `json:"comment_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpVoted   int32         `json:"up_voted"`
	DownVoted int32         `json:"down_voted"`
	UserName  string        `json:"user_name"`
	UserEmail string        `json:"user_email"`
	UserImage string        `json:"user_image"`
}

func FormatDatabaseComment(dbCom database.Comment) Comment {
	return Comment{
		ID:        dbCom.ID,
		Comment:   dbCom.Comment,
		UserID:    dbCom.AuthorID,
		PostID:    dbCom.PostID,
		CommentID: dbCom.ParentCommentID,
		UpVoted:   dbCom.UpVoted,
		DownVoted: dbCom.DownVoted,
		CreatedAt: dbCom.CreatedAt,
	}
}

func FormatDatabaseGetPostCommentsRow(dbCom database.GetPostCommentsRow) Comment {
	return Comment{
		ID:        dbCom.ID,
		Comment:   dbCom.Comment,
		UserID:    dbCom.AuthorID,
		PostID:    dbCom.PostID,
		CommentID: dbCom.ParentCommentID,
		UpVoted:   dbCom.UpVoted,
		DownVoted: dbCom.DownVoted,
		UserName:  dbCom.UserName.String,
		UserEmail: dbCom.UserEmail.String,
		UserImage: dbCom.UserImage.String,
		CreatedAt: dbCom.CreatedAt,
	}
}
