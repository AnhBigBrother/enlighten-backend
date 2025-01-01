-- +goose Up
ALTER TABLE users
ADD bio VARCHAR(255)
;

-- +goose Down
ALTER TABLE users
DROP COLUMN bio
;