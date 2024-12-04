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
  "image" = $4
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