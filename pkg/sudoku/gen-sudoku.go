package sudoku

import "math/rand/v2"

func Gen(hide int) [][]int {
	board := [][]int{}
	for i := 0; i < 9; i++ {
		board = append(board, make([]int, 9))
	}
	box1 := rand.Perm(9)
	idx := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			board[i][j] = box1[idx] + 1
			idx++
		}
	}
	box2 := rand.Perm(9)
	idx = 0
	for i := 3; i < 6; i++ {
		for j := 3; j < 6; j++ {
			board[i][j] = box2[idx] + 1
			idx++
		}
	}
	box3 := rand.Perm(9)
	idx = 0
	for i := 6; i < 9; i++ {
		for j := 6; j < 9; j++ {
			board[i][j] = box3[idx] + 1
			idx++
		}
	}
	board, _ = Solve(board)
	hidePos := rand.Perm(81)
	for i := 0; i < hide && i < 81; i++ {
		x, y := hidePos[i]/9, hidePos[i]%9
		board[x][y] = 0
	}
	return board
}
