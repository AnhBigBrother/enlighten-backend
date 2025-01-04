// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
  users (
    "id",
    "email",
    "name",
    "password",
    "image",
    "refresh_token",
    "created_at",
    "updated_at"
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
  id, email, name, password, image, refresh_token, created_at, updated_at, bio
`

type CreateUserParams struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Password     string
	Image        sql.NullString
	RefreshToken sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Name,
		arg.Password,
		arg.Image,
		arg.RefreshToken,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Image,
		&i.RefreshToken,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Bio,
	)
	return i, err
}

const deleteUserInfo = `-- name: DeleteUserInfo :exec
DELETE FROM users
WHERE
  email = $1
`

func (q *Queries) DeleteUserInfo(ctx context.Context, email string) error {
	_, err := q.db.ExecContext(ctx, deleteUserInfo, email)
	return err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT
  id, email, name, password, image, refresh_token, created_at, updated_at, bio
FROM
  users
WHERE
  email = $1
LIMIT
  1
`

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Image,
		&i.RefreshToken,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Bio,
	)
	return i, err
}

const getTopAuthor = `-- name: GetTopAuthor :many
SELECT
  a.author_id, a.total_posts, a.total_upvoted,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  (
    SELECT
      author_id,
      COUNT(id) AS total_posts,
      SUM(up_voted) AS total_upvoted
    FROM
      posts
    GROUP BY
      author_id
  ) a
  INNER JOIN users u ON author_id = u.id
ORDER BY
  total_upvoted DESC,
  total_posts DESC
LIMIT
  $1
OFFSET
  $2
`

type GetTopAuthorParams struct {
	Limit  int32
	Offset int32
}

type GetTopAuthorRow struct {
	AuthorID     uuid.UUID
	TotalPosts   int64
	TotalUpvoted int64
	AuthorName   string
	AuthorEmail  string
	AuthorImage  sql.NullString
}

func (q *Queries) GetTopAuthor(ctx context.Context, arg GetTopAuthorParams) ([]GetTopAuthorRow, error) {
	rows, err := q.db.QueryContext(ctx, getTopAuthor, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopAuthorRow
	for rows.Next() {
		var i GetTopAuthorRow
		if err := rows.Scan(
			&i.AuthorID,
			&i.TotalPosts,
			&i.TotalUpvoted,
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

const getUserHotPosts = `-- name: GetUserHotPosts :many
SELECT
  p.id, p.title, p.content, p.author_id, p.up_voted, p.down_voted, p.comments_count, p.created_at, p.updated_at,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image,
  p.up_voted + p.down_voted + p.comments_count AS total_interactions
FROM
  posts p
  INNER JOIN (
    SELECT
      id, email, name, password, image, refresh_token, created_at, updated_at, bio
    FROM
      users
    WHERE
      users.id = $1
  ) AS u ON p.author_id = u.id
ORDER BY
  total_interactions DESC
LIMIT
  $2
OFFSET
  $3
`

type GetUserHotPostsParams struct {
	ID     uuid.UUID
	Limit  int32
	Offset int32
}

type GetUserHotPostsRow struct {
	ID                uuid.UUID
	Title             string
	Content           string
	AuthorID          uuid.UUID
	UpVoted           int32
	DownVoted         int32
	CommentsCount     int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
	AuthorName        string
	AuthorEmail       string
	AuthorImage       sql.NullString
	TotalInteractions int32
}

func (q *Queries) GetUserHotPosts(ctx context.Context, arg GetUserHotPostsParams) ([]GetUserHotPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserHotPosts, arg.ID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserHotPostsRow
	for rows.Next() {
		var i GetUserHotPostsRow
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
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserNewPosts = `-- name: GetUserNewPosts :many
SELECT
  p.id, p.title, p.content, p.author_id, p.up_voted, p.down_voted, p.comments_count, p.created_at, p.updated_at,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  posts p
  INNER JOIN (
    SELECT
      id, email, name, password, image, refresh_token, created_at, updated_at, bio
    FROM
      users
    WHERE
      users.id = $1
  ) AS u ON p.author_id = u.id
ORDER BY
  p.created_at DESC
LIMIT
  $2
OFFSET
  $3
`

type GetUserNewPostsParams struct {
	ID     uuid.UUID
	Limit  int32
	Offset int32
}

type GetUserNewPostsRow struct {
	ID            uuid.UUID
	Title         string
	Content       string
	AuthorID      uuid.UUID
	UpVoted       int32
	DownVoted     int32
	CommentsCount int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AuthorName    string
	AuthorEmail   string
	AuthorImage   sql.NullString
}

func (q *Queries) GetUserNewPosts(ctx context.Context, arg GetUserNewPostsParams) ([]GetUserNewPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserNewPosts, arg.ID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserNewPostsRow
	for rows.Next() {
		var i GetUserNewPostsRow
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

const getUserOverview = `-- name: GetUserOverview :one
SELECT
  u.id,
  u.name,
  u.email,
  u.image,
  u.bio,
  u.created_at,
  u.updated_at,
  a.total_posts,
  a.total_upvoted,
  a.total_downvoted
FROM
  users u
  INNER JOIN (
    SELECT
      author_id,
      COUNT(*) AS total_posts,
      SUM(up_voted) AS total_upvoted,
      SUM(down_voted) AS total_downvoted
    FROM
      posts
    WHERE
      author_id = $1
    GROUP BY
      author_id
  ) a ON u.id = a.author_id
`

type GetUserOverviewRow struct {
	ID             uuid.UUID
	Name           string
	Email          string
	Image          sql.NullString
	Bio            sql.NullString
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TotalPosts     int64
	TotalUpvoted   int64
	TotalDownvoted int64
}

func (q *Queries) GetUserOverview(ctx context.Context, authorID uuid.UUID) (GetUserOverviewRow, error) {
	row := q.db.QueryRowContext(ctx, getUserOverview, authorID)
	var i GetUserOverviewRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Image,
		&i.Bio,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TotalPosts,
		&i.TotalUpvoted,
		&i.TotalDownvoted,
	)
	return i, err
}

const getUserTopPosts = `-- name: GetUserTopPosts :many
SELECT
  p.id, p.title, p.content, p.author_id, p.up_voted, p.down_voted, p.comments_count, p.created_at, p.updated_at,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  posts p
  INNER JOIN (
    SELECT
      id, email, name, password, image, refresh_token, created_at, updated_at, bio
    FROM
      users
    WHERE
      users.id = $1
  ) AS u ON p.author_id = u.id
ORDER BY
  p.up_voted DESC,
  p.down_voted ASC
LIMIT
  $2
OFFSET
  $3
`

type GetUserTopPostsParams struct {
	ID     uuid.UUID
	Limit  int32
	Offset int32
}

type GetUserTopPostsRow struct {
	ID            uuid.UUID
	Title         string
	Content       string
	AuthorID      uuid.UUID
	UpVoted       int32
	DownVoted     int32
	CommentsCount int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AuthorName    string
	AuthorEmail   string
	AuthorImage   sql.NullString
}

func (q *Queries) GetUserTopPosts(ctx context.Context, arg GetUserTopPostsParams) ([]GetUserTopPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserTopPosts, arg.ID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserTopPostsRow
	for rows.Next() {
		var i GetUserTopPostsRow
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

const updateUserInfo = `-- name: UpdateUserInfo :one
UPDATE users
SET
  "name" = $2,
  "password" = $3,
  "image" = $4,
  "bio" = $5,
  "updated_at" = $6
WHERE
  email = $1
RETURNING
  id, email, name, password, image, refresh_token, created_at, updated_at, bio
`

type UpdateUserInfoParams struct {
	Email     string
	Name      string
	Password  string
	Image     sql.NullString
	Bio       sql.NullString
	UpdatedAt time.Time
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserInfo,
		arg.Email,
		arg.Name,
		arg.Password,
		arg.Image,
		arg.Bio,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Image,
		&i.RefreshToken,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Bio,
	)
	return i, err
}

const updateUserRefreshToken = `-- name: UpdateUserRefreshToken :one
UPDATE users
SET
  refresh_token = $2
WHERE
  email = $1
RETURNING
  id, email, name, password, image, refresh_token, created_at, updated_at, bio
`

type UpdateUserRefreshTokenParams struct {
	Email        string
	RefreshToken sql.NullString
}

func (q *Queries) UpdateUserRefreshToken(ctx context.Context, arg UpdateUserRefreshTokenParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserRefreshToken, arg.Email, arg.RefreshToken)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Image,
		&i.RefreshToken,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Bio,
	)
	return i, err
}
