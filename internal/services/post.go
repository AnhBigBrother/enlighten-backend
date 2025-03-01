package services

import (
	"context"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/pb"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostServer struct {
	pb.UnimplementedPostServer
	DBQueries    *database.Queries
	DBConnection *pgxpool.Pool
}

func NewPostServer() *PostServer {
	return &PostServer{
		DBQueries:    cfg.DBQueries,
		DBConnection: cfg.DBConnection,
	}
}

func (server *PostServer) GetFollowedPosts(ctx context.Context, req *pb.GetFollowedPostsRequest) (*pb.GetFollowedPostsResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_uuid, err := uuid.Parse(userSession["jti"].(string))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ans := []*pb.PostData{}
	if req.GetSort() == "hot" {
		posts, err := server.DBQueries.GetHotFollowedPosts(ctx, database.GetHotFollowedPostsParams{
			FollowerID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			Limit:      int32(req.GetLimit()),
			Offset:     int32(req.GetOffset()),
		})
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		for _, p := range posts {
			ans = append(ans, &pb.PostData{
				Id:      p.ID.String(),
				Title:   p.Title,
				Content: p.Content,
				Author: &pb.UserBaseInfo{
					Id:    p.AuthorID.String(),
					Name:  p.AuthorName.String,
					Email: p.AuthorEmail.String,
					Image: p.AuthorImage.String,
				},
				Upvote:    uint32(p.UpVoted),
				Downvote:  uint32(p.DownVoted),
				Comments:  uint32(p.CommentsCount),
				CreatedAt: uint64(p.CreatedAt.Time.Second()),
				UpdatedAt: uint64(p.UpdatedAt.Time.Second()),
			})
		}
	} else if req.GetSort() == "top" {
		posts, err := server.DBQueries.GetTopFollowedPosts(ctx, database.GetTopFollowedPostsParams{
			FollowerID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			Limit:      int32(req.GetLimit()),
			Offset:     int32(req.GetOffset()),
		})
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		for _, p := range posts {
			ans = append(ans, &pb.PostData{
				Id:      p.ID.String(),
				Title:   p.Title,
				Content: p.Content,
				Author: &pb.UserBaseInfo{
					Id:    p.AuthorID.String(),
					Name:  p.AuthorName.String,
					Email: p.AuthorEmail.String,
					Image: p.AuthorImage.String,
				},
				Upvote:    uint32(p.UpVoted),
				Downvote:  uint32(p.DownVoted),
				Comments:  uint32(p.CommentsCount),
				CreatedAt: uint64(p.CreatedAt.Time.Second()),
				UpdatedAt: uint64(p.UpdatedAt.Time.Second()),
			})
		}
	} else {
		posts, err := server.DBQueries.GetNewFollowedPosts(ctx, database.GetNewFollowedPostsParams{
			FollowerID: pgtype.UUID{Bytes: user_uuid, Valid: true},
			Limit:      int32(req.GetLimit()),
			Offset:     int32(req.GetOffset()),
		})
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		for _, p := range posts {
			ans = append(ans, &pb.PostData{
				Id:      p.ID.String(),
				Title:   p.Title,
				Content: p.Content,
				Author: &pb.UserBaseInfo{
					Id:    p.AuthorID.String(),
					Name:  p.AuthorName.String,
					Email: p.AuthorEmail.String,
					Image: p.AuthorImage.String,
				},
				Upvote:    uint32(p.UpVoted),
				Downvote:  uint32(p.DownVoted),
				Comments:  uint32(p.CommentsCount),
				CreatedAt: uint64(p.CreatedAt.Time.Second()),
				UpdatedAt: uint64(p.UpdatedAt.Time.Second()),
			})
		}
	}

	return &pb.GetFollowedPostsResponse{
		Posts: ans,
	}, nil
}

func (server *PostServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	currentTime := time.Now()
	post, err := server.DBQueries.CreatePost(ctx, database.CreatePostParams{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Title:     req.GetTitle(),
		Content:   req.GetContent(),
		AuthorID:  pgtype.UUID{Bytes: user_uuid, Valid: true},
		CreatedAt: pgtype.Timestamp{Time: currentTime, InfinityModifier: pgtype.Finite, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: currentTime, InfinityModifier: pgtype.Finite, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.CreatePostResponse{
		Created: &pb.CreatePostResponse_CreatedPost{
			Id:        post.ID.String(),
			Title:     post.Title,
			Content:   post.Content,
			AuthorId:  post.AuthorID.String(),
			CreatedAt: uint64(currentTime.Second()),
			UpdatedAt: uint64(currentTime.Second()),
		},
	}, nil
}

func (server *PostServer) UpVotePost(ctx context.Context, req *pb.UpVotePostRequest) (*pb.UpVotePostResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post_uuid, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = server.DBQueries.VotePost(
		ctx,
		server.DBConnection,
		pgtype.UUID{Bytes: user_uuid, Valid: true},
		pgtype.UUID{Bytes: post_uuid, Valid: true},
		"up",
	)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.UpVotePostResponse{Message: "success"}, nil
}

func (server *PostServer) DownVotePost(ctx context.Context, req *pb.DownVotePostRequest) (*pb.DownVotePostResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post_uuid, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = server.DBQueries.VotePost(
		ctx,
		server.DBConnection,
		pgtype.UUID{Bytes: user_uuid, Valid: true},
		pgtype.UUID{Bytes: post_uuid, Valid: true},
		"down",
	)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.DownVotePostResponse{Message: "success"}, nil
}

func (server *PostServer) SavePost(ctx context.Context, req *pb.SavePostRequest) (*pb.SavePostResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post_uuid, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	savedPost, err := server.DBQueries.CreateSavedPost(
		ctx,
		database.CreateSavedPostParams{
			ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
			UserID:    pgtype.UUID{Bytes: user_uuid, Valid: true},
			PostID:    pgtype.UUID{Bytes: post_uuid, Valid: true},
			CreatedAt: pgtype.Timestamp{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
		},
	)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.SavePostResponse{
		SavedPostId: savedPost.ID.String(),
		CreatedAt:   uint64(savedPost.CreatedAt.Time.Second()),
	}, nil
}

func (server *PostServer) UnSavePost(ctx context.Context, req *pb.UnSavePostRequest) (*pb.UnSavePostResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post_uuid, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = server.DBQueries.DeleteSavedPost(ctx, database.DeleteSavedPostParams{
		UserID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		PostID: pgtype.UUID{Bytes: post_uuid, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.UnSavePostResponse{
		Message: "success",
	}, nil
}

func (server *PostServer) GetAllSavedPosts(ctx context.Context, req *pb.GetAllSavedPostsRequest) (*pb.GetAllSavedPostsResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	savedPosts, err := server.DBQueries.GetAllSavedPosts(ctx, database.GetAllSavedPostsParams{
		UserID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		Limit:  int32(req.GetLimit()),
		Offset: int32(req.GetOffset()),
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	ans := []*pb.PostData{}
	for _, p := range savedPosts {
		ans = append(ans, &pb.PostData{
			Id:      p.ID.String(),
			Title:   p.Title.String,
			Content: p.Content.String,
			Author: &pb.UserBaseInfo{
				Id:    p.AuthorID.String(),
				Name:  p.AuthorName.String,
				Email: p.AuthorEmail.String,
				Image: p.AuthorImage.String,
			},
			Upvote:    uint32(p.UpVoted.Int32),
			Downvote:  uint32(p.DownVoted.Int32),
			Comments:  uint32(p.CommentsCount.Int32),
			CreatedAt: uint64(p.CreatedAt.Time.Second()),
			UpdatedAt: uint64(p.UpdatedAt.Time.Second()),
		})
	}
	return &pb.GetAllSavedPostsResponse{
		Posts: ans,
	}, nil
}

func (server *PostServer) CheckPostInteracted(ctx context.Context, req *pb.CheckPostInteractedRequest) (*pb.CheckPostInteractedResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post_uuid, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	checked, err := server.DBQueries.CheckPostInteracted(ctx, database.CheckPostInteractedParams{
		VoterID: pgtype.UUID{Bytes: user_uuid, Valid: true},
		PostID:  pgtype.UUID{Bytes: post_uuid, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	res := &pb.CheckPostInteractedResponse{
		PostId: checked.PostID.String(),
		Saved:  checked.Saved,
	}
	if checked.Voted == "up" {
		res.Voted = pb.Voted_UP
	} else {
		res.Voted = pb.Voted_DOWN
	}

	return res, nil
}

func (server *PostServer) AddPostComment(ctx context.Context, req *pb.AddPostCommentRequest) (*pb.AddPostCommentResponse, error) {
	userSession, _ := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	user_id, _ := userSession["jti"].(string)
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post_uuid, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	comment, err := server.DBQueries.AddComment(
		ctx,
		server.DBConnection,
		req.GetComment(),
		pgtype.UUID{Bytes: user_uuid, Valid: true},
		pgtype.UUID{Bytes: post_uuid, Valid: true},
		pgtype.UUID{Valid: false},
	)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.AddPostCommentResponse{
		Created: &pb.AddPostCommentResponse_CreatedComment{
			Id:              comment.ID.String(),
			Comment:         comment.Comment,
			AuthorId:        comment.AuthorID.String(),
			PostId:          comment.PostID.String(),
			ParentCommentId: comment.ParentCommentID.String(),
			CreatedAt:       uint64(comment.CreatedAt.Time.Second()),
		},
	}, nil
}
