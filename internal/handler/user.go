package handler

import (
	"enlighten-backend/cfg"
	"enlighten-backend/internal/database"
	"enlighten-backend/pkg/req"
	"enlighten-backend/pkg/resp"
	"fmt"
	"net/http"
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
	params := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := req.ParseBody(r, &params)
	if err != nil {
		resp.Err(w, 400, err.Error())
	}
	fmt.Println(params)
	resp.Json(w, 201, "")
}
