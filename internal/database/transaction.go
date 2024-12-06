package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func (q *Queries) VoteComment(ctx context.Context, db *sql.DB, voter_id, comment_id uuid.UUID, vote string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := q.WithTx(tx)
	vc, err := qtx.FindCommentVote(ctx, FindCommentVoteParams{
		VoterID:   voter_id,
		CommentID: comment_id,
	})

	if err != nil {
		err = qtx.CreateVoteComment(ctx, CreateVoteCommentParams{
			VoterID:   voter_id,
			CommentID: comment_id,
			ID:        uuid.New(),
			Voted:     Voted(vote),
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		if vote == "up" {
			err = qtx.IncreCommentUpVoted(ctx, comment_id)
			if err != nil {
				return err
			}
		} else {
			err = qtx.IncreCommentDownVoted(ctx, comment_id)
			if err != nil {
				return err
			}
		}

		return tx.Commit()
	}

	if vc.Voted == Voted(vote) {
		err = qtx.DeleteVoteComment(ctx, vc.ID)
		if err != nil {
			return err
		}
		if vote == "up" {
			err = qtx.DecreCommentUpVoted(ctx, comment_id)
			if err != nil {
				return err
			}
		} else {
			err = qtx.DecreCommentDownVoted(ctx, comment_id)
			if err != nil {
				return err
			}
		}
		return tx.Commit()
	}

	err = qtx.ChangeVoteComment(ctx, vc.ID)
	if err != nil {
		return err
	}

	if vote == "up" {
		err = qtx.IncreCommentUpVoted(ctx, comment_id)
		if err != nil {
			return err
		}
		err = qtx.DecreCommentDownVoted(ctx, comment_id)
		if err != nil {
			return err
		}
	} else {
		err = qtx.IncreCommentDownVoted(ctx, comment_id)
		if err != nil {
			return err
		}
		err = qtx.DecreCommentUpVoted(ctx, comment_id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (q *Queries) VotePost(ctx context.Context, db *sql.DB, voter_id, post_id uuid.UUID, vote string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := q.WithTx(tx)
	vp, err := qtx.FindPostVote(ctx, FindPostVoteParams{
		VoterID: voter_id,
		PostID:  post_id,
	})

	if err != nil {
		err = qtx.CreateVotePost(ctx, CreateVotePostParams{
			VoterID:   voter_id,
			PostID:    post_id,
			ID:        uuid.New(),
			Voted:     Voted(vote),
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		if vote == "up" {
			err = qtx.IncrePostUpVoted(ctx, post_id)
			if err != nil {
				return err
			}
		} else {
			err = qtx.IncrePostDownVoted(ctx, post_id)
			if err != nil {
				return err
			}
		}

		return tx.Commit()
	}

	if vp.Voted == Voted(vote) {
		err = qtx.DeleteVotePost(ctx, vp.ID)
		if err != nil {
			return err
		}
		if vote == "up" {
			err = qtx.DecrePostUpVoted(ctx, post_id)
			if err != nil {
				return err
			}
		} else {
			err = qtx.DecrePostDownVoted(ctx, post_id)
			if err != nil {
				return err
			}
		}
		return tx.Commit()
	}

	err = qtx.ChangeVotePost(ctx, vp.ID)
	if err != nil {
		return err
	}

	if vote == "up" {
		err = qtx.IncrePostUpVoted(ctx, post_id)
		if err != nil {
			return err
		}
		err = qtx.DecrePostDownVoted(ctx, post_id)
		if err != nil {
			return err
		}
	} else {
		err = qtx.IncrePostDownVoted(ctx, post_id)
		if err != nil {
			return err
		}
		err = qtx.DecrePostUpVoted(ctx, post_id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (q *Queries) AddComment(ctx context.Context, db *sql.DB, comment string, author_id, post_id uuid.UUID, parent_comment_id uuid.NullUUID) (Comment, error) {
	tx, err := db.Begin()
	if err != nil {
		return Comment{}, err
	}
	defer tx.Rollback()

	qtx := q.WithTx(tx)
	comm, err := qtx.CreateComment(ctx, CreateCommentParams{
		ID:              uuid.New(),
		Comment:         comment,
		AuthorID:        author_id,
		PostID:          post_id,
		ParentCommentID: parent_comment_id,
		CreatedAt:       time.Now(),
	})
	if err != nil {
		return Comment{}, err
	}
	err = qtx.IncrePostCommentCount(ctx, post_id)
	if err != nil {
		return Comment{}, err
	}

	return comm, tx.Commit()
}
