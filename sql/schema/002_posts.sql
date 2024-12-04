-- +goose Up
CREATE TABLE
  posts (
    "id" UUID NOT NULL PRIMARY KEY,
    "title" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "author_id" UUID NOT NULL REFERENCES users (id),
    "up_voted" INT NOT NULL DEFAULT (0),
    "down_voted" INT NOT NULL DEFAULT (0),
    "comments_count" INT NOT NULL DEFAULT (0),
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
  )
;

-- +goose Down
DROP TABLE posts
;