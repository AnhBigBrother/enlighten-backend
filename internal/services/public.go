package services

import (
	"context"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/pb"
	jwtToken "github.com/AnhBigBrother/enlighten-backend/internal/pkg/jwt-token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PublicServer struct {
	DBQueries    *database.Queries
	DBConnection *pgxpool.Pool
	pb.UnimplementedPublicServer
}

func NewPublicServer() *PublicServer {
	return &PublicServer{
		DBQueries:    cfg.DBQueries,
		DBConnection: cfg.DBConnection,
	}
}

func (server *PublicServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if req.GetEmail() == "" || req.GetName() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing some of parameters: Email, Name, Password")
	}
	userId := uuid.New()
	currentTime := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)

	refresh_token, err := jwtToken.Sign(jwtToken.Claims{
		Email: req.GetEmail(),
		Name:  req.GetName(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.RefreshTokenAge) * time.Second)),
			ID:        userId.String(),
			Subject:   "refresh_token",
		},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	createUserParams := database.CreateUserParams{
		ID:           pgtype.UUID{Bytes: userId, Valid: true},
		Email:        req.GetEmail(),
		Name:         req.GetName(),
		Password:     string(hashedPassword),
		RefreshToken: pgtype.Text{String: refresh_token, Valid: true},
		CreatedAt:    pgtype.Timestamp{Time: currentTime, InfinityModifier: pgtype.Finite, Valid: true},
		UpdatedAt:    pgtype.Timestamp{Time: currentTime, InfinityModifier: pgtype.Finite, Valid: true},
	}
	if req.GetImage() != "" {
		createUserParams.Image = pgtype.Text{String: req.GetImage(), Valid: true}
	}
	_, err = server.DBQueries.CreateUser(ctx, createUserParams)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	access_token, err := jwtToken.Sign(jwtToken.Claims{
		Email: req.GetEmail(),
		Name:  req.GetName(),
		Image: req.GetImage(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.AccessTokenAge) * time.Second)),
			ID:        userId.String(),
			Subject:   "access_token",
		},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return &pb.SignUpResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

func (server *PublicServer) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing some of parameters: Email, Name, Password")
	}
	user, err := server.DBQueries.FindUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword())); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	currentTime := time.Now()
	access_token, err := jwtToken.Sign(jwtToken.Claims{
		Email: req.GetEmail(),
		Name:  user.Name,
		Image: user.Image.String,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.AccessTokenAge) * time.Second)),
			ID:        user.ID.String(),
			Subject:   "access_token",
		},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	refresh_token, err := jwtToken.Sign(jwtToken.Claims{
		Email: req.GetEmail(),
		Name:  user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.RefreshTokenAge) * time.Second)),
			ID:        user.ID.String(),
			Subject:   "refresh_token",
		},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	_, err = server.DBQueries.UpdateUserRefreshToken(ctx, database.UpdateUserRefreshTokenParams{
		Email:        req.GetEmail(),
		RefreshToken: pgtype.Text{String: refresh_token, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.SignInResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

func (server *PublicServer) GetAccessToken(ctx context.Context, req *pb.GetAccessTokenRequest) (*pb.GetAccessTokenResponse, error) {
	refresh_token := req.GetRefreshToken()
	if refresh_token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing refresh_token")
	}
	claims, err := jwtToken.ParseAndValidate(refresh_token)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	user, err := server.DBQueries.FindUserByEmail(ctx, claims["email"].(string))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	currentTime := time.Now()
	access_token, err := jwtToken.Sign(jwtToken.Claims{
		Email: user.Email,
		Name:  user.Name,
		Image: user.Image.String,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.AccessTokenAge) * time.Second)),
			ID:        claims["jti"].(string),
			Subject:   "access_token",
		},
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.GetAccessTokenResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

func (server *PublicServer) GetUserOverview(ctx context.Context, req *pb.GetUserOverviewRequest) (*pb.GetUserOverviewResponse, error) {
	user_id := req.GetUserId()
	if user_id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	overview, err := server.DBQueries.GetUserOverview(ctx, pgtype.UUID{Bytes: user_uuid, Valid: true})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.GetUserOverviewResponse{
		Id:            overview.ID.String(),
		Name:          overview.Name,
		Email:         overview.Email,
		Image:         overview.Image.String,
		Bio:           overview.Bio.String,
		TotalPosts:    uint32(overview.TotalPosts),
		TotalUpvote:   uint32(overview.TotalUpvoted),
		TotalDownvote: uint32(overview.TotalDownvoted),
		Follower:      uint32(overview.Follower),
		Following:     uint32(overview.Following),
		CreatedAt:     uint64(overview.CreatedAt.Time.Second()),
		UpdatedAt:     uint64(overview.UpdatedAt.Time.Second()),
	}, nil
}

func (server *PublicServer) GetUserPosts(ctx context.Context, req *pb.GetUserPostsRequest) (*pb.GetUserPostsResponse, error) {
	user_id := req.GetUserId()
	if user_id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res := []*pb.PostData{}
	if req.GetSort() == "hot" {
		posts, _ := server.DBQueries.GetUserHotPosts(
			ctx,
			database.GetUserHotPostsParams{
				ID:     pgtype.UUID{Bytes: user_uuid, Valid: true},
				Limit:  int32(req.GetLimit()),
				Offset: int32(req.GetOffset()),
			},
		)
		for _, p := range posts {
			res = append(res, &pb.PostData{
				Id:      p.ID.String(),
				Title:   p.Title,
				Content: p.Content,
				Author: &pb.UserBaseInfo{
					Id:    p.AuthorID.String(),
					Name:  p.AuthorName,
					Email: p.AuthorEmail,
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
		posts, _ := server.DBQueries.GetUserTopPosts(
			ctx,
			database.GetUserTopPostsParams{
				ID:     pgtype.UUID{Bytes: user_uuid, Valid: true},
				Limit:  int32(req.GetLimit()),
				Offset: int32(req.GetOffset()),
			},
		)
		for _, p := range posts {
			res = append(res, &pb.PostData{
				Id:      p.ID.String(),
				Title:   p.Title,
				Content: p.Content,
				Author: &pb.UserBaseInfo{
					Id:    p.AuthorID.String(),
					Name:  p.AuthorName,
					Email: p.AuthorEmail,
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
		posts, _ := server.DBQueries.GetUserNewPosts(
			ctx,
			database.GetUserNewPostsParams{
				ID:     pgtype.UUID{Bytes: user_uuid, Valid: true},
				Limit:  int32(req.GetLimit()),
				Offset: int32(req.GetOffset()),
			},
		)
		for _, p := range posts {
			res = append(res, &pb.PostData{
				Id:      p.ID.String(),
				Title:   p.Title,
				Content: p.Content,
				Author: &pb.UserBaseInfo{
					Id:    p.AuthorID.String(),
					Name:  p.AuthorName,
					Email: p.AuthorEmail,
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

	return &pb.GetUserPostsResponse{
		Posts: res,
	}, nil
}

func (server *PublicServer) GetAllPosts(ctx context.Context, req *pb.GetAllPostsRequest) (*pb.GetAllPostsResponse, error) {
	res := []*pb.PostData{}
	if req.GetSort() == "hot" {
		posts, _ := server.DBQueries.GetHotPosts(ctx, database.GetHotPostsParams{
			Limit:  int32(req.GetLimit()),
			Offset: int32(req.GetOffset()),
		})
		for _, p := range posts {
			res = append(res, &pb.PostData{
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
		posts, _ := server.DBQueries.GetTopPosts(ctx, database.GetTopPostsParams{
			Limit:  int32(req.GetLimit()),
			Offset: int32(req.GetOffset()),
		})
		for _, p := range posts {
			res = append(res, &pb.PostData{
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
		posts, _ := server.DBQueries.GetNewPosts(ctx, database.GetNewPostsParams{
			Limit:  int32(req.GetLimit()),
			Offset: int32(req.GetOffset()),
		})
		for _, p := range posts {
			res = append(res, &pb.PostData{
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
	return &pb.GetAllPostsResponse{Posts: res}, nil
}

func (server *PublicServer) GetPostById(ctx context.Context, req *pb.GetPostByIdRequest) (*pb.GetPostByIdResponse, error) {
	post_id := req.GetPostId()
	post_uuid, err := uuid.Parse(post_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	post, err := server.DBQueries.GetPostById(ctx, pgtype.UUID{Bytes: post_uuid, Valid: true})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.GetPostByIdResponse{
		Post: &pb.PostData{
			Id:      post.ID.String(),
			Title:   post.Title,
			Content: post.Content,
			Author: &pb.UserBaseInfo{
				Id:    post.AuthorID.String(),
				Name:  post.AuthorName.String,
				Email: post.AuthorEmail.String,
				Image: post.AuthorImage.String,
			},
			Upvote:    uint32(post.UpVoted),
			Downvote:  uint32(post.DownVoted),
			Comments:  uint32(post.CommentsCount),
			CreatedAt: uint64(post.CreatedAt.Time.Second()),
			UpdatedAt: uint64(post.UpdatedAt.Time.Second()),
		},
	}, nil
}

func (server *PublicServer) GetPostComments(ctx context.Context, req *pb.GetPostCommentsRequest) (*pb.GetPostCommentsResponse, error) {
	post_uuid, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	comments, err := server.DBQueries.GetPostComments(ctx, database.GetPostCommentsParams{
		PostID: pgtype.UUID{Bytes: post_uuid, Valid: true},
		Limit:  int32(req.GetLimit()),
		Offset: int32(req.GetOffset()),
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	resComments := []*pb.CommentData{}
	for _, c := range comments {
		resComments = append(resComments, &pb.CommentData{
			Id:      c.ID.String(),
			Comment: c.Comment,
			Author: &pb.UserBaseInfo{
				Id:    c.AuthorID.String(),
				Name:  c.AuthorName.String,
				Email: c.AuthorEmail.String,
				Image: c.AuthorImage.String,
			},
			PostId:          c.PostID.String(),
			ParentCommentId: c.ParentCommentID.String(),
			Upvote:          uint32(c.UpVoted),
			Downvote:        uint32(c.DownVoted),
			CreatedAt:       uint64(c.CreatedAt.Time.Second()),
		})
	}
	return &pb.GetPostCommentsResponse{
		Comments: resComments,
	}, nil
}

func (server *PublicServer) GetCommentReplies(ctx context.Context, req *pb.GetCommentRepliesRequest) (*pb.GetCommentRepliesResponse, error) {
	post_uuid, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	parent_comment_uuid, err := uuid.Parse(req.GetCommentId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	replies, err := server.DBQueries.GetCommentsReplies(ctx, database.GetCommentsRepliesParams{
		PostID:          pgtype.UUID{Bytes: post_uuid, Valid: true},
		ParentCommentID: pgtype.UUID{Bytes: parent_comment_uuid, Valid: true},
		Limit:           int32(req.GetLimit()),
		Offset:          int32(req.GetOffset()),
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	resReplies := []*pb.CommentData{}
	for _, r := range replies {
		resReplies = append(resReplies, &pb.CommentData{
			Id:      r.ID.String(),
			Comment: r.Comment,
			Author: &pb.UserBaseInfo{
				Id:    r.AuthorID.String(),
				Name:  r.AuthorName.String,
				Email: r.AuthorEmail.String,
				Image: r.AuthorImage.String,
			},
			PostId:          r.PostID.String(),
			ParentCommentId: r.ParentCommentID.String(),
			Upvote:          uint32(r.UpVoted),
			Downvote:        uint32(r.DownVoted),
			CreatedAt:       uint64(r.CreatedAt.Time.Second()),
		})
	}
	return &pb.GetCommentRepliesResponse{
		Replies: resReplies,
	}, nil
}

func (server *PublicServer) GetTopUsers(ctx context.Context, req *pb.GetTopUsersRequest) (*pb.GetTopUsersResponse, error) {
	users, err := server.DBQueries.GetTopAuthor(ctx, database.GetTopAuthorParams{
		Limit:  int32(req.GetLimit()),
		Offset: int32(req.GetOffset()),
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	res := []*pb.UserBaseInfo{}
	for _, u := range users {
		res = append(res, &pb.UserBaseInfo{
			Id:    u.ID.String(),
			Name:  u.Name,
			Email: u.Email,
			Image: u.Image.String,
		})
	}
	return &pb.GetTopUsersResponse{
		Users: res,
	}, nil
}
