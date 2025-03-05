package sudoku

import "errors"

func Solve(board [][]int) ([][]int, error) {
	cloneBoard := [][]int{}
	for i := 0; i < 9; i++ {
		row := []int{}
		for j := 0; j < 9; j++ {
			row = append(row, board[i][j])
		}
		cloneBoard = append(cloneBoard, row)
	}
	row, col, box := [][]bool{}, [][]bool{}, [][]bool{}
	for i := 0; i < 9; i++ {
		row = append(row, make([]bool, 10))
		col = append(col, make([]bool, 10))
		box = append(box, make([]bool, 10))
	}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			b := (i/3)*3 + j/3
			if cloneBoard[i][j] != 0 {
				if row[i][cloneBoard[i][j]] || col[j][cloneBoard[i][j]] || box[b][cloneBoard[i][j]] {
					return nil, errors.New("cloneBoard cannot be solved")
				}
				row[i][cloneBoard[i][j]] = true
				col[j][cloneBoard[i][j]] = true
				box[b][cloneBoard[i][j]] = true
			}
		}
	}
	flag := false
	var backtrack func(idx int)
	backtrack = func(idx int) {
		if idx > 80 || flag {
			flag = true
			return
		}
		i, j := idx/9, idx%9
		b := (i/3)*3 + j/3
		if cloneBoard[i][j] == 0 {
			for x := 1; x <= 9; x++ {
				if !row[i][x] && !col[j][x] && !box[b][x] {
					cloneBoard[i][j] = x
					row[i][x] = true
					col[j][x] = true
					box[b][x] = true
					backtrack(idx + 1)
					if flag {
						return
					}
					cloneBoard[i][j] = 0
					row[i][x] = false
					col[j][x] = false
					box[b][x] = false
				}
			}
		} else {
			backtrack(idx + 1)
		}
	}
	backtrack(0)
	if !flag {
		return nil, errors.New("board cannot be solved")
	}
	return cloneBoard, nil
}
