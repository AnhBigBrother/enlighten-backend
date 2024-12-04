-- name: FindPostVote :one
SELECT
  *
FROM
  post_votes
WHERE
  voter_id = $1
  AND post_id = $2
;

-- name: CreateVotePost :exec
INSERT INTO
  post_votes (
    "voter_id",
    "post_id",
    "id",
    "voted",
    "created_at"
  )
VALUES
  ($1, $2, $3, $4, $5)
;

-- name: ChangeVotePost :exec
UPDATE post_votes
SET
  voted = CASE
    WHEN voted = 'up'::VOTED THEN 'down'::VOTED
    ELSE 'up'::VOTED
  END
WHERE
  id = $1
;

-- name: DeleteVotePost :exec
DELETE FROM post_votes
WHERE
  id = $1
;

-- name: IncrePostUpVoted :exec
UPDATE posts
SET
  up_voted = up_voted + 1
WHERE
  id = $1
;

-- name: IncrePostDownVoted :exec
UPDATE posts
SET
  down_voted = down_voted + 1
WHERE
  id = $1
;

-- name: DecrePostUpVoted :exec
UPDATE posts
SET
  up_voted = up_voted - 1
WHERE
  id = $1
  AND up_voted > 0
;

-- name: DecrePostDownVoted :exec
UPDATE posts
SET
  down_voted = down_voted - 1
WHERE
  id = $1
  AND down_voted > 0
;