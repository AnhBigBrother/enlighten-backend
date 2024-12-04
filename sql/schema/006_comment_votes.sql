-- +goose Up
CREATE TABLE
  comment_votes (
    "id" UUID NOT NULL PRIMARY KEY,
    "voted" VOTED NOT NULL,
    "voter_id" UUID NOT NULL REFERENCES users (id),
    "comment_id" UUID NOT NULL REFERENCES COMMENTS (id),
    "created_at" TIMESTAMP NOT NULL,
    UNIQUE (voter_id, comment_id)
  )
;

-- +goose Down
DROP TABLE comment_votes
;