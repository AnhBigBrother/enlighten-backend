-- name: CreateFollows :exec
INSERT INTO
  user_follows ("id", "follower_id", "author_id", "created_at")
VALUES
  ($1, $2, $3, $4)
;

-- name: DeleteFollows :exec
DELETE FROM user_follows
WHERE
  author_id = $1
  AND follower_id = $2
;

-- name: GetFollows :one
SELECT
  *
FROM
  user_follows
WHERE
  follower_id = $1
  AND author_id = $2
;

-- name: GetFollowedPosts :many
SELECT
  fa.author_email,
  fa.author_name,
  fa.author_image,
  p.*
FROM
  (
    SELECT
      u.id,
      u.name AS author_name,
      u.email AS author_email,
      u.image AS author_image
    FROM
      (
        SELECT
          *
        FROM
          user_follows
        WHERE
          user_follows.follower_id = $1
      ) uf
      LEFT JOIN users u ON uf.author_id = u.id
  ) fa
  LEFT JOIN posts p ON fa.id = p.author_id
ORDER BY
  p.created_at DESC
LIMIT
  $2
OFFSET
  $3
;

-- name: GetNewFollowedPosts :many
SELECT
  ap.*,
  CASE
    WHEN uf.id IS NOT NULL THEN 'Followed'
    ELSE 'Recommend'
  END AS status
FROM
  (
    SELECT
      p.*,
      u.email AS author_email,
      u.name AS author_name,
      u.image AS author_image
    FROM
      posts p
      LEFT JOIN users u ON p.author_id = u.id
  ) ap
  LEFT JOIN (
    SELECT
      *
    FROM
      user_follows
    WHERE
      user_follows.follower_id = $1
  ) uf ON ap.author_id = uf.author_id
ORDER BY
  status ASC,
  ap.created_at DESC
LIMIT
  $2
OFFSET
  $3
;

-- name: GetTopFollowedPosts :many
SELECT
  ap.*,
  CASE
    WHEN uf.id IS NOT NULL THEN 'Followed'
    ELSE 'Recommend'
  END AS status
FROM
  (
    SELECT
      p.*,
      u.email AS author_email,
      u.name AS author_name,
      u.image AS author_image
    FROM
      posts p
      LEFT JOIN users u ON p.author_id = u.id
  ) ap
  LEFT JOIN (
    SELECT
      *
    FROM
      user_follows
    WHERE
      user_follows.follower_id = $1
  ) uf ON ap.author_id = uf.author_id
ORDER BY
  status ASC,
  ap.up_voted DESC
LIMIT
  $2
OFFSET
  $3
;

-- name: GetHotFollowedPosts :many
SELECT
  ap.*,
  ap.up_voted + ap.down_voted + ap.comments_count AS total_interactions,
  CASE
    WHEN uf.id IS NOT NULL THEN 'Followed'
    ELSE 'Recommend'
  END AS status
FROM
  (
    SELECT
      p.*,
      u.email AS author_email,
      u.name AS author_name,
      u.image AS author_image
    FROM
      posts p
      LEFT JOIN users u ON p.author_id = u.id
  ) ap
  LEFT JOIN (
    SELECT
      *
    FROM
      user_follows
    WHERE
      user_follows.follower_id = $1
  ) uf ON ap.author_id = uf.author_id
ORDER BY
  status ASC,
  total_interactions DESC
LIMIT
  $2
OFFSET
  $3
;

-- name: GetFollowedAuthor :many
SELECT
  u.id,
  u.email,
  u.name,
  u.image
FROM
  (
    SELECT
      *
    FROM
      user_follows uf
    WHERE
      uf.follower_id = $1
  ) fu
  INNER JOIN users u ON fu.author_id = u.id
ORDER BY
  fu.created_at
LIMIT
  $2
OFFSET
  $3
;