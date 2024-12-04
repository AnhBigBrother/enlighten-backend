-- +goose Up
CREATE TYPE VOTED AS ENUM('up', 'down')
;

-- +goose Down
DROP TYPE VOTED
;