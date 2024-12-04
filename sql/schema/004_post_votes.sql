-- +goose Up
CREATE TABLE
  post_votes (
    "id" UUID NOT NULL PRIMARY KEY,
    "voted" VOTED NOT NULL,
    "voter_id" UUID NOT NULL REFERENCES users (id),
    "post_id" UUID NOT NULL REFERENCES posts (id),
    "created_at" TIMESTAMP NOT NULL,
    UNIQUE (voter_id, post_id)
  )
;

-- +goose Down
DROP TABLE post_votes
;