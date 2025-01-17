-- +goose Up
CREATE TABLE
  saved_posts (
    "id" UUID NOT NULL UNIQUE,
    "user_id" UUID NOT NULL REFERENCES users (id),
    "post_id" UUID NOT NULL REFERENCES posts (id),
    "created_at" TIMESTAMP NOT NULL,
    UNIQUE (user_id, post_id)
  )
;

-- +goose Down
DROP TABLE saved_posts
;