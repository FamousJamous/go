package game

import (
  "fmt"
  "strings"
)

type Row [8]*Piece

type Rows [8]Row

type Board struct {
  rows *Rows
  whiteKingPos *Coord
  blackKingPos *Coord
  whitePieces map[int]*Piece
  blackPieces map[int]*Piece
  whitePoints int
  blackPoints int
  stringKey []byte
}

type BoardView interface {
  Get(coord *Coord) *Piece
  StringKey() string
  GetPieces(color Color) map[int]*Piece
  GetPoints(color Color) int
}

// Returns kingPos, otherKingPos, the first kingPos will match the given color
func (board *Board) getKingPositions(color Color) (*Coord, *Coord) {
  if color == Black {
    return board.blackKingPos, board.whiteKingPos
  }
  return board.whiteKingPos, board.blackKingPos
}

func (board *Board) GetPieces(color Color) map[int]*Piece {
  if color == Black {
    return board.blackPieces
  }
  return board.whitePieces
}

func (board *Board) GetPoints(color Color) int {
  if color == Black {
    return board.blackPoints
  }
  return board.whitePoints
}

func MakeBoard() *Board {
  board := EmptyBoard()
  initNonPawns(0, White, board)
  initPawns(1, White, board)
  initPawns(6, Black, board)
  initNonPawns(7, Black, board)
  return board
}

func EmptyBoard() *Board {
  board := &Board{
    &Rows{}, nil, nil, make(map[int]*Piece), make(map[int]*Piece), 0, 0,
    make([]byte, 64)}
  for i := 0; i < 64; i++ {
    board.stringKey[i] = ' '
  }
  return board
}

func (board *Board) Get(coord *Coord) *Piece {
  if coord != nil && coord.InRange() {
    return board.rows[coord.row][coord.col]
  }
  return nil
}

func (board *Board) Set(coord *Coord, piece *Piece) {
  if !coord.InRange() {
    panic(fmt.Sprintf("board:\n%v\ncoord out of range %v", board, coord))
  }
  if existing := board.rows[coord.row][coord.col]; existing != nil {
    if existing.color == Black {
      board.blackPoints -= existing.GetPoints()
      delete(board.blackPieces, coord.toKey())
    } else {
      board.whitePoints -= existing.GetPoints()
      delete(board.whitePieces, coord.toKey())
    }
  }
  if piece != nil {
    if piece.color == Black {
      board.blackPieces[coord.toKey()] = piece
      board.blackPoints += piece.GetPoints()
      if piece.name == 'k' {
        board.blackKingPos = coord
      }
    } else {
      board.whitePieces[coord.toKey()] = piece
      board.whitePoints += piece.GetPoints()
      if piece.name == 'k' {
        board.whiteKingPos = coord
      }
    }
  }
  board.stringKey[coord.toArrayIndex()] = piece.toByte()
  board.rows[coord.row][coord.col] = piece
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
  board.Set(&Coord{row, 0}, &Piece{'r', color})
  board.Set(&Coord{row, 1}, &Piece{'n', color})
  board.Set(&Coord{row, 2}, &Piece{'b', color})
  board.Set(&Coord{row, 3}, &Piece{'q', color})
  board.Set(&Coord{row, 4}, &Piece{'k', color})
  board.Set(&Coord{row, 5}, &Piece{'b', color})
  board.Set(&Coord{row, 6}, &Piece{'n', color})
  board.Set(&Coord{row, 7}, &Piece{'r', color})
}

func initPawns(row int, color Color, board *Board) {
  for col := 0; col < 8; col++ {
    board.Set(&Coord{row, col}, &Piece{'p', color})
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
    builder.WriteString(
      fmt.Sprintf(" %v |", board.Get(&Coord{row, col}).String()))
  }
  builder.WriteString(fmt.Sprintln(" ", row + 1))
}

func (board *Board) StringKey() string {
  return string(board.stringKey[:])
}
