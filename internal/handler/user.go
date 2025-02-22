package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/dto"
	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/parser"
	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/resp"
	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	DB *database.Queries
}

func NewUsersHandler() UsersHandler {
	return UsersHandler{
		DB: cfg.DBQueries,
	}
}

func (usersHandler *UsersHandler) SignUp(w http.ResponseWriter, r *http.Request) {
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
		ID:           pgtype.UUID{Bytes: userId, Valid: true},
		Email:        params.Email,
		Name:         params.Name,
		Password:     string(hashedPassword),
		RefreshToken: pgtype.Text{String: refresh_token, Valid: true},
		CreatedAt:    pgtype.Timestamp{Time: currentTime, InfinityModifier: pgtype.Finite, Valid: true},
		UpdatedAt:    pgtype.Timestamp{Time: currentTime, InfinityModifier: pgtype.Finite, Valid: true},
	}
	if params.Image != "" {
		createUserParams.Image = pgtype.Text{String: params.Image, Valid: true}
	}

	_, err = usersHandler.DB.CreateUser(r.Context(), createUserParams)
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

func (usersHandler *UsersHandler) SignIn(w http.ResponseWriter, r *http.Request) {
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

	user, err := usersHandler.DB.FindUserByEmail(r.Context(), params.Email)
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

	_, err = usersHandler.DB.UpdateUserRefreshToken(r.Context(), database.UpdateUserRefreshTokenParams{
		Email:        params.Email,
		RefreshToken: pgtype.Text{String: refresh_token, Valid: true},
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

func (usersHandler *UsersHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	sessionEmail := session["email"].(string)

	_, err := usersHandler.DB.UpdateUserRefreshToken(r.Context(), database.UpdateUserRefreshTokenParams{
		Email:        sessionEmail,
		RefreshToken: pgtype.Text{Valid: false},
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

func (usersHandler *UsersHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	sessionEmail := session["email"].(string)
	currUser, err := usersHandler.DB.FindUserByEmail(r.Context(), sessionEmail)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}

	resp.Json(w, 200, currUser)
}

func (usersHandler *UsersHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
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

	session := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	sessionEmail := session["email"].(string)
	user, err := usersHandler.DB.FindUserByEmail(r.Context(), sessionEmail)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	updateUserInfoParams := database.UpdateUserInfoParams{
		Email:     sessionEmail,
		Name:      user.Name,
		Image:     user.Image,
		Password:  user.Password,
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
	}
	if len(params.Password) > 0 {
		updateUserInfoParams.Password = params.Password
	}
	if len(params.Name) > 0 {
		updateUserInfoParams.Name = params.Name
	}
	if len(params.Image) > 0 {
		updateUserInfoParams.Image = pgtype.Text{String: params.Image, Valid: true}
	}
	if len(params.Bio) > 0 {
		updateUserInfoParams.Bio = pgtype.Text{String: params.Bio, Valid: true}
	}
	_, err = usersHandler.DB.UpdateUserInfo(r.Context(), updateUserInfoParams)
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

func (usersHandler *UsersHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	session := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	sessionEmail := session["email"].(string)

	user, err := usersHandler.DB.FindUserByEmail(r.Context(), sessionEmail)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}

	err = usersHandler.DB.DeleteUserInfo(r.Context(), sessionEmail)
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

func (usersHandler *UsersHandler) GetSesion(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	resp.Json(w, 200, session)
}

func (usersHandler *UsersHandler) GetAccessToken(w http.ResponseWriter, r *http.Request) {
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
		return
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

	user, err := usersHandler.DB.UpdateUserRefreshToken(r.Context(), database.UpdateUserRefreshTokenParams{
		Email:        claims["email"].(string),
		RefreshToken: pgtype.Text{String: new_refresh_token, Valid: true},
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

func (usersHandler *UsersHandler) GetMyOverview(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	userId := user["jti"].(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	overview, err := usersHandler.DB.GetUserOverview(r.Context(), pgtype.UUID{Bytes: userUUID, Valid: true})
	if err != nil {
		resp.Err(w, 404, "user not found")
		return
	}
	resp.Json(w, 200, struct {
		ID             pgtype.UUID      `json:"id"`
		Name           string           `json:"name"`
		Email          string           `json:"email"`
		Image          string           `json:"image"`
		Bio            string           `json:"bio"`
		TotalPosts     int32            `json:"total_posts"`
		TotalUpvoted   int32            `json:"total_upvoted"`
		TotalDownvoted int32            `json:"total_downvoted"`
		Follower       int32            `json:"follower"`
		Following      int32            `json:"following"`
		CreatedAt      pgtype.Timestamp `json:"created_at"`
		UpdatedAt      pgtype.Timestamp `json:"updated_at"`
	}{
		ID:             overview.ID,
		Name:           overview.Name,
		Email:          overview.Email,
		Image:          overview.Image.String,
		Bio:            overview.Bio.String,
		TotalPosts:     overview.TotalPosts,
		TotalUpvoted:   overview.TotalUpvoted,
		TotalDownvoted: overview.TotalDownvoted,
		Follower:       overview.Follower,
		Following:      overview.Following,
		CreatedAt:      overview.CreatedAt,
		UpdatedAt:      overview.UpdatedAt,
	})
}

func (usersHandler *UsersHandler) GetMyPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	userId := user["jti"].(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	queryParams := r.URL.Query()
	sort, limitStr, offsetStr := queryParams.Get("sort"), queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	type JsonPost struct {
		ID            pgtype.UUID      `json:"id"`
		Title         string           `json:"title"`
		Content       string           `json:"content"`
		AuthorID      pgtype.UUID      `json:"author_id"`
		UpVoted       int32            `json:"up_voted"`
		DownVoted     int32            `json:"down_voted"`
		CommentsCount int32            `json:"comments_count"`
		CreatedAt     pgtype.Timestamp `json:"created_at"`
		UpdatedAt     pgtype.Timestamp `json:"updated_at"`
		AuthorName    string           `json:"author_name"`
		AuthorEmail   string           `json:"author_email"`
		AuthorImage   string           `json:"author_image"`
	}
	if sort == "hot" {
		posts, _ := usersHandler.DB.GetUserHotPosts(r.Context(), database.GetUserHotPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(limit), Offset: int32(offset)})
		jsonPosts := []JsonPost{}
		for _, p := range posts {
			jsonPosts = append(jsonPosts, JsonPost{
				ID:            p.ID,
				Title:         p.Title,
				Content:       p.Content,
				AuthorID:      p.AuthorID,
				UpVoted:       p.UpVoted,
				DownVoted:     p.DownVoted,
				CommentsCount: p.CommentsCount,
				CreatedAt:     p.CreatedAt,
				UpdatedAt:     p.UpdatedAt,
				AuthorName:    p.AuthorName,
				AuthorEmail:   p.AuthorEmail,
				AuthorImage:   p.AuthorImage.String,
			})
		}
		resp.Json(w, 200, jsonPosts)
		return
	}
	if sort == "top" {
		posts, _ := usersHandler.DB.GetUserTopPosts(r.Context(), database.GetUserTopPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(limit), Offset: int32(offset)})
		jsonPosts := []JsonPost{}
		for _, p := range posts {
			jsonPosts = append(jsonPosts, JsonPost{
				ID:            p.ID,
				Title:         p.Title,
				Content:       p.Content,
				AuthorID:      p.AuthorID,
				UpVoted:       p.UpVoted,
				DownVoted:     p.DownVoted,
				CommentsCount: p.CommentsCount,
				CreatedAt:     p.CreatedAt,
				UpdatedAt:     p.UpdatedAt,
				AuthorName:    p.AuthorName,
				AuthorEmail:   p.AuthorEmail,
				AuthorImage:   p.AuthorImage.String,
			})
		}
		resp.Json(w, 200, jsonPosts)
		return
	}
	posts, _ := usersHandler.DB.GetUserNewPosts(r.Context(), database.GetUserNewPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(limit), Offset: int32(offset)})
	jsonPosts := []JsonPost{}
	for _, p := range posts {
		jsonPosts = append(jsonPosts, JsonPost{
			ID:            p.ID,
			Title:         p.Title,
			Content:       p.Content,
			AuthorID:      p.AuthorID,
			UpVoted:       p.UpVoted,
			DownVoted:     p.DownVoted,
			CommentsCount: p.CommentsCount,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			AuthorName:    p.AuthorName,
			AuthorEmail:   p.AuthorEmail,
			AuthorImage:   p.AuthorImage.String,
		})
	}
	resp.Json(w, 200, jsonPosts)
}

func (usersHandler *UsersHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("user_id")
	if user_id == "" {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	overview, err := usersHandler.DB.GetUserOverview(r.Context(), pgtype.UUID{Bytes: userUUID, Valid: true})
	if err != nil {
		resp.Err(w, 404, "user not found")
		return
	}
	resp.Json(w, 200, struct {
		ID             pgtype.UUID      `json:"id"`
		Name           string           `json:"name"`
		Email          string           `json:"email"`
		Image          string           `json:"image"`
		Bio            string           `json:"bio"`
		TotalPosts     int32            `json:"total_posts"`
		TotalUpvoted   int32            `json:"total_upvoted"`
		TotalDownvoted int32            `json:"total_downvoted"`
		Follower       int32            `json:"follower"`
		Following      int32            `json:"following"`
		CreatedAt      pgtype.Timestamp `json:"created_at"`
		UpdatedAt      pgtype.Timestamp `json:"updated_at"`
	}{
		ID:             overview.ID,
		Name:           overview.Name,
		Email:          overview.Email,
		Image:          overview.Image.String,
		Bio:            overview.Bio.String,
		TotalPosts:     overview.TotalPosts,
		TotalUpvoted:   overview.TotalUpvoted,
		TotalDownvoted: overview.TotalDownvoted,
		Follower:       overview.Follower,
		Following:      overview.Following,
		CreatedAt:      overview.CreatedAt,
		UpdatedAt:      overview.UpdatedAt,
	})
}

func (usersHandler *UsersHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("user_id")
	if user_id == "" {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	queryParams := r.URL.Query()
	sort, limitStr, offsetStr := queryParams.Get("sort"), queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	type JsonPost struct {
		ID            pgtype.UUID      `json:"id"`
		Title         string           `json:"title"`
		Content       string           `json:"content"`
		AuthorID      pgtype.UUID      `json:"author_id"`
		UpVoted       int32            `json:"up_voted"`
		DownVoted     int32            `json:"down_voted"`
		CommentsCount int32            `json:"comments_count"`
		CreatedAt     pgtype.Timestamp `json:"created_at"`
		UpdatedAt     pgtype.Timestamp `json:"updated_at"`
		AuthorName    string           `json:"author_name"`
		AuthorEmail   string           `json:"author_email"`
		AuthorImage   string           `json:"author_image"`
	}
	if sort == "hot" {
		posts, _ := usersHandler.DB.GetUserHotPosts(r.Context(), database.GetUserHotPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(limit), Offset: int32(offset)})
		jsonPosts := []JsonPost{}
		for _, p := range posts {
			jsonPosts = append(jsonPosts, JsonPost{
				ID:            p.ID,
				Title:         p.Title,
				Content:       p.Content,
				AuthorID:      p.AuthorID,
				UpVoted:       p.UpVoted,
				DownVoted:     p.DownVoted,
				CommentsCount: p.CommentsCount,
				CreatedAt:     p.CreatedAt,
				UpdatedAt:     p.UpdatedAt,
				AuthorName:    p.AuthorName,
				AuthorEmail:   p.AuthorEmail,
				AuthorImage:   p.AuthorImage.String,
			})
		}
		resp.Json(w, 200, jsonPosts)
		return
	}
	if sort == "top" {
		posts, _ := usersHandler.DB.GetUserTopPosts(r.Context(), database.GetUserTopPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(limit), Offset: int32(offset)})
		jsonPosts := []JsonPost{}
		for _, p := range posts {
			jsonPosts = append(jsonPosts, JsonPost{
				ID:            p.ID,
				Title:         p.Title,
				Content:       p.Content,
				AuthorID:      p.AuthorID,
				UpVoted:       p.UpVoted,
				DownVoted:     p.DownVoted,
				CommentsCount: p.CommentsCount,
				CreatedAt:     p.CreatedAt,
				UpdatedAt:     p.UpdatedAt,
				AuthorName:    p.AuthorName,
				AuthorEmail:   p.AuthorEmail,
				AuthorImage:   p.AuthorImage.String,
			})
		}
		resp.Json(w, 200, jsonPosts)
		return
	}
	posts, _ := usersHandler.DB.GetUserNewPosts(r.Context(), database.GetUserNewPostsParams{ID: pgtype.UUID{Bytes: userUUID, Valid: true}, Limit: int32(limit), Offset: int32(offset)})
	jsonPosts := []JsonPost{}
	for _, p := range posts {
		jsonPosts = append(jsonPosts, JsonPost{
			ID:            p.ID,
			Title:         p.Title,
			Content:       p.Content,
			AuthorID:      p.AuthorID,
			UpVoted:       p.UpVoted,
			DownVoted:     p.DownVoted,
			CommentsCount: p.CommentsCount,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			AuthorName:    p.AuthorName,
			AuthorEmail:   p.AuthorEmail,
			AuthorImage:   p.AuthorImage.String,
		})
	}
	resp.Json(w, 200, jsonPosts)
}

func (usersHandler *UsersHandler) Follow(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	follower_id := session["jti"].(string)
	follower_uuid, _ := uuid.Parse(follower_id)
	user_id := r.PathValue("user_id")
	if user_id == "" {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	err = usersHandler.DB.CreateFollows(r.Context(), database.CreateFollowsParams{
		ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
		AuthorID:   pgtype.UUID{Bytes: userUUID, Valid: true},
		FollowerID: pgtype.UUID{Bytes: follower_uuid, Valid: true},
		CreatedAt:  pgtype.Timestamp{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 201, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (usersHandler *UsersHandler) UnFollow(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	follower_id := session["jti"].(string)
	follower_uuid, _ := uuid.Parse(follower_id)
	user_id := r.PathValue("user_id")
	if user_id == "" {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	err = usersHandler.DB.DeleteFollows(r.Context(), database.DeleteFollowsParams{
		AuthorID:   pgtype.UUID{Bytes: userUUID, Valid: true},
		FollowerID: pgtype.UUID{Bytes: follower_uuid, Valid: true},
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 200, struct {
		Message string `json:"message"`
	}{Message: "success"})
}

func (usersHandler *UsersHandler) CheckFollowed(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	follower_id := session["jti"].(string)
	follower_uuid, _ := uuid.Parse(follower_id)
	user_id := r.PathValue("user_id")
	userUUID, err := uuid.Parse(user_id)
	if err != nil {
		resp.Err(w, 400, "invalid user_id")
		return
	}
	follow, err := usersHandler.DB.GetFollows(r.Context(), database.GetFollowsParams{
		AuthorID:   pgtype.UUID{Bytes: userUUID, Valid: true},
		FollowerID: pgtype.UUID{Bytes: follower_uuid, Valid: true},
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 200, struct {
		ID         pgtype.UUID      `json:"id"`
		FollowerID pgtype.UUID      `json:"follower_id"`
		AuthorID   pgtype.UUID      `json:"author_id"`
		CreatedAt  pgtype.Timestamp `json:"created_at"`
	}{
		ID:         follow.ID,
		FollowerID: follow.FollowerID,
		AuthorID:   follow.AuthorID,
		CreatedAt:  follow.CreatedAt,
	})
}

func (usersHandler *UsersHandler) GetFollowedAuthor(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(cfg.CtxKeys.User).(map[string]interface{})
	userId := user["jti"].(string)
	userUUID, _ := uuid.Parse(userId)
	queryParams := r.URL.Query()
	limitStr, offsetStr := queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	followedAuthors, err := usersHandler.DB.GetFollowedAuthor(r.Context(), database.GetFollowedAuthorParams{
		FollowerID: pgtype.UUID{Bytes: userUUID, Valid: true},
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 200, followedAuthors)
}

func (usersHandler *UsersHandler) GetTopAuthor(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limitStr, offsetStr := queryParams.Get("limit"), queryParams.Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	topAuthors, err := usersHandler.DB.GetTopAuthor(r.Context(), database.GetTopAuthorParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 200, topAuthors)
}
