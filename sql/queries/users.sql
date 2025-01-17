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
      CASE
        WHEN a.total_posts IS NULL THEN 0
        ELSE a.total_posts
      END::INTEGER AS total_posts,
      CASE
        WHEN a.total_upvoted IS NULL THEN 0
        ELSE a.total_upvoted
      END::INTEGER AS total_upvoted,
      CASE
        WHEN a.total_downvoted IS NULL THEN 0
        ELSE a.total_downvoted
      END::INTEGER AS total_downvoted
    FROM
      (
        SELECT
          *
        FROM
          users u1
        WHERE
          u1.id = $1
      ) u
      LEFT JOIN (
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
      ) a ON u.id = a.author_id
  ),
  total_follows AS (
    SELECT
      *,
      CASE
        WHEN fr.author_id IS NULL THEN fg.follower_id
        ELSE fr.author_id
      END AS user_id
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
      FULL OUTER JOIN (
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
  CASE
    WHEN total_follows.follower IS NULL THEN 0
    ELSE total_follows.follower
  END::INTEGER AS follower,
  CASE
    WHEN total_follows.following IS NULL THEN 0
    ELSE total_follows.following
  END::INTEGER AS "following"
FROM
  total_interactions
  LEFT JOIN total_follows ON total_interactions.id = total_follows.user_id
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
WITH
  count_followers AS (
    SELECT
      uf.author_id,
      COUNT(id) AS followers
    FROM
      user_follows uf
    GROUP BY
      uf.author_id
  ),
  authors AS (
    SELECT
      p.author_id,
      COUNT(id) AS total_posts,
      SUM(up_voted) AS total_upvoted
    FROM
      posts p
    GROUP BY
      p.author_id
  )
SELECT
  u.id,
  u.name,
  u.email,
  u.image,
  CASE
    WHEN a.total_posts IS NULL THEN 0
    ELSE a.total_posts
  END::INTEGER AS total_posts,
  CASE
    WHEN a.total_upvoted IS NULL THEN 0
    ELSE a.total_upvoted
  END::INTEGER AS total_upvoted,
  CASE
    WHEN cf.followers IS NULL THEN 0
    ELSE cf.followers
  END::INTEGER AS total_follower
FROM
  users u
  LEFT JOIN authors a ON u.id = a.author_id
  LEFT JOIN count_followers cf ON u.id = cf.author_id
ORDER BY
  total_upvoted DESC,
  total_follower DESC,
  total_posts DESC
LIMIT
  $1
OFFSET
  $2
;