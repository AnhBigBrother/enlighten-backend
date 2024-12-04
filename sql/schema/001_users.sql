-- +goose Up
CREATE TABLE
  users (
    "id" UUID NOT NULL UNIQUE,
    "email" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "image" TEXT,
    "refresh_token" TEXT,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
  )
;

-- +goose Down
DROP TABLE users
;