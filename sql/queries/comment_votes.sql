-- name: FindCommentVote :one
SELECT
  *
FROM
  comment_votes
WHERE
  voter_id = $1
  AND comment_id = $2
;

-- name: CreateVoteComment :exec
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
;

-- name: ChangeVoteComment :exec
UPDATE comment_votes
SET
  voted = CASE
    WHEN voted = 'up'::VOTED THEN 'down'::VOTED
    ELSE 'up'::VOTED
  END
WHERE
  id = $1
;

-- name: DeleteVoteComment :exec
DELETE FROM comment_votes
WHERE
  id = $1
;

-- name: IncreCommentUpVoted :exec
UPDATE COMMENTS
SET
  up_voted = up_voted + 1
WHERE
  id = $1
;

-- name: IncreCommentDownVoted :exec
UPDATE COMMENTS
SET
  down_voted = down_voted + 1
WHERE
  id = $1
;

-- name: DecreCommentUpVoted :exec
UPDATE COMMENTS
SET
  up_voted = up_voted - 1
WHERE
  id = $1
  AND up_voted > 0
;

-- name: DecreCommentDownVoted :exec
UPDATE COMMENTS
SET
  down_voted = down_voted - 1
WHERE
  id = $1
  AND down_voted > 0
;

-- name: GetCommentVotes :one
SELECT
  *
FROM
  comment_votes
WHERE
  comment_id = $1
  AND voter_id = $2
LIMIT
  1
;