package game

import "fmt"

type Move struct {
  from *Coord
  to *Coord
  // Need to clean this up. Makes no sense here. move is used in Event as the
  // rook movement.
  promoteTo byte
}

func (move *Move) apply(board *Board) bool {
  piece := board.Get(move.from)
  if piece == nil {
    return false
  }
  board.Set(move.from, nil)
  if move.promoteTo == 0 {
    board.Set(move.to, piece)
  } else {
    board.Set(move.to, &Piece{move.promoteTo, piece.color})
  }
  return true
}

func (move *Move) undo(board *Board) bool {
  piece := board.Get(move.to)
  if piece == nil {
    panic(fmt.Sprintf("no piece at to %v in undo move", move.to))
  }
  if move.promoteTo == 0 {
    board.Set(move.from, piece)
  } else {
    board.Set(move.from, &Piece{'p', piece.color})
  }
  board.Set(move.to,  nil)
  return true
}

func MakeMove(from *Coord, to *Coord) *Move {
  return &Move{from, to, 0}
}

func MakePromo(from *Coord, to *Coord, promoteTo byte) *Move {
  return &Move{from, to, promoteTo}
}

func (move *Move) Diff() (int, int) {
  from := move.from
  to := move.to
  return abs(to.row - from.row), abs(to.col - from.col)
}

func (move *Move) InRange() bool {
  return move.from.InRange() && move.to.InRange()
}

func (move *Move) String() string {
  return fmt.Sprintf("%v%v", move.from, move.to)
}

func (move *Move) Reverse() *Move {
  return &Move{move.to, move.from, move.promoteTo}
}
