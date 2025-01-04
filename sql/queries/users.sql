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
WITH
  total_interactions AS (
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
      (
        SELECT
          author_id,
          COUNT(*) AS total_posts,
          SUM(up_voted) AS total_upvoted,
          SUM(down_voted) AS total_downvoted
        FROM
          posts p
        WHERE
          p.author_id = $1
        GROUP BY
          p.author_id
      ) a
      INNER JOIN users u ON u.id = a.author_id
  ),
  total_follows AS (
    SELECT
      *
    FROM
      (
        SELECT
          author_id,
          COUNT(id) AS follower
        FROM
          user_follows f1
        WHERE
          f1.author_id = $1
        GROUP BY
          author_id
      ) fr
      JOIN (
        SELECT
          follower_id,
          COUNT(id) AS "following"
        FROM
          user_follows f2
        WHERE
          f2.follower_id = $1
        GROUP BY
          f2.follower_id
      ) fg ON fr.author_id = fg.follower_id
  )
SELECT
  total_interactions.*,
  total_follows.follower,
  total_follows.following
FROM
  total_interactions
  LEFT JOIN total_follows ON total_interactions.id = total_follows.author_id
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

-- name: GetTopAuthor :many
SELECT
  a.*,
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
;