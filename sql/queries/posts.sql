-- name: CreatePost :one
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
  *
;

-- name: GetPostById :one
SELECT
  p.*,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  (
    SELECT
      *
    FROM
      posts p1
    WHERE
      p1.id = $1
    LIMIT
      1
  ) p
  LEFT JOIN users u ON p.author_id = u.id
;

-- name: GetPostByAuthor :many
SELECT
  *
FROM
  posts
WHERE
  author_id = $1
;

-- name: GetAllPosts :many
SELECT
  p.*,
  u.name AS author_name,
  u.email AS author_email,
  u.image AS author_image
FROM
  posts p
  LEFT JOIN users u ON p.author_id = u.id
ORDER BY
  p.created_at DESC
;

-- name: GetPostComments :many
SELECT
  pc.*,
  u.email AS user_email,
  u.name AS user_name,
  u.image AS user_image
FROM
  (
    SELECT
      *
    FROM
      COMMENTS c
    WHERE
      c.post_id = $1
      AND c.parent_comment_id IS NULL
  ) pc
  LEFT JOIN users u ON pc.author_id = u.id
ORDER BY
  pc.created_at DESC
;

-- name: IncrePostCommentCount :exec
UPDATE posts
SET
  comments_count = comments_count + 1
WHERE
  id = $1
;