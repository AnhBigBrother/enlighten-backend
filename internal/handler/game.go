package handler

import (
	"net/http"
	"strconv"

	"github.com/AnhBigBrother/enlighten-backend/pkg/parser"
	"github.com/AnhBigBrother/enlighten-backend/pkg/resp"
	"github.com/AnhBigBrother/enlighten-backend/pkg/sudoku"
)

type Sudoku struct{}

func (s *Sudoku) SolveSudoku(w http.ResponseWriter, r *http.Request) {
	su := struct {
		Board [][]int `json:"board"`
	}{}
	err := parser.ParseBody(r.Body, &su)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	result, err := sudoku.Solve(su.Board)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	resp.Json(w, 200, struct {
		Result [][]int `json:"result"`
	}{Result: result})
}

func (s *Sudoku) GenerateSudoku(w http.ResponseWriter, r *http.Request) {
	hideStr := r.URL.Query().Get("hide")
	hide, err := strconv.Atoi(hideStr)
	if err != nil || hide < 0 || hide > 81 {
		hide = 45
	}
	board := sudoku.Gen(hide)
	resp.Json(w, 200, struct {
		Board [][]int `json:"board"`
	}{Board: board})
}

func (s *Sudoku) CheckSolvable(w http.ResponseWriter, r *http.Request) {
	su := struct {
		Board [][]int `json:"board"`
	}{}
	err := parser.ParseBody(r.Body, &su)
	if err != nil {
		resp.Err(w, 400, err.Error())
		return
	}
	_, err = sudoku.Solve(su.Board)
	if err != nil {
		resp.Json(w, 200, struct {
			Solvable bool `json:"solvable"`
		}{
			Solvable: false,
		})
		return
	}
	resp.Json(w, 200, struct {
		Solvable bool `json:"solvable"`
	}{
		Solvable: true,
	})
}
