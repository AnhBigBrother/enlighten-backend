package handler

import (
	"bytes"
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

type oauthToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

func NewOauthApi() OauthApi {
	return OauthApi{
		DB: cfg.DBQueries,
	}
}

func (oauthApi *OauthApi) HandleGoogleOauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tokenType, accessToken, err := getGoogleAccessToken(code)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	userData, err := getGoogleUserInfo(tokenType, accessToken)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	access_token, refresh_token, err := oauthApi.signInOauthUser(userData)
	if err != nil {
		if err.Error() == "unregistered user" {
			http.Redirect(w, r, fmt.Sprintf("%s/signup?email=%s&name=%s&image=%s", cfg.Frontend, userData.Email, userData.Name, userData.Picture), http.StatusSeeOther)
			return
		}
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	http.Redirect(w, r, cfg.Frontend, http.StatusSeeOther)
}

func (oauthApi *OauthApi) HandleGithubOauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tokenType, accessToken, err := getGithubAccessToken(code)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	userData, err := getGithubUserInfo(tokenType, accessToken)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	access_token, refresh_token, err := oauthApi.signInOauthUser(userData)
	if err != nil {
		if err.Error() == "unregistered user" {
			http.Redirect(w, r, fmt.Sprintf("%s/signup?email=%s&name=%s&image=%s", cfg.Frontend, userData.Email, userData.Name, userData.Picture), http.StatusSeeOther)
			return
		}
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	http.Redirect(w, r, cfg.Frontend, http.StatusSeeOther)
}

func (oauthApi *OauthApi) HandleMicrosoftOauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tokenType, accessToken, err := getMicrosoftAccessToken(code)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	userData, err := getMicrosoftUserInfo(tokenType, accessToken)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	access_token, refresh_token, err := oauthApi.signInOauthUser(userData)
	if err != nil {
		if err.Error() == "unregistered user" {
			http.Redirect(w, r, fmt.Sprintf("%s/signup?email=%s&name=%s&image=%s", cfg.Frontend, userData.Email, userData.Name, userData.Picture), http.StatusSeeOther)
			return
		}
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	http.Redirect(w, r, cfg.Frontend, http.StatusSeeOther)
}

func (oauthApi *OauthApi) HandleDiscordOauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tokenType, accessToken, err := getDiscordAccessToken(code)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	userData, err := getDiscordUserInfo(tokenType, accessToken)
	if err != nil {
		resp.Err(w, 404, err.Error())
		return
	}
	access_token, refresh_token, err := oauthApi.signInOauthUser(userData)
	if err != nil {
		if err.Error() == "unregistered user" {
			http.Redirect(w, r, fmt.Sprintf("%s/signup?email=%s&name=%s&image=%s", cfg.Frontend, userData.Email, userData.Name, userData.Picture), http.StatusSeeOther)
			return
		}
	}

	resp.SetCookie(w, "access_token", access_token)
	resp.SetCookie(w, "refresh_token", refresh_token)

	http.Redirect(w, r, cfg.Frontend, http.StatusSeeOther)
}

func getGoogleAccessToken(code string) (string, string, error) {
	client := &http.Client{}
	reqBody := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s", code, cfg.GoogleClientId, cfg.GoogleClientSecret, cfg.GoogleCallbackUrl, "authorization_code")
	req, _ := http.NewRequest("POST", cfg.GoogleGetTokenUrl, bytes.NewBufferString(reqBody))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", cfg.GoogleClientId)
	q.Add("client_secret", cfg.GoogleClientSecret)
	q.Add("redirect_uri", cfg.GoogleCallbackUrl)
	q.Add("grant_type", "authorization_code")
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	token := oauthToken{}
	err = parser.ParseBody(res.Body, &token)
	if err != nil {
		return "", "", err
	}
	return token.TokenType, token.AccessToken, nil
}

func getGithubAccessToken(code string) (string, string, error) {
	client := &http.Client{}
	reqBody := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s", code, cfg.GithubClientId, cfg.GithubClientSecret, cfg.GithubCallbackUrl, "authorization_code")
	req, _ := http.NewRequest("POST", cfg.GithubGetTokenUrl, bytes.NewBufferString(reqBody))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", cfg.GithubClientId)
	q.Add("client_secret", cfg.GithubClientSecret)
	q.Add("redirect_uri", cfg.GithubCallbackUrl)
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	token := oauthToken{}
	err = parser.ParseBody(res.Body, &token)
	if err != nil {
		return "", "", err
	}
	return token.TokenType, token.AccessToken, nil
}

func getMicrosoftAccessToken(code string) (string, string, error) {
	client := &http.Client{}
	reqBody := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s", code, cfg.MicrosoftClientId, cfg.MicrosoftClientSecret, cfg.MicrosoftCallbackUrl, "authorization_code")
	req, _ := http.NewRequest("POST", cfg.MicrosoftGetTokenUrl, bytes.NewBufferString(reqBody))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", cfg.MicrosoftClientId)
	q.Add("client_secret", cfg.MicrosoftClientSecret)
	q.Add("redirect_uri", cfg.MicrosoftCallbackUrl)
	q.Add("grant_type", "authorization_code")
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	token := oauthToken{}
	err = parser.ParseBody(res.Body, &token)
	if err != nil {
		return "", "", err
	}
	return token.TokenType, token.AccessToken, nil
}

func getDiscordAccessToken(code string) (string, string, error) {
	client := &http.Client{}
	reqBody := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s", code, cfg.DiscordClientId, cfg.DiscordClientSecret, cfg.DiscordCallbackUrl, "authorization_code")
	req, _ := http.NewRequest("POST", cfg.DiscordGetTokenUrl, bytes.NewBufferString(reqBody))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", cfg.DiscordClientId)
	q.Add("client_secret", cfg.DiscordClientSecret)
	q.Add("redirect_uri", cfg.DiscordCallbackUrl)
	q.Add("grant_type", "authorization_code")
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	token := oauthToken{}
	err = parser.ParseBody(res.Body, &token)
	if err != nil {
		return "", "", err
	}
	return token.TokenType, token.AccessToken, nil
}

func getGoogleUserInfo(tokenType, accessToken string) (oauthUserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return oauthUserInfo{}, errors.New("oauth token failed")
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

func getGithubUserInfo(tokenType, accessToken string) (oauthUserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return oauthUserInfo{}, errors.New("oauth token failed")
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

func getMicrosoftUserInfo(tokenType, accessToken string) (oauthUserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return oauthUserInfo{}, errors.New("oauth token failed")
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

func getDiscordUserInfo(tokenType, accessToken string) (oauthUserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return oauthUserInfo{}, errors.New("oauth token failed")
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
	access_token, err := token.Sign(token.Claims{
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

	_, err = oauthApi.DB.UpdateUserRefreshToken(context.Background(), database.UpdateUserRefreshTokenParams{
		Email:        dbUser.Email,
		RefreshToken: sql.NullString{String: refresh_token, Valid: true},
	})
	if err != nil {
		return "", "", nil
	}

	return access_token, refresh_token, nil
}
