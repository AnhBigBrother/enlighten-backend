package services

import (
	"context"

	"github.com/AnhBigBrother/enlighten-backend/cfg"
	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/AnhBigBrother/enlighten-backend/internal/pb"
	"github.com/AnhBigBrother/enlighten-backend/internal/pkg/sudoku"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameServer struct {
	pb.UnimplementedGameServer
	DBQueries    *database.Queries
	DBConnection *pgxpool.Pool
}

func NewGameServer() *GameServer {
	return &GameServer{
		DBQueries:    cfg.DBQueries,
		DBConnection: cfg.DBConnection,
	}
}

func (server *GameServer) SolveSudoku(ctx context.Context, req *pb.SolveSudokuRequest) (*pb.SolveSudokuResponse, error) {
	board := [][]int32{}
	board = append(board, req.GetLine1())
	board = append(board, req.GetLine2())
	board = append(board, req.GetLine3())
	board = append(board, req.GetLine4())
	board = append(board, req.GetLine5())
	board = append(board, req.GetLine6())
	board = append(board, req.GetLine7())
	board = append(board, req.GetLine8())
	board = append(board, req.GetLine9())
	result, err := sudoku.Solve(board)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &pb.SolveSudokuResponse{
		Line1: result[0],
		Line2: result[1],
		Line3: result[2],
		Line4: result[3],
		Line5: result[4],
		Line6: result[5],
		Line7: result[6],
		Line8: result[7],
		Line9: result[8],
	}, nil
}

func (server *GameServer) GenerateSudoku(ctx context.Context, req *pb.GenerateSudokuRequest) (*pb.GenerateSudokuResponse, error) {
	board := sudoku.Gen(int32(req.GetHide()))
	return &pb.GenerateSudokuResponse{
		Line1: board[0],
		Line2: board[1],
		Line3: board[2],
		Line4: board[3],
		Line5: board[4],
		Line6: board[5],
		Line7: board[6],
		Line8: board[7],
		Line9: board[8],
	}, nil
}

func (server *GameServer) CheckSudokuSolvable(ctx context.Context, req *pb.CheckSudokuSolvableRequest) (*pb.CheckSudokuSolvableResponse, error) {
	board := [][]int32{}
	board = append(board, req.GetLine1())
	board = append(board, req.GetLine2())
	board = append(board, req.GetLine3())
	board = append(board, req.GetLine4())
	board = append(board, req.GetLine5())
	board = append(board, req.GetLine6())
	board = append(board, req.GetLine7())
	board = append(board, req.GetLine8())
	board = append(board, req.GetLine9())
	_, err := sudoku.Solve(board)
	if err != nil {
		return &pb.CheckSudokuSolvableResponse{
			Solvable: false,
		}, nil
	}
	return &pb.CheckSudokuSolvableResponse{
		Solvable: true,
	}, nil
}
