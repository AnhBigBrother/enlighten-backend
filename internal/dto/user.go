package dto

import (
	"errors"
	"regexp"
	"strings"
)

type UserSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user *UserSignUp) ValidateInput() error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	errArr := []string{}
	if !emailRegex.MatchString(user.Email) {
		errArr = append(errArr, "invalid email")
	}
	if len(user.Name) < 3 {
		errArr = append(errArr, "name too short")
	}
	if len(user.Password) < 6 {
		errArr = append(errArr, "password too short")
	}
	if len(errArr) > 0 {
		errMsg := strings.Join(errArr, ", ")
		return errors.New(errMsg)
	}
	return nil
}

func (user *UserLogIn) ValidateInput() error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	errArr := []string{}
	if !emailRegex.MatchString(user.Email) {
		errArr = append(errArr, "invalid email")
	}
	if len(user.Password) < 6 {
		errArr = append(errArr, "password too short")
	}
	if len(errArr) > 0 {
		errMsg := strings.Join(errArr, ", ")
		return errors.New(errMsg)
	}
	return nil
}
