package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
)

type UserInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func GetGoogleUserInfo(tokenType, accessToken string) (UserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return UserInfo{}, errors.New("oauth token failed")
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", cfg.GoogleGetUserDataUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	res, err := client.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	userInfo := UserInfo{}
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func GetGithubUserInfo(tokenType, accessToken string) (UserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return UserInfo{}, errors.New("oauth token failed")
	}

	client := &http.Client{}

	req, _ := http.NewRequest("GET", cfg.GithubGetUserDataUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	res, err := client.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	userInfo := struct {
		Name      string `json:"name"`
		AvatarUrl string `json:"avatar_url"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return UserInfo{}, err
	}

	emailReq, _ := http.NewRequest("GET", cfg.GithubGetUserDataUrl+"/emails", nil)
	emailReq.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	emailRes, err := client.Do(emailReq)
	if err != nil {
		return UserInfo{}, err
	}
	emailResData := []interface{}{}
	err = json.NewDecoder(emailRes.Body).Decode(&emailResData)
	if err != nil {
		return UserInfo{}, err
	}

	return UserInfo{
		Email:   emailResData[0].(map[string]interface{})["email"].(string),
		Name:    userInfo.Name,
		Picture: userInfo.AvatarUrl,
	}, nil
}

func GetMicrosoftUserInfo(tokenType, accessToken string) (UserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return UserInfo{}, errors.New("oauth token failed")
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", cfg.MicrosoftGetUserDataUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	res, err := client.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	userInfo := UserInfo{}
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return UserInfo{}, err
	}
	return userInfo, nil
}

func GetDiscordUserInfo(tokenType, accessToken string) (UserInfo, error) {
	if tokenType == "" || accessToken == "" {
		return UserInfo{}, errors.New("oauth token failed")
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", cfg.DiscordGetUserDataUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	res, err := client.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	userInfo := struct {
		Id         string `json:"id"`
		Email      string `json:"email"`
		GlobalName string `json:"global_name"`
		Avatar     string `json:"avatar"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return UserInfo{}, err
	}
	return UserInfo{
		Email:   userInfo.Email,
		Name:    userInfo.GlobalName,
		Picture: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s", userInfo.Id, userInfo.Avatar),
	}, nil
}
