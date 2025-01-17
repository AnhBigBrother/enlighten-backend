// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkPostInteracted = `-- name: CheckPostInteracted :one
SELECT
  r.post_id,
  CASE
    WHEN l.voted IS NULL THEN 'none'
    ELSE l.voted::TEXT
  END AS voted,
  CASE
    WHEN l.saved IS NULL THEN FALSE
    ELSE TRUE
  END AS saved
FROM
  (
    SELECT
      $2 AS post_id,
      voted,
      saved.id::TEXT AS saved
    FROM
      (
        SELECT
          voted,
          post_id
        FROM
          post_votes
        WHERE
          voter_id = $1
          AND post_id = $2
      ) voted
      FULL OUTER JOIN (
        SELECT
          id,
          post_id
        FROM
          saved_posts
        WHERE
          user_id = $1
          AND post_id = $2
      ) saved ON voted.post_id = saved.post_id
  ) l
  RIGHT JOIN (
    SELECT
      $2::UUID AS post_id
  ) r ON l.post_id = r.post_id
`

type CheckPostInteractedParams struct {
	VoterID pgtype.UUID `json:"voter_id"`
	PostID  pgtype.UUID `json:"post_id"`
}

type CheckPostInteractedRow struct {
	PostID pgtype.UUID `json:"post_id"`
	Voted  string      `json:"voted"`
	Saved  bool        `json:"saved"`
}

func (q *Queries) CheckPostInteracted(ctx context.Context, arg CheckPostInteractedParams) (CheckPostInteractedRow, error) {
	row := q.db.QueryRow(ctx, checkPostInteracted, arg.VoterID, arg.PostID)
	var i CheckPostInteractedRow
	err := row.Scan(&i.PostID, &i.Voted, &i.Saved)
	return i, err
}

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
	ID        pgtype.UUID      `json:"id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	AuthorID  pgtype.UUID      `json:"author_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRow(ctx, createPost,
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

const getHotPosts = `-- name: GetHotPosts :many
SELECT
  p.id, p.title, p.content, p.author_id, p.up_voted, p.down_voted, p.comments_count, p.created_at, p.updated_at,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image,
  p.up_voted + p.down_voted + p.comments_count AS total_interactions
FROM
  posts p
  LEFT JOIN users u ON p.author_id = u.id
ORDER BY
  total_interactions DESC
LIMIT
  $1
OFFSET
  $2
`

type GetHotPostsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetHotPostsRow struct {
	ID                pgtype.UUID      `json:"id"`
	Title             string           `json:"title"`
	Content           string           `json:"content"`
	AuthorID          pgtype.UUID      `json:"author_id"`
	UpVoted           int32            `json:"up_voted"`
	DownVoted         int32            `json:"down_voted"`
	CommentsCount     int32            `json:"comments_count"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	UpdatedAt         pgtype.Timestamp `json:"updated_at"`
	AuthorName        pgtype.Text      `json:"author_name"`
	AuthorEmail       pgtype.Text      `json:"author_email"`
	AuthorImage       pgtype.Text      `json:"author_image"`
	TotalInteractions int32            `json:"total_interactions"`
}

func (q *Queries) GetHotPosts(ctx context.Context, arg GetHotPostsParams) ([]GetHotPostsRow, error) {
	rows, err := q.db.Query(ctx, getHotPosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetHotPostsRow{}
	for rows.Next() {
		var i GetHotPostsRow
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
			&i.TotalInteractions,
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

const getNewPosts = `-- name: GetNewPosts :many
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
LIMIT
  $1
OFFSET
  $2
`

type GetNewPostsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetNewPostsRow struct {
	ID            pgtype.UUID      `json:"id"`
	Title         string           `json:"title"`
	Content       string           `json:"content"`
	AuthorID      pgtype.UUID      `json:"author_id"`
	UpVoted       int32            `json:"up_voted"`
	DownVoted     int32            `json:"down_voted"`
	CommentsCount int32            `json:"comments_count"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
	AuthorName    pgtype.Text      `json:"author_name"`
	AuthorEmail   pgtype.Text      `json:"author_email"`
	AuthorImage   pgtype.Text      `json:"author_image"`
}

func (q *Queries) GetNewPosts(ctx context.Context, arg GetNewPostsParams) ([]GetNewPostsRow, error) {
	rows, err := q.db.Query(ctx, getNewPosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetNewPostsRow{}
	for rows.Next() {
		var i GetNewPostsRow
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

func (q *Queries) GetPostByAuthor(ctx context.Context, authorID pgtype.UUID) ([]Post, error) {
	rows, err := q.db.Query(ctx, getPostByAuthor, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
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
	ID            pgtype.UUID      `json:"id"`
	Title         string           `json:"title"`
	Content       string           `json:"content"`
	AuthorID      pgtype.UUID      `json:"author_id"`
	UpVoted       int32            `json:"up_voted"`
	DownVoted     int32            `json:"down_voted"`
	CommentsCount int32            `json:"comments_count"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
	AuthorName    pgtype.Text      `json:"author_name"`
	AuthorEmail   pgtype.Text      `json:"author_email"`
	AuthorImage   pgtype.Text      `json:"author_image"`
}

func (q *Queries) GetPostById(ctx context.Context, id pgtype.UUID) (GetPostByIdRow, error) {
	row := q.db.QueryRow(ctx, getPostById, id)
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

const getTopPosts = `-- name: GetTopPosts :many
SELECT
  p.id, p.title, p.content, p.author_id, p.up_voted, p.down_voted, p.comments_count, p.created_at, p.updated_at,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  posts p
  LEFT JOIN users u ON p.author_id = u.id
ORDER BY
  p.up_voted DESC,
  p.down_voted ASC
LIMIT
  $1
OFFSET
  $2
`

type GetTopPostsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetTopPostsRow struct {
	ID            pgtype.UUID      `json:"id"`
	Title         string           `json:"title"`
	Content       string           `json:"content"`
	AuthorID      pgtype.UUID      `json:"author_id"`
	UpVoted       int32            `json:"up_voted"`
	DownVoted     int32            `json:"down_voted"`
	CommentsCount int32            `json:"comments_count"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
	AuthorName    pgtype.Text      `json:"author_name"`
	AuthorEmail   pgtype.Text      `json:"author_email"`
	AuthorImage   pgtype.Text      `json:"author_image"`
}

func (q *Queries) GetTopPosts(ctx context.Context, arg GetTopPostsParams) ([]GetTopPostsRow, error) {
	rows, err := q.db.Query(ctx, getTopPosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTopPostsRow{}
	for rows.Next() {
		var i GetTopPostsRow
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

func (q *Queries) IncrePostCommentCount(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, increPostCommentCount, id)
	return err
}
