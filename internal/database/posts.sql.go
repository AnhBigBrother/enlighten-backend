// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO
  posts (
    "id",
    "title",
    "content",
    "author_id",
    "created_at",
    "updated_at"
  )
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING
  id, title, content, author_id, up_voted, down_voted, comments_count, created_at, updated_at
`

type CreatePostParams struct {
	ID        uuid.UUID
	Title     string
	Content   string
	AuthorID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.AuthorID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.AuthorID,
		&i.UpVoted,
		&i.DownVoted,
		&i.CommentsCount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllPosts = `-- name: GetAllPosts :many
SELECT
  p.id, p.title, p.content, p.author_id, p.up_voted, p.down_voted, p.comments_count, p.created_at, p.updated_at,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  posts p
  LEFT JOIN users u ON p.author_id = u.id
ORDER BY
  p.created_at DESC
`

type GetAllPostsRow struct {
	ID            uuid.UUID
	Title         string
	Content       string
	AuthorID      uuid.UUID
	UpVoted       int32
	DownVoted     int32
	CommentsCount int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AuthorName    sql.NullString
	AuthorEmail   sql.NullString
	AuthorImage   sql.NullString
}

func (q *Queries) GetAllPosts(ctx context.Context) ([]GetAllPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllPostsRow
	for rows.Next() {
		var i GetAllPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.AuthorID,
			&i.UpVoted,
			&i.DownVoted,
			&i.CommentsCount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AuthorName,
			&i.AuthorEmail,
			&i.AuthorImage,
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

const getPostByAuthor = `-- name: GetPostByAuthor :many
SELECT
  id, title, content, author_id, up_voted, down_voted, comments_count, created_at, updated_at
FROM
  posts
WHERE
  author_id = $1
`

func (q *Queries) GetPostByAuthor(ctx context.Context, authorID uuid.UUID) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostByAuthor, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.AuthorID,
			&i.UpVoted,
			&i.DownVoted,
			&i.CommentsCount,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getPostById = `-- name: GetPostById :one
SELECT
  p.id, p.title, p.content, p.author_id, p.up_voted, p.down_voted, p.comments_count, p.created_at, p.updated_at,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  (
    SELECT
      id, title, content, author_id, up_voted, down_voted, comments_count, created_at, updated_at
    FROM
      posts p1
    WHERE
      p1.id = $1
    LIMIT
      1
  ) p
  LEFT JOIN users u ON p.author_id = u.id
`

type GetPostByIdRow struct {
	ID            uuid.UUID
	Title         string
	Content       string
	AuthorID      uuid.UUID
	UpVoted       int32
	DownVoted     int32
	CommentsCount int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AuthorName    sql.NullString
	AuthorEmail   sql.NullString
	AuthorImage   sql.NullString
}

func (q *Queries) GetPostById(ctx context.Context, id uuid.UUID) (GetPostByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getPostById, id)
	var i GetPostByIdRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.AuthorID,
		&i.UpVoted,
		&i.DownVoted,
		&i.CommentsCount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AuthorName,
		&i.AuthorEmail,
		&i.AuthorImage,
	)
	return i, err
}

const getPostComments = `-- name: GetPostComments :many
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
      AND c.parent_comment_id IS NULL
  ) pc
  LEFT JOIN users u ON pc.author_id = u.id
ORDER BY
  pc.created_at DESC
`

type GetPostCommentsRow struct {
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

func (q *Queries) GetPostComments(ctx context.Context, postID uuid.UUID) ([]GetPostCommentsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostComments, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostCommentsRow
	for rows.Next() {
		var i GetPostCommentsRow
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

const increPostCommentCount = `-- name: IncrePostCommentCount :exec
UPDATE posts
SET
  comments_count = comments_count + 1
WHERE
  id = $1
`

func (q *Queries) IncrePostCommentCount(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, increPostCommentCount, id)
	return err
}