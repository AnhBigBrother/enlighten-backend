package services

import (
	"context"
	"strings"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/pb"
	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct {
	DB *database.Queries
}

func NewUserServer() *UserServer {
	return &UserServer{
		DB: cfg.DBQueries,
	}
}

func (server *UserServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	if req.GetEmail() == "" || req.GetName() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing some of parameters: Email, Name, Password")
	}
	userId := uuid.New()
	currentTime := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)

	refresh_token, err := token.Sign(token.Claims{
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
	_, err = server.DB.CreateUser(ctx, createUserParams)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	access_token, err := token.Sign(token.Claims{
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

	return &pb.SignupResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

func (server *UserServer) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing some of parameters: Email, Name, Password")
	}
	user, err := server.DB.FindUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword())); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	currentTime := time.Now()
	access_token, err := token.Sign(token.Claims{
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
	refresh_token, err := token.Sign(token.Claims{
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

	_, err = server.DB.UpdateUserRefreshToken(ctx, database.UpdateUserRefreshTokenParams{
		Email:        req.GetEmail(),
		RefreshToken: pgtype.Text{String: refresh_token, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.SigninResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

func (server *UserServer) Signout(ctx context.Context, req *pb.SignoutRequest) (*pb.SignoutResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	userEmail := userSession["email"].(string)
	_, err := server.DB.UpdateUserRefreshToken(ctx, database.UpdateUserRefreshTokenParams{
		Email:        userEmail,
		RefreshToken: pgtype.Text{Valid: false},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return &pb.SignoutResponse{
		Mesasge: "success",
	}, nil
}

func (server *UserServer) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	userEmail := userSession["email"].(string)
	user, err := server.DB.FindUserByEmail(ctx, userEmail)
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
		CreatedAt:    timestamppb.New(user.CreatedAt.Time),
		UpdateMedAt:  timestamppb.New(user.UpdatedAt.Time),
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
	user, err := server.DB.FindUserByEmail(ctx, userEmail)
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
	_, err = server.DB.UpdateUserInfo(ctx, updateUserInfoParams)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}
	access_token, err := token.Sign(token.Claims{
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
		Exp:   timestamppb.New(userSession["exp"].(*jwt.NumericDate).Time),
		Iat:   timestamppb.New(userSession["iat"].(*jwt.NumericDate).Time),
		Email: userSession["email"].(string),
		Name:  userSession["name"].(string),
		Image: userSession["image"].(string),
	}, nil
}

func (server *UserServer) GetAccessToken(ctx context.Context, req *pb.GetAccessTokenRequest) (*pb.GetAccessTokenResponse, error) {
	refresh_token := req.GetRefreshToken()
	if refresh_token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing refresh_token")
	}
	claims, err := token.Parse(refresh_token)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	currentTime := time.Now()
	if int64(claims["exp"].(float64)) < currentTime.Unix() {
		return nil, status.Error(codes.Unavailable, "refresh_token expired")
	}
	user, err := server.DB.FindUserByEmail(ctx, claims["email"].(string))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	access_token, err := token.Sign(token.Claims{
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
		RefreshToken: user.RefreshToken.String,
	}, nil
}

func (server *UserServer) GetMyOverview(ctx context.Context, req *pb.GetMyOverviewRequest) (*pb.GetMyOverviewResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	userId := userSession["jti"].(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	overview, err := server.DB.GetUserOverview(ctx, pgtype.UUID{Bytes: userUUID, Valid: true})
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
		CreatedAt:     timestamppb.New(overview.CreatedAt.Time),
		UpdateMedAt:   timestamppb.New(overview.UpdatedAt.Time),
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
		posts, _ := server.DB.GetUserHotPosts(ctx, database.GetUserHotPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(req.GetLimit()), Offset: int32(req.GetOffset())})
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
				CreatedAt: timestamppb.New(p.CreatedAt.Time),
				UpdatedAt: timestamppb.New(p.UpdatedAt.Time),
			})
		}
		return &pb.GetMyPostResponse{
			Posts: res,
		}, nil
	}
	if req.GetSort() == "top" {
		posts, _ := server.DB.GetUserTopPosts(ctx, database.GetUserTopPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(req.GetLimit()), Offset: int32(req.GetOffset())})
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
				CreatedAt: timestamppb.New(p.CreatedAt.Time),
				UpdatedAt: timestamppb.New(p.UpdatedAt.Time),
			})
		}
		return &pb.GetMyPostResponse{
			Posts: res,
		}, nil
	}
	posts, _ := server.DB.GetUserNewPosts(ctx, database.GetUserNewPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(req.GetLimit()), Offset: int32(req.GetOffset())})
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
			CreatedAt: timestamppb.New(p.CreatedAt.Time),
			UpdatedAt: timestamppb.New(p.UpdatedAt.Time),
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
	err = server.DB.CreateFollows(ctx, database.CreateFollowsParams{
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

	err = server.DB.DeleteFollows(ctx, database.DeleteFollowsParams{
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
	follow, err := server.DB.GetFollows(ctx, database.GetFollowsParams{
		AuthorID:   pgtype.UUID{Bytes: user_uuid, Valid: true},
		FollowerID: pgtype.UUID{Bytes: follower_uuid, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return &pb.CheckUserFollowedResponse{
		FollowId:  follow.ID.String(),
		CreatedAt: timestamppb.New(follow.CreatedAt.Time),
	}, nil
}

func (server *UserServer) GetFollowedUser(ctx context.Context, req *pb.GetFollowedUsersRequest) (*pb.GetFollowedUsersResponse, error) {
	userSession := ctx.Value(cfg.CtxKeys.User).(map[string]interface{})
	follower_id := userSession["jti"].(string)
	follower_uuid, _ := uuid.Parse(follower_id)
	followedAuthors, err := server.DB.GetFollowedAuthor(ctx, database.GetFollowedAuthorParams{
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
