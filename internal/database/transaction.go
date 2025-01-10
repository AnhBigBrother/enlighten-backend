package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (q *Queries) VoteComment(ctx context.Context, conn *pgxpool.Pool, voter_id, comment_id pgtype.UUID, vote string) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := q.WithTx(tx)
	vc, err := qtx.FindCommentVote(ctx, FindCommentVoteParams{
		VoterID:   voter_id,
		CommentID: comment_id,
	})

	if err != nil {
		err = qtx.CreateVoteComment(ctx, CreateVoteCommentParams{
			VoterID:   voter_id,
			CommentID: comment_id,
			ID: pgtype.UUID{
				Bytes: uuid.New(),
				Valid: true,
			},
			Voted: Voted(vote),
			CreatedAt: pgtype.Timestamp{
				Time:             time.Now(),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
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

		return tx.Commit(ctx)
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
		return tx.Commit(ctx)
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

	return tx.Commit(ctx)
}

func (q *Queries) VotePost(ctx context.Context, conn *pgxpool.Pool, voter_id, post_id pgtype.UUID, vote string) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := q.WithTx(tx)
	vp, err := qtx.FindPostVote(ctx, FindPostVoteParams{
		VoterID: voter_id,
		PostID:  post_id,
	})

	if err != nil {
		err = qtx.CreateVotePost(ctx, CreateVotePostParams{
			VoterID: voter_id,
			PostID:  post_id,
			ID: pgtype.UUID{
				Bytes: uuid.New(),
				Valid: true,
			},
			Voted: Voted(vote),
			CreatedAt: pgtype.Timestamp{
				Time:             time.Now(),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
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

		return tx.Commit(ctx)
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
		return tx.Commit(ctx)
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

	return tx.Commit(ctx)
}

func (q *Queries) AddComment(ctx context.Context, conn *pgxpool.Pool, comment string, author_id, post_id, parent_comment_id pgtype.UUID) (Comment, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return Comment{}, err
	}
	defer tx.Rollback(ctx)

	qtx := q.WithTx(tx)
	comm, err := qtx.CreateComment(ctx, CreateCommentParams{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Comment:         comment,
		AuthorID:        author_id,
		PostID:          post_id,
		ParentCommentID: parent_comment_id,
		CreatedAt: pgtype.Timestamp{
			Time:             time.Now(),
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
	})
	if err != nil {
		return Comment{}, err
	}
	err = qtx.IncrePostCommentCount(ctx, post_id)
	if err != nil {
		return Comment{}, err
	}

	return comm, tx.Commit(ctx)
}
