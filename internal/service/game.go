package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/resp"
	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/sudoku"
)

type Sudoku struct{}

func NewSudokuService() Sudoku {
	return Sudoku{}
}

func (s *Sudoku) SolveSudoku(w http.ResponseWriter, r *http.Request) {
	su := struct {
		Board [][]int `json:"board"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&su)
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
	err := json.NewDecoder(r.Body).Decode(&su)
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
