package services

import (
	"context"
	"strings"
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

type UserServer struct {
	pb.UnimplementedUserServer
	DBQueries    *database.Queries
	DBConnection *pgxpool.Pool
}

func NewUserServer() *UserServer {
	return &UserServer{
		DBQueries:    cfg.DBQueries,
		DBConnection: cfg.DBConnection,
	}
}

func (server *UserServer) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	userEmail := userSession["email"].(string)
	_, err := server.DBQueries.UpdateUserRefreshToken(ctx, database.UpdateUserRefreshTokenParams{
		Email:        userEmail,
		RefreshToken: pgtype.Text{Valid: false},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return &pb.SignOutResponse{
		Message: "success",
	}, nil
}

func (server *UserServer) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	userEmail := userSession["email"].(string)
	user, err := server.DBQueries.FindUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.GetMeResponse{
		Id:           user.ID.String(),
		Email:        user.Email,
		Name:         user.Name,
		Image:        user.Image.String,
		Bio:          user.Bio.String,
		RefreshToken: user.RefreshToken.String,
		CreatedAt:    uint64(user.CreatedAt.Time.Second()),
		UpdateMedAt:  uint64(user.UpdatedAt.Time.Second()),
	}, nil
}

func (server *UserServer) UpdateMe(ctx context.Context, req *pb.UpdateMeRequest) (*pb.UpdateMeResponse, error) {
	errArr := []string{}
	if len(req.GetName()) > 0 && len(req.GetName()) < 3 {
		errArr = append(errArr, "name too short")
	}
	if len(req.GetPassword()) > 0 && len(req.GetPassword()) < 6 {
		errArr = append(errArr, "password too short")
	}
	if len(req.GetBio()) > 255 {
		errArr = append(errArr, "bio must less than 255 characters")
	}
	if len(req.GetPassword()) == 0 && len(req.GetName()) == 0 && len(req.GetImage()) == 0 && len(req.GetBio()) == 0 {
		errArr = append(errArr, "nothing has changed")
	}
	if len(errArr) > 0 {
		errMsg := strings.Join(errArr, ", ")
		return nil, status.Error(codes.InvalidArgument, errMsg)
	}
	if len(req.GetPassword()) > 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
		req.Password = string(hashedPassword)
	}
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	userEmail := userSession["email"].(string)
	user, err := server.DBQueries.FindUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	updateUserInfoParams := database.UpdateUserInfoParams{
		Email:     user.Email,
		Name:      user.Name,
		Image:     user.Image,
		Password:  user.Password,
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
	}
	if len(req.GetPassword()) > 0 {
		updateUserInfoParams.Password = req.GetPassword()
	}
	if len(req.GetName()) > 0 {
		updateUserInfoParams.Name = req.GetName()
	}
	if len(req.GetImage()) > 0 {
		updateUserInfoParams.Image = pgtype.Text{String: req.GetImage(), Valid: true}
	}
	if len(req.GetBio()) > 0 {
		updateUserInfoParams.Bio = pgtype.Text{String: req.GetBio(), Valid: true}
	}
	_, err = server.DBQueries.UpdateUserInfo(ctx, updateUserInfoParams)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	access_token, err := jwtToken.Sign(jwtToken.Claims{
		Email: user.Email,
		Name:  user.Name,
		Image: user.Image.String,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Unix(int64(userSession["iat"].(float64)), 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(int64(userSession["exp"].(float64)), 0)),
			ID:        user.ID.String(),
			Subject:   "access_token",
		},
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.UpdateMeResponse{
		AccessToken:  access_token,
		RefreshToken: user.RefreshToken.String,
	}, nil
}

func (server *UserServer) GetSession(ctx context.Context, req *pb.GetSessionRequest) (*pb.GetSessionResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	return &pb.GetSessionResponse{
		Jti:   userSession["jti"].(string),
		Sub:   userSession["sub"].(string),
		Exp:   uint64(userSession["exp"].(float64)),
		Iat:   uint64(userSession["iat"].(float64)),
		Email: userSession["email"].(string),
		Name:  userSession["name"].(string),
		Image: userSession["image"].(string),
	}, nil
}

func (server *UserServer) GetMyOverview(ctx context.Context, req *pb.GetMyOverviewRequest) (*pb.GetMyOverviewResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	userId := userSession["jti"].(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	overview, err := server.DBQueries.GetUserOverview(ctx, pgtype.UUID{Bytes: userUUID, Valid: true})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetMyOverviewResponse{
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

func (server *UserServer) GetMyPost(ctx context.Context, req *pb.GetMyPostRequest) (*pb.GetMyPostResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	userId := userSession["jti"].(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.GetSort() == "hot" {
		posts, _ := server.DBQueries.GetUserHotPosts(ctx, database.GetUserHotPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(req.GetLimit()), Offset: int32(req.GetOffset())})
		res := []*pb.PostData{}
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
		return &pb.GetMyPostResponse{
			Posts: res,
		}, nil
	}
	if req.GetSort() == "top" {
		posts, _ := server.DBQueries.GetUserTopPosts(ctx, database.GetUserTopPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(req.GetLimit()), Offset: int32(req.GetOffset())})
		res := []*pb.PostData{}
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
		return &pb.GetMyPostResponse{
			Posts: res,
		}, nil
	}
	posts, _ := server.DBQueries.GetUserNewPosts(ctx, database.GetUserNewPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(req.GetLimit()), Offset: int32(req.GetOffset())})
	res := []*pb.PostData{}
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
	return &pb.GetMyPostResponse{
		Posts: res,
	}, nil
}

func (server *UserServer) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	follower_id := userSession["jti"].(string)
	follower_uuid, _ := uuid.Parse(follower_id)
	user_id := req.GetUserId()
	if user_id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = server.DBQueries.CreateFollows(ctx, database.CreateFollowsParams{
		ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
		AuthorID:   pgtype.UUID{Bytes: user_uuid, Valid: true},
		FollowerID: pgtype.UUID{Bytes: follower_uuid, Valid: true},
		CreatedAt:  pgtype.Timestamp{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.FollowUserResponse{
		Message: "success",
	}, nil
}

func (server *UserServer) UnFollowUser(ctx context.Context, req *pb.UnFollowUserRequest) (*pb.UnFollowUserResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	follower_id := userSession["jti"].(string)
	follower_uuid, _ := uuid.Parse(follower_id)
	user_id := req.GetUserId()
	if user_id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = server.DBQueries.DeleteFollows(ctx, database.DeleteFollowsParams{
		AuthorID:   pgtype.UUID{Bytes: user_uuid, Valid: true},
		FollowerID: pgtype.UUID{Bytes: follower_uuid, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	return &pb.UnFollowUserResponse{
		Message: "success",
	}, nil
}

func (server *UserServer) CheckUserFollowed(ctx context.Context, req *pb.CheckUserFollowedRequest) (*pb.CheckUserFollowedResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	follower_id := userSession["jti"].(string)
	follower_uuid, _ := uuid.Parse(follower_id)
	user_id := req.GetUserId()
	if user_id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	follow, err := server.DBQueries.GetFollows(ctx, database.GetFollowsParams{
		AuthorID:   pgtype.UUID{Bytes: user_uuid, Valid: true},
		FollowerID: pgtype.UUID{Bytes: follower_uuid, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return &pb.CheckUserFollowedResponse{
		FollowId:  follow.ID.String(),
		CreatedAt: uint64(follow.CreatedAt.Time.Second()),
	}, nil
}

func (server *UserServer) GetFollowedUsers(ctx context.Context, req *pb.GetFollowedUsersRequest) (*pb.GetFollowedUsersResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	follower_id := userSession["jti"].(string)
	follower_uuid, _ := uuid.Parse(follower_id)
	followedAuthors, err := server.DBQueries.GetFollowedAuthor(ctx, database.GetFollowedAuthorParams{
		FollowerID: pgtype.UUID{Bytes: follower_uuid, Valid: true},
		Limit:      int32(req.GetLimit()),
		Offset:     int32(req.GetOffset()),
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	res := []*pb.UserBaseInfo{}
	for _, a := range followedAuthors {
		res = append(res, &pb.UserBaseInfo{
			Id:    a.ID.String(),
			Name:  a.Name,
			Email: a.Email,
			Image: a.Image.String,
		})
	}

	return &pb.GetFollowedUsersResponse{Followed: res}, nil
}
