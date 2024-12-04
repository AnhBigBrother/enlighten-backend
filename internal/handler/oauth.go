package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/pkg/parser"
	"github.com/AnhBigBrother/enlighten-backend/pkg/resp"
	"github.com/AnhBigBrother/enlighten-backend/pkg/token"
	"github.com/golang-jwt/jwt/v5"
)

type OauthApi struct {
	DB *database.Queries
}

type oauthUserInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewOauthApi() OauthApi {
	return OauthApi{
		DB: cfg.DBQueries,
	}
}

func (oauthApi *OauthApi) OauthGoogle(w http.ResponseWriter, r *http.Request) {
	tokenType := r.URL.Query().Get("token_type")
	accessToken := r.URL.Query().Get("access_token")
	userData, err := oauthApi.callGoogle(tokenType, accessToken)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	access_token, refresh_token, err := oauthApi.signInOauthUser(userData)
	if err != nil {
		if err.Error() == "unregistered user" {
			resp.Json(w, 202, struct {
				Message  string        `json:"message"`
				UserInfo oauthUserInfo `json:"user_info"`
			}{
				Message:  "user need to registered",
				UserInfo: userData,
			})
			return
		}
		resp.Err(w, 500, err.Error())
		return
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	resp.Json(w, 200, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: refresh_token})
}

func (oauthApi *OauthApi) OauthGithub(w http.ResponseWriter, r *http.Request) {
	tokenType := r.URL.Query().Get("token_type")
	accessToken := r.URL.Query().Get("access_token")
	userData, err := oauthApi.callGithub(tokenType, accessToken)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	access_token, refresh_token, err := oauthApi.signInOauthUser(userData)
	if err != nil {
		if err.Error() == "unregistered user" {
			resp.Json(w, 202, struct {
				Message  string        `json:"message"`
				UserInfo oauthUserInfo `json:"user_info"`
			}{
				Message:  "user need to registered",
				UserInfo: userData,
			})
			return
		}
		resp.Err(w, 500, err.Error())
		return
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	resp.Json(w, 200, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: refresh_token})
}

func (oauthApi *OauthApi) OauthMicrosoft(w http.ResponseWriter, r *http.Request) {
	tokenType := r.URL.Query().Get("token_type")
	accessToken := r.URL.Query().Get("access_token")
	userData, err := oauthApi.callMicrosoft(tokenType, accessToken)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	access_token, refresh_token, err := oauthApi.signInOauthUser(userData)
	if err != nil {
		if err.Error() == "unregistered user" {
			resp.Json(w, 202, struct {
				Message  string        `json:"message"`
				UserInfo oauthUserInfo `json:"user_info"`
			}{
				Message:  "user need to registered",
				UserInfo: userData,
			})
			return
		}
		resp.Err(w, 500, err.Error())
		return
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	resp.Json(w, 200, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: refresh_token})
}

func (oauthApi *OauthApi) OauthDiscord(w http.ResponseWriter, r *http.Request) {
	tokenType := r.URL.Query().Get("token_type")
	accessToken := r.URL.Query().Get("access_token")
	userData, err := oauthApi.callDiscord(tokenType, accessToken)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	access_token, refresh_token, err := oauthApi.signInOauthUser(userData)
	if err != nil {
		if err.Error() == "unregistered user" {
			resp.Json(w, 202, struct {
				Message  string        `json:"message"`
				UserInfo oauthUserInfo `json:"user_info"`
			}{
				Message:  "user need to registered",
				UserInfo: userData,
			})
			return
		}
		resp.Err(w, 500, err.Error())
		return
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	resp.Json(w, 200, struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{AccessToken: access_token, RefreshToken: refresh_token})
}

func (oauthApi *OauthApi) callGoogle(tokenType, accessToken string) (oauthUserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return oauthUserInfo{}, errors.New("token_type and access_token are required")
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", cfg.GoogleGetUserDataUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	res, err := client.Do(req)
	if err != nil {
		return oauthUserInfo{}, err
	}
	userInfo := oauthUserInfo{}
	err = parser.ParseBody(res.Body, &userInfo)
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (oauthApi *OauthApi) callGithub(tokenType, accessToken string) (oauthUserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return oauthUserInfo{}, errors.New("token_type and access_token are required")
	}

	client := &http.Client{}

	req, _ := http.NewRequest("GET", cfg.GithubGetUserDataUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	res, err := client.Do(req)
	if err != nil {
		return oauthUserInfo{}, err
	}
	userInfo := struct {
		Name      string `json:"name"`
		AvatarUrl string `json:"avatar_url"`
	}{}
	err = parser.ParseBody(res.Body, &userInfo)
	if err != nil {
		return oauthUserInfo{}, err
	}

	emailReq, _ := http.NewRequest("GET", cfg.GithubGetUserDataUrl+"/emails", nil)
	emailReq.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	emailRes, err := client.Do(emailReq)
	if err != nil {
		return oauthUserInfo{}, err
	}
	emailResData := []interface{}{}
	err = parser.ParseBody(emailRes.Body, &emailResData)
	if err != nil {
		return oauthUserInfo{}, err
	}

	return oauthUserInfo{
		Email:   emailResData[0].(map[string]interface{})["email"].(string),
		Name:    userInfo.Name,
		Picture: userInfo.AvatarUrl,
	}, nil
}

func (oauthApi *OauthApi) callMicrosoft(tokenType, accessToken string) (oauthUserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return oauthUserInfo{}, errors.New("token_type and access_token are required")
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", cfg.MicrosoftGetUserDataUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	res, err := client.Do(req)
	if err != nil {
		return oauthUserInfo{}, err
	}
	userInfo := oauthUserInfo{}
	err = parser.ParseBody(res.Body, &userInfo)
	if err != nil {
		return oauthUserInfo{}, err
	}
	return userInfo, nil
}

func (oauthApi *OauthApi) callDiscord(tokenType, accessToken string) (oauthUserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return oauthUserInfo{}, errors.New("token_type and access_token are required")
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", cfg.DiscordGetUserDataUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	res, err := client.Do(req)
	if err != nil {
		return oauthUserInfo{}, err
	}
	userInfo := struct {
		Id         string `json:"id"`
		Email      string `json:"email"`
		GlobalName string `json:"global_name"`
		Avatar     string `json:"avatar"`
	}{}
	err = parser.ParseBody(res.Body, &userInfo)
	if err != nil {
		return oauthUserInfo{}, err
	}
	return oauthUserInfo{
		Email:   userInfo.Email,
		Name:    userInfo.GlobalName,
		Picture: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s", userInfo.Id, userInfo.Avatar),
	}, nil
}

func (oauthApi *OauthApi) signInOauthUser(user oauthUserInfo) (string, string, error) {
	dbUser, err := oauthApi.DB.FindUserByEmail(context.Background(), user.Email)

	if err != nil {
		return "", "", errors.New("unregistered user")
	}

	currentTime := time.Now()
	refresh_token, err := token.Sign(token.Claims{
		Email: user.Email,
		Name:  user.Name,
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
	access_token, err := token.Sign(token.Claims{
		Email: user.Email,
		Name:  user.Name,
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

	_, err = oauthApi.DB.UpdateUserRefreshToken(context.Background(), database.UpdateUserRefreshTokenParams{
		Email:        dbUser.Email,
		RefreshToken: sql.NullString{String: refresh_token, Valid: true},
	})
	if err != nil {
		return "", "", nil
	}

	return access_token, refresh_token, nil
}
