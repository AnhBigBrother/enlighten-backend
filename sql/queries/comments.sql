-- name: CreateComment :exec
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
  ) pc
  LEFT JOIN users u ON pc.author_id = u.id
;