-- name: CreateUser :one
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
  *
;

-- name: FindUserByEmail :one
SELECT
  *
FROM
  users
WHERE
  email = $1
LIMIT
  1
;

-- name: UpdateUserRefreshToken :one
UPDATE users
SET
  refresh_token = $2
WHERE
  email = $1
RETURNING
  *
;

-- name: UpdateUserInfo :one
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
  *
;

-- name: DeleteUserInfo :exec
DELETE FROM users
WHERE
  email = $1
;

-- name: GetUserOverview :one
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
;

-- name: GetUserNewPosts :many
SELECT
  p.*,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  posts p
  INNER JOIN (
    SELECT
      *
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
;

-- name: GetUserTopPosts :many
SELECT
  p.*,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  posts p
  INNER JOIN (
    SELECT
      *
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
;

-- name: GetUserHotPosts :many
SELECT
  p.*,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image,
  p.up_voted + p.down_voted + p.comments_count AS total_interactions
FROM
  posts p
  INNER JOIN (
    SELECT
      *
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
;