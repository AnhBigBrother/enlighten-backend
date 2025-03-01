package services

import (
	"context"
	"errors"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/pb"
	jwtToken "github.com/AnhBigBrother/enlighten-backend/internal/pkg/jwt-token"
	oauthprovider "github.com/AnhBigBrother/enlighten-backend/internal/pkg/oauth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OauthServer struct {
	pb.UnimplementedOauthServer
	DBQueries    *database.Queries
	DBConnection *pgxpool.Pool
}

func NewOauthServer() *OauthServer {
	return &OauthServer{
		DBQueries:    cfg.DBQueries,
		DBConnection: cfg.DBConnection,
	}
}

func (server *OauthServer) OauthGoogle(ctx context.Context, req *pb.OauthGoogleRequest) (*pb.OauthGoogleResponse, error) {
	userData, err := oauthprovider.GetGoogleUserInfo(req.GetTokenType(), req.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	access_token, refresh_token, err := server.signInOauthUser(userData)
	if err != nil {
		return &pb.OauthGoogleResponse{
			AccessToken:  "",
			RefreshToken: "",
			User: &pb.OauthUserData{
				Email: userData.Email,
				Name:  userData.Name,
				Image: userData.Picture,
			},
		}, nil
	}
	return &pb.OauthGoogleResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User: &pb.OauthUserData{
			Email: userData.Email,
			Name:  userData.Name,
			Image: userData.Picture,
		},
	}, nil
}

func (server *OauthServer) OauthGithub(ctx context.Context, req *pb.OauthGithubRequest) (*pb.OauthGithubResponse, error) {
	userData, err := oauthprovider.GetGithubUserInfo(req.GetTokenType(), req.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	access_token, refresh_token, err := server.signInOauthUser(userData)
	if err != nil {
		return &pb.OauthGithubResponse{
			AccessToken:  "",
			RefreshToken: "",
			User: &pb.OauthUserData{
				Email: userData.Email,
				Name:  userData.Name,
				Image: userData.Picture,
			},
		}, nil
	}
	return &pb.OauthGithubResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User: &pb.OauthUserData{
			Email: userData.Email,
			Name:  userData.Name,
			Image: userData.Picture,
		},
	}, nil
}

func (server *OauthServer) OauthMicrosoft(ctx context.Context, req *pb.OauthMicrosoftRequest) (*pb.OauthMicrosoftResponse, error) {
	userData, err := oauthprovider.GetMicrosoftUserInfo(req.GetTokenType(), req.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	access_token, refresh_token, err := server.signInOauthUser(userData)
	if err != nil {
		return &pb.OauthMicrosoftResponse{
			AccessToken:  "",
			RefreshToken: "",
			User: &pb.OauthUserData{
				Email: userData.Email,
				Name:  userData.Name,
				Image: userData.Picture,
			},
		}, nil
	}
	return &pb.OauthMicrosoftResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User: &pb.OauthUserData{
			Email: userData.Email,
			Name:  userData.Name,
			Image: userData.Picture,
		},
	}, nil
}

func (server *OauthServer) OauthDiscord(ctx context.Context, req *pb.OauthDiscordRequest) (*pb.OauthDiscordResponse, error) {
	userData, err := oauthprovider.GetDiscordUserInfo(req.GetTokenType(), req.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	access_token, refresh_token, err := server.signInOauthUser(userData)
	if err != nil {
		return &pb.OauthDiscordResponse{
			AccessToken:  "",
			RefreshToken: "",
			User: &pb.OauthUserData{
				Email: userData.Email,
				Name:  userData.Name,
				Image: userData.Picture,
			},
		}, nil
	}
	return &pb.OauthDiscordResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User: &pb.OauthUserData{
			Email: userData.Email,
			Name:  userData.Name,
			Image: userData.Picture,
		},
	}, nil
}

func (server *OauthServer) signInOauthUser(user oauthprovider.UserInfo) (accessToken string, refreshToken string, err error) {
	dbUser, err := server.DBQueries.FindUserByEmail(context.Background(), user.Email)

	if err != nil {
		return "", "", errors.New("unregistered user")
	}

	currentTime := time.Now()
	refresh_token, err := jwtToken.Sign(jwtToken.Claims{
		Email: dbUser.Email,
		Name:  dbUser.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.RefreshTokenAge) * time.Second)),
			ID:        dbUser.ID.String(),
			Subject:   "refresh_token",
		},
	})
	if err != nil {
		return "", "", err
	}
	access_token, err := jwtToken.Sign(jwtToken.Claims{
		Email: dbUser.Email,
		Name:  dbUser.Name,
		Image: dbUser.Image.String,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.AccessTokenAge) * time.Second)),
			ID:        dbUser.ID.String(),
			Subject:   "access_token",
		},
	})
	if err != nil {
		return "", "", nil
	}

	_, err = server.DBQueries.UpdateUserRefreshToken(context.Background(), database.UpdateUserRefreshTokenParams{
		Email:        dbUser.Email,
		RefreshToken: pgtype.Text{String: refresh_token, Valid: true},
	})
	if err != nil {
		return "", "", nil
	}

	return access_token, refresh_token, nil
}
