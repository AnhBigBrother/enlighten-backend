-- +goose Up
CREATE TABLE
  COMMENTS (
    "id" UUID NOT NULL PRIMARY KEY,
    "comment" TEXT NOT NULL,
    "author_id" UUID NOT NULL REFERENCES users (id),
    "post_id" UUID NOT NULL REFERENCES posts (id),
    "parent_comment_id" UUID REFERENCES COMMENTS (id),
    "up_voted" INT NOT NULL DEFAULT (0),
    "down_voted" INT NOT NULL DEFAULT (0),
    "created_at" TIMESTAMP NOT NULL,
    UNIQUE (author_id, post_id, parent_comment_id)
  )
;

-- +goose Down
DROP TABLE COMMENTS
;