package game

import (
  "fmt"
  "strings"
)

type Row [8]*Piece

type Board [8]Row

type BoardView interface {
  Get(coord *Coord) *Piece
}

func MakeBoard() *Board {
  board := &Board{}
  initNonPawns(0, White, board)
  initPawns(1, White, board)
  initPawns(6, Black, board)
  initNonPawns(7, Black, board)
  return board
}

func EmptyBoard() *Board {
  return &Board{}
}

func (board *Board) Get(coord *Coord) *Piece {
  if coord != nil && coord.InRange() {
    return board[coord.row][coord.col]
  }
  return nil
}

func (board *Board) Set(coord *Coord, piece *Piece) {
  board[coord.row][coord.col] = piece
}

func (board *Board) MovePiece(move *Move) error {
  from := board.Get(move.from)
  if from == nil {
    return &GameError{fmt.Sprint("No piece at ", move.from)}
  }
  board.Set(move.to, from)
  board.Set(move.from, nil)
  return nil
}


func (board *Board) String() string {
  builder := &strings.Builder{}
  printCols(builder)
  for row := 7; row >= 0; row-- {
    printLine(builder)
    printRow(row, board, builder)
  }
  printLine(builder)
  printCols(builder)
  return builder.String()
}


func initNonPawns(row int, color Color, board *Board) {
  board[row][0] = &Piece{'r', color}
  board[row][1] = &Piece{'n', color}
  board[row][2] = &Piece{'b', color}
  board[row][3] = &Piece{'q', color}
  board[row][4] = &Piece{'k', color}
  board[row][5] = &Piece{'b', color}
  board[row][6] = &Piece{'n', color}
  board[row][7] = &Piece{'r', color}
}

func initPawns(row int, color Color, board *Board) {
  for col := 0; col < 8; col++ {
    board[row][col] = &Piece{'p', color}
  }
}

func printLine(builder *strings.Builder) {
  builder.WriteString("   +")
  for col := 0; col < 8; col++ {
    builder.WriteString("---+")
  }
  builder.WriteByte('\n')
}

func printCols(builder *strings.Builder) {
  builder.WriteString("  ")
  for col := 'a'; col <= 'h'; col++ {
    builder.WriteString(fmt.Sprintf("   %c", col))
  }
  builder.WriteByte('\n')
}

func printRow(row int, board *Board, builder *strings.Builder) {
  builder.WriteString(fmt.Sprintf(" %v |", row + 1))
  for col := 0; col < 8; col++ {
    builder.WriteString(fmt.Sprintf(" %v |", board[row][col]))
  }
  builder.WriteString(fmt.Sprintln(" ", row + 1))
}

func (board *Board) stringKey() string {
  builder := &strings.Builder{}
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      builder.WriteString(board[row][col].String())
    }
  }
  return builder.String()
}
