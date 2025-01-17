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

-- name: GetNewPosts :many
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
LIMIT
  $1
OFFSET
  $2
;

-- name: GetTopPosts :many
SELECT
  p.*,
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
;

-- name: GetHotPosts :many
SELECT
  p.*,
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
;

-- name: IncrePostCommentCount :exec
UPDATE posts
SET
  comments_count = comments_count + 1
WHERE
  id = $1
;

-- name: CheckPostInteracted :one
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
;