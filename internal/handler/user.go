package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/dto"
	"github.com/AnhBigBrother/enlighten-backend/pkg/req"
	"github.com/AnhBigBrother/enlighten-backend/pkg/resp"
	"github.com/AnhBigBrother/enlighten-backend/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserApi struct {
	DB *database.Queries
}

func NewUserApi() UserApi {
	return UserApi{
		DB: cfg.DBQueries,
	}
}

func (userApi *UserApi) SignUp(w http.ResponseWriter, r *http.Request) {
	params := dto.UserSignUp{}
	err := req.ParseBody(r, &params)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	if err := params.ValidateInput(); err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	uuid := uuid.New()
	currentTime := time.Now()

	refresh_token, err := token.Sign(token.Claims{
		Email: params.Email,
		Name:  params.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.RefreshTokenAge) * time.Second)),
			ID:        uuid.String(),
			Subject:   "refresh_token",
		},
	})
	if err != nil {
		resp.Err(w, 500, err.Error())
		return
	}

	_, err = userApi.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid,
		Email:        params.Email,
		Name:         params.Name,
		Password:     params.Password,
		RefreshToken: sql.NullString{String: refresh_token, Valid: true},
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	access_token, err := token.Sign(token.Claims{
		Email: params.Email,
		Name:  params.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.AccessTokenAge) * time.Second)),
			ID:        uuid.String(),
			Subject:   "access_token",
		},
	})
	if err != nil {
		resp.Err(w, 500, err.Error())
		return
	}

	resp.SetCookie(w, "access_token", access_token)

	resp.Json(w, 201, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: refresh_token})
}
