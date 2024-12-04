// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: comments.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createComment = `-- name: CreateComment :exec
INSERT INTO
  COMMENTS (
    "id",
    "comment",
    "author_id",
    "post_id",
    "parent_comment_id",
    "created_at"
  )
VALUES
  ($1, $2, $3, $4, $5, $6)
`

type CreateCommentParams struct {
	ID              uuid.UUID
	Comment         string
	AuthorID        uuid.UUID
	PostID          uuid.UUID
	ParentCommentID uuid.NullUUID
	CreatedAt       time.Time
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) error {
	_, err := q.db.ExecContext(ctx, createComment,
		arg.ID,
		arg.Comment,
		arg.AuthorID,
		arg.PostID,
		arg.ParentCommentID,
		arg.CreatedAt,
	)
	return err
}

const getCommentsReplies = `-- name: GetCommentsReplies :many
SELECT
  pc.id, pc.comment, pc.author_id, pc.post_id, pc.parent_comment_id, pc.up_voted, pc.down_voted, pc.created_at,
  u.email AS user_email,
  u.name AS user_name,
  u.image AS user_image
FROM
  (
    SELECT
      id, comment, author_id, post_id, parent_comment_id, up_voted, down_voted, created_at
    FROM
      COMMENTS c
    WHERE
      c.post_id = $1
      AND c.parent_comment_id = $2
  ) pc
  LEFT JOIN users u ON pc.author_id = u.id
`

type GetCommentsRepliesParams struct {
	PostID          uuid.UUID
	ParentCommentID uuid.NullUUID
}

type GetCommentsRepliesRow struct {
	ID              uuid.UUID
	Comment         string
	AuthorID        uuid.UUID
	PostID          uuid.UUID
	ParentCommentID uuid.NullUUID
	UpVoted         int32
	DownVoted       int32
	CreatedAt       time.Time
	UserEmail       sql.NullString
	UserName        sql.NullString
	UserImage       sql.NullString
}

func (q *Queries) GetCommentsReplies(ctx context.Context, arg GetCommentsRepliesParams) ([]GetCommentsRepliesRow, error) {
	rows, err := q.db.QueryContext(ctx, getCommentsReplies, arg.PostID, arg.ParentCommentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommentsRepliesRow
	for rows.Next() {
		var i GetCommentsRepliesRow
		if err := rows.Scan(
			&i.ID,
			&i.Comment,
			&i.AuthorID,
			&i.PostID,
			&i.ParentCommentID,
			&i.UpVoted,
			&i.DownVoted,
			&i.CreatedAt,
			&i.UserEmail,
			&i.UserName,
			&i.UserImage,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
