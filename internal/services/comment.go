package services

import (
	"context"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/pb"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentServer struct {
	pb.UnimplementedCommentServer
	DBQueries    *database.Queries
	DBConnection *pgxpool.Pool
}

func NewCommentServer() *CommentServer {
	return &CommentServer{
		DBQueries:    cfg.DBQueries,
		DBConnection: cfg.DBConnection,
	}
}

func (server *CommentServer) UpVoteComment(ctx context.Context, req *pb.UpVoteCommentRequest) (*pb.UpVoteCommentResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	comment_id := req.GetCommentId()
	comment_uuid, err := uuid.Parse(comment_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = server.DBQueries.VoteComment(
		ctx,
		server.DBConnection,
		pgtype.UUID{Bytes: user_uuid, Valid: true},
		pgtype.UUID{Bytes: comment_uuid, Valid: true},
		"up",
	)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.UpVoteCommentResponse{
		Message: "success",
	}, nil
}

func (server *CommentServer) DownVoteComment(ctx context.Context, req *pb.DownVoteCommentRequest) (*pb.DownVoteCommentResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	comment_id := req.GetCommentId()
	comment_uuid, err := uuid.Parse(comment_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = server.DBQueries.VoteComment(
		ctx,
		server.DBConnection,
		pgtype.UUID{Bytes: user_uuid, Valid: true},
		pgtype.UUID{Bytes: comment_uuid, Valid: true},
		"down",
	)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.DownVoteCommentResponse{
		Message: "success",
	}, nil
}

func (server *CommentServer) CheckCommentInteracted(ctx context.Context, req *pb.CheckCommentInteractedRequest) (*pb.CheckCommentInteractedResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	comment_id := req.GetCommentId()
	comment_uuid, err := uuid.Parse(comment_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	check, err := server.DBQueries.GetCommentVotes(
		ctx,
		database.GetCommentVotesParams{
			CommentID: pgtype.UUID{Bytes: comment_uuid, Valid: true},
			VoterID:   pgtype.UUID{Bytes: user_uuid, Valid: true},
		},
	)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	ans := &pb.CheckCommentInteractedResponse{}
	if check.Voted == "up" {
		ans.Voted = pb.Voted_UP
	} else {
		ans.Voted = pb.Voted_DOWN
	}

	return ans, nil
}

func (server *CommentServer) ReplyComment(ctx context.Context, req *pb.ReplyCommentRequest) (*pb.ReplyCommentResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	comment_uuid, err := uuid.Parse(req.GetCommentId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post_uuid, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	reply, err := server.DBQueries.AddComment(
		ctx,
		server.DBConnection,
		req.GetReplyBody(),
		pgtype.UUID{Bytes: user_uuid, Valid: true},
		pgtype.UUID{Bytes: post_uuid, Valid: true},
		pgtype.UUID{Bytes: comment_uuid, Valid: true},
	)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.ReplyCommentResponse{
		Created: &pb.ReplyCommentResponse_CreatedReply{
			Id:              reply.ID.String(),
			Reply:           reply.Comment,
			AuthorId:        reply.AuthorID.String(),
			PostId:          reply.PostID.String(),
			ParentCommentId: reply.ParentCommentID.String(),
			CreatedAt:       uint64(reply.CreatedAt.Time.Second()),
		},
	}, nil
}
