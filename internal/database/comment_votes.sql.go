// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: comment_votes.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const changeVoteComment = `-- name: ChangeVoteComment :exec
UPDATE comment_votes
SET
  voted = CASE
    WHEN voted = 'up'::VOTED THEN 'down'::VOTED
    ELSE 'up'::VOTED
  END
WHERE
  id = $1
`

func (q *Queries) ChangeVoteComment(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, changeVoteComment, id)
	return err
}

const createVoteComment = `-- name: CreateVoteComment :exec
INSERT INTO
  comment_votes (
    "voter_id",
    "comment_id",
    "id",
    "voted",
    "created_at"
  )
VALUES
  ($1, $2, $3, $4, $5)
`

type CreateVoteCommentParams struct {
	VoterID   pgtype.UUID      `json:"voter_id"`
	CommentID pgtype.UUID      `json:"comment_id"`
	ID        pgtype.UUID      `json:"id"`
	Voted     Voted            `json:"voted"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) CreateVoteComment(ctx context.Context, arg CreateVoteCommentParams) error {
	_, err := q.db.Exec(ctx, createVoteComment,
		arg.VoterID,
		arg.CommentID,
		arg.ID,
		arg.Voted,
		arg.CreatedAt,
	)
	return err
}

const decreCommentDownVoted = `-- name: DecreCommentDownVoted :exec
UPDATE COMMENTS
SET
  down_voted = down_voted - 1
WHERE
  id = $1
  AND down_voted > 0
`

func (q *Queries) DecreCommentDownVoted(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, decreCommentDownVoted, id)
	return err
}

const decreCommentUpVoted = `-- name: DecreCommentUpVoted :exec
UPDATE COMMENTS
SET
  up_voted = up_voted - 1
WHERE
  id = $1
  AND up_voted > 0
`

func (q *Queries) DecreCommentUpVoted(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, decreCommentUpVoted, id)
	return err
}

const deleteVoteComment = `-- name: DeleteVoteComment :exec
DELETE FROM comment_votes
WHERE
  id = $1
`

func (q *Queries) DeleteVoteComment(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteVoteComment, id)
	return err
}

const findCommentVote = `-- name: FindCommentVote :one
SELECT
  id, voted, voter_id, comment_id, created_at
FROM
  comment_votes
WHERE
  voter_id = $1
  AND comment_id = $2
`

type FindCommentVoteParams struct {
	VoterID   pgtype.UUID `json:"voter_id"`
	CommentID pgtype.UUID `json:"comment_id"`
}

func (q *Queries) FindCommentVote(ctx context.Context, arg FindCommentVoteParams) (CommentVote, error) {
	row := q.db.QueryRow(ctx, findCommentVote, arg.VoterID, arg.CommentID)
	var i CommentVote
	err := row.Scan(
		&i.ID,
		&i.Voted,
		&i.VoterID,
		&i.CommentID,
		&i.CreatedAt,
	)
	return i, err
}

const getCommentVotes = `-- name: GetCommentVotes :one
SELECT
  id, voted, voter_id, comment_id, created_at
FROM
  comment_votes
WHERE
  comment_id = $1
  AND voter_id = $2
LIMIT
  1
`

type GetCommentVotesParams struct {
	CommentID pgtype.UUID `json:"comment_id"`
	VoterID   pgtype.UUID `json:"voter_id"`
}

func (q *Queries) GetCommentVotes(ctx context.Context, arg GetCommentVotesParams) (CommentVote, error) {
	row := q.db.QueryRow(ctx, getCommentVotes, arg.CommentID, arg.VoterID)
	var i CommentVote
	err := row.Scan(
		&i.ID,
		&i.Voted,
		&i.VoterID,
		&i.CommentID,
		&i.CreatedAt,
	)
	return i, err
}

const increCommentDownVoted = `-- name: IncreCommentDownVoted :exec
UPDATE COMMENTS
SET
  down_voted = down_voted + 1
WHERE
  id = $1
`

func (q *Queries) IncreCommentDownVoted(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, increCommentDownVoted, id)
	return err
}

const increCommentUpVoted = `-- name: IncreCommentUpVoted :exec
UPDATE COMMENTS
SET
  up_voted = up_voted + 1
WHERE
  id = $1
`

func (q *Queries) IncreCommentUpVoted(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, increCommentUpVoted, id)
	return err
}
