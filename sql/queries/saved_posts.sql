-- name: CreateSavedPost :one
INSERT INTO
  saved_posts ("id", "user_id", "post_id", "created_at")
VALUES
  ($1, $2, $3, $4)
RETURNING
  *
;

-- name: DeleteSavedPost :exec
DELETE FROM saved_posts
WHERE
  user_id = $1
  AND post_id = $2
;

-- name: GetAllSavedPosts :many
SELECT
  sp.*,
  u.email AS author_email,
  u.name AS author_name,
  u.image AS author_image
FROM
  (
    SELECT
      s.created_at AS saved_at,
      p.*
    FROM
      (
        SELECT
          *
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
;