-- name: CreateComment :one
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
RETURNING
  *
;

-- name: GetPostComments :many
SELECT
  pc.*,
  u.email AS author_email,
  u.name AS author_name,
  u.image AS author_image
FROM
  (
    SELECT
      *
    FROM
      COMMENTS c
    WHERE
      c.post_id = $1
      AND c.parent_comment_id IS NULL
    ORDER BY
      c.created_at DESC
    LIMIT
      $2
    OFFSET
      $3
  ) pc
  LEFT JOIN users u ON pc.author_id = u.id
ORDER BY
  pc.created_at DESC
;

-- name: GetCommentsReplies :many
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
      AND c.parent_comment_id = $2
    ORDER BY
      c.created_at DESC
    LIMIT
      $3
    OFFSET
      $4
  ) pc
  LEFT JOIN users u ON pc.author_id = u.id
;