package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/dto"
	"github.com/AnhBigBrother/enlighten-backend/internal/models"
	"github.com/AnhBigBrother/enlighten-backend/pkg/parser"
	"github.com/AnhBigBrother/enlighten-backend/pkg/resp"
	"github.com/AnhBigBrother/enlighten-backend/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	err := parser.ParseBody(r.Body, &params)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	if err := params.ValidateInput(); err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	userId := uuid.New()
	currentTime := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		resp.Err(w, 500, err.Error())
		return
	}

	refresh_token, err := token.Sign(token.Claims{
		Email: params.Email,
		Name:  params.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.RefreshTokenAge) * time.Second)),
			ID:        userId.String(),
			Subject:   "refresh_token",
		},
	})
	if err != nil {
		resp.Err(w, 500, err.Error())
		return
	}

	createUserParams := database.CreateUserParams{
		ID:           userId,
		Email:        params.Email,
		Name:         params.Name,
		Password:     string(hashedPassword),
		RefreshToken: sql.NullString{String: refresh_token, Valid: true},
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	if params.Image != "" {
		createUserParams.Image = sql.NullString{String: params.Image, Valid: true}
	}

	_, err = userApi.DB.CreateUser(r.Context(), createUserParams)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	access_token, err := token.Sign(token.Claims{
		Email: params.Email,
		Name:  params.Name,
		Image: params.Image,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.AccessTokenAge) * time.Second)),
			ID:        userId.String(),
			Subject:   "access_token",
		},
	})
	if err != nil {
		resp.Err(w, 500, err.Error())
		return
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	resp.Json(w, 201, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: refresh_token})
}

func (userApi *UserApi) SignIn(w http.ResponseWriter, r *http.Request) {
	params := dto.UserLogIn{}
	err := parser.ParseBody(r.Body, &params)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	if err := params.ValidateInput(); err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	user, err := userApi.DB.FindUserByEmail(r.Context(), params.Email)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		resp.Err(w, 401, err.Error())
		return
	}

	currentTime := time.Now()
	access_token, err := token.Sign(token.Claims{
		Email: params.Email,
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
		resp.Err(w, 500, err.Error())
		return
	}
	refresh_token, err := token.Sign(token.Claims{
		Email: params.Email,
		Name:  user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.RefreshTokenAge) * time.Second)),
			ID:        user.ID.String(),
			Subject:   "refresh_token",
		},
	})
	if err != nil {
		resp.Err(w, 500, err.Error())
		return
	}

	_, err = userApi.DB.UpdateUserRefreshToken(r.Context(), database.UpdateUserRefreshTokenParams{
		Email:        params.Email,
		RefreshToken: sql.NullString{String: refresh_token, Valid: true},
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	resp.Json(w, 200, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: refresh_token})
}

func (userApi *UserApi) SignOut(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
	sessionEmail := session["email"].(string)

	_, err := userApi.DB.UpdateUserRefreshToken(r.Context(), database.UpdateUserRefreshTokenParams{
		Email:        sessionEmail,
		RefreshToken: sql.NullString{String: "", Valid: false},
	})

	if err != nil {
		log.Println(err.Error())
	}

	resp.DeleteCookie(w, "access_token")
	resp.DeleteCookie(w, "refresh_token")

	resp.Json(w, 200, struct {
		Message string `json:"message"`
	}{Message: "Signed out"})
}

func (userApi *UserApi) GetMe(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
	sessionEmail := session["email"].(string)
	currUser, err := userApi.DB.FindUserByEmail(r.Context(), sessionEmail)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}

	resp.Json(w, 200, models.FormatDatabaseUser(currUser))
}

func (userApi *UserApi) UpdateMe(w http.ResponseWriter, r *http.Request) {
	params := dto.UserUpdate{}
	err := parser.ParseBody(r.Body, &params)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	if err := params.ValidateInput(); err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	if len(params.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			resp.Err(w, 400, err.Error())
			return
		}
		params.Password = string(hashedPassword)
	}

	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
	sessionEmail := session["email"].(string)
	user, err := userApi.DB.FindUserByEmail(r.Context(), sessionEmail)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	updateUserInfoParams := database.UpdateUserInfoParams{
		Email:     sessionEmail,
		Name:      user.Name,
		Image:     user.Image,
		Password:  user.Password,
		UpdatedAt: time.Now(),
	}
	updateUserInfoParams.UpdatedAt = time.Now()
	if len(params.Password) > 0 {
		updateUserInfoParams.Password = params.Password
	}
	if len(params.Name) > 0 {
		updateUserInfoParams.Name = params.Name
	}
	if len(params.Image) > 0 {
		updateUserInfoParams.Image = sql.NullString{String: params.Image, Valid: true}
	}
	if len(params.Bio) > 0 {
		updateUserInfoParams.Bio = sql.NullString{String: params.Bio, Valid: true}
	}
	_, err = userApi.DB.UpdateUserInfo(r.Context(), updateUserInfoParams)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	access_token, err := token.Sign(token.Claims{
		Email: sessionEmail,
		Name:  params.Name,
		Image: params.Image,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Unix(int64(session["iat"].(float64)), 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(int64(session["exp"].(float64)), 0)),
			ID:        user.ID.String(),
			Subject:   "access_token",
		},
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	resp.SetCookie(w, "access_token", access_token)

	resp.Json(w, 200, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: user.RefreshToken.String})
}

func (userApi *UserApi) DeleteMe(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
	sessionEmail := session["email"].(string)

	user, err := userApi.DB.FindUserByEmail(r.Context(), sessionEmail)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	err = userApi.DB.DeleteUserInfo(r.Context(), sessionEmail)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	resp.DeleteCookie(w, "access_token")
	resp.DeleteCookie(w, "refresh_token")

	resp.Json(w, 200, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (userApi *UserApi) GetSesion(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		log.Println("Server error: route must nested inside auth middleware")
		resp.Json(w, 500, "server error: something went wrong")
		return
	}
	resp.Json(w, 200, session)
}

func (userApi *UserApi) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	refresh_token := r.URL.Query().Get("refresh_token")
	if refresh_token == "" {
		cookie, err := r.Cookie("refresh_token")
		if err == nil {
			refresh_token = cookie.Value
		}
	}
	if refresh_token == "" {
		resp.Err(w, 400, "Missing parameter: refresh_token")
		return
	}
	claims, err := token.Parse(refresh_token)
	if err != nil {
		resp.Err(w, 400, err.Error())
	}

	currentTime := time.Now()
	if int64(claims["exp"].(float64)) < currentTime.Unix() {
		resp.Err(w, 403, "refresh_token expired")
		return
	}

	new_refresh_token, err := token.Sign(token.Claims{
		Email: claims["email"].(string),
		Name:  claims["name"].(string),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.RefreshTokenAge) * time.Second)),
			ID:        claims["jti"].(string),
			Subject:   "refresh_token",
		},
	})
	if err != nil {
		resp.Err(w, 403, err.Error())
		return
	}

	user, err := userApi.DB.UpdateUserRefreshToken(r.Context(), database.UpdateUserRefreshTokenParams{
		Email:        claims["email"].(string),
		RefreshToken: sql.NullString{String: new_refresh_token, Valid: true},
	})
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}

	access_token, err := token.Sign(token.Claims{
		Email: claims["email"].(string),
		Name:  claims["name"].(string),
		Image: user.Image.String,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(cfg.AccessTokenAge) * time.Second)),
			ID:        claims["jti"].(string),
			Subject:   "access_token",
		},
	})
	if err != nil {
		resp.Err(w, 500, err.Error())
		return
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", new_refresh_token)

	resp.Json(w, 201, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: new_refresh_token})
}
