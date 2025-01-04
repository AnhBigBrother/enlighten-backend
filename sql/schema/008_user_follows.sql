-- +goose Up
CREATE TABLE
  user_follows (
    "id" UUID NOT NULL UNIQUE,
    "follower_id" UUID NOT NULL REFERENCES users (id),
    "author_id" UUID NOT NULL REFERENCES users (id),
    "created_at" TIMESTAMP NOT NULL,
    UNIQUE (follower_id, author_id)
  )
;

-- +goose Down
DROP TABLE user_follows
;