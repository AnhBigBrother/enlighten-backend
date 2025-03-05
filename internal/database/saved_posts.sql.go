// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: saved_posts.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSavedPost = `-- name: CreateSavedPost :one
INSERT INTO
  saved_posts ("id", "user_id", "post_id", "created_at")
VALUES
  ($1, $2, $3, $4)
RETURNING
  id, user_id, post_id, created_at
`

type CreateSavedPostParams struct {
	ID        pgtype.UUID      `json:"id"`
	UserID    pgtype.UUID      `json:"user_id"`
	PostID    pgtype.UUID      `json:"post_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) CreateSavedPost(ctx context.Context, arg CreateSavedPostParams) (SavedPost, error) {
	row := q.db.QueryRow(ctx, createSavedPost,
		arg.ID,
		arg.UserID,
		arg.PostID,
		arg.CreatedAt,
	)
	var i SavedPost
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteSavedPost = `-- name: DeleteSavedPost :exec
DELETE FROM saved_posts
WHERE
  user_id = $1
  AND post_id = $2
`

type DeleteSavedPostParams struct {
	UserID pgtype.UUID `json:"user_id"`
	PostID pgtype.UUID `json:"post_id"`
}

func (q *Queries) DeleteSavedPost(ctx context.Context, arg DeleteSavedPostParams) error {
	_, err := q.db.Exec(ctx, deleteSavedPost, arg.UserID, arg.PostID)
	return err
}

const getAllSavedPosts = `-- name: GetAllSavedPosts :many
SELECT
  sp.saved_at, sp.id, sp.title, sp.content, sp.author_id, sp.up_voted, sp.down_voted, sp.comments_count, sp.created_at, sp.updated_at,
  u.email AS author_email,
  u.name AS author_name,
  u.image AS author_image
FROM
  (
    SELECT
      s.created_at AS saved_at,
      p.id, p.title, p.content, p.author_id, p.up_voted, p.down_voted, p.comments_count, p.created_at, p.updated_at
    FROM
      (
        SELECT
          id, user_id, post_id, created_at
        FROM
          saved_posts
        WHERE
          saved_posts.user_id = $1
      ) s
      LEFT JOIN posts p ON s.post_id = p.id
  ) sp
  LEFT JOIN users u ON sp.author_id = u.id
ORDER BY
  sp.saved_at DESC
LIMIT
  $2
OFFSET
  $3
`

type GetAllSavedPostsParams struct {
	UserID pgtype.UUID `json:"user_id"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

type GetAllSavedPostsRow struct {
	SavedAt       pgtype.Timestamp `json:"saved_at"`
	ID            pgtype.UUID      `json:"id"`
	Title         pgtype.Text      `json:"title"`
	Content       pgtype.Text      `json:"content"`
	AuthorID      pgtype.UUID      `json:"author_id"`
	UpVoted       pgtype.Int4      `json:"up_voted"`
	DownVoted     pgtype.Int4      `json:"down_voted"`
	CommentsCount pgtype.Int4      `json:"comments_count"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
	AuthorEmail   pgtype.Text      `json:"author_email"`
	AuthorName    pgtype.Text      `json:"author_name"`
	AuthorImage   pgtype.Text      `json:"author_image"`
}

func (q *Queries) GetAllSavedPosts(ctx context.Context, arg GetAllSavedPostsParams) ([]GetAllSavedPostsRow, error) {
	rows, err := q.db.Query(ctx, getAllSavedPosts, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllSavedPostsRow{}
	for rows.Next() {
		var i GetAllSavedPostsRow
		if err := rows.Scan(
			&i.SavedAt,
			&i.ID,
			&i.Title,
			&i.Content,
			&i.AuthorID,
			&i.UpVoted,
			&i.DownVoted,
			&i.CommentsCount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AuthorEmail,
			&i.AuthorName,
			&i.AuthorImage,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
