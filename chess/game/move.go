package game

import "fmt"

type FromTo struct {
  from *Coord
  to *Coord
}

func (fromTo *FromTo) Diff() (int, int) {
  from := fromTo.from
  to := fromTo.to
  return abs(to.row - from.row), abs(to.col - from.col)
}

func (fromTo *FromTo) InRange() bool {
  return fromTo.from.InRange() && fromTo.to.InRange()
}

func MakeFromTo(from *Coord, to *Coord) *FromTo {
  return &FromTo{from, to}
}

func (fromTo *FromTo) String() string {
  return fmt.Sprintf("%v%v", fromTo.from, fromTo.to)
}

func (fromTo *FromTo) Reverse() *FromTo {
  return MakeFromTo(fromTo.to, fromTo.from)
}

type Captured struct {
  piece *Piece
  coord *Coord
}

func (captured *Captured) String() string {
  return fmt.Sprintf("{ piece: %v, coord: %v }", captured.piece, captured.coord)
}

type Move struct {
  fromTo *FromTo
  captured *Captured
}

func (move *Move) GetFromTo() *FromTo {
  return move.fromTo
}

func (move *Move) String() string {
  return fmt.Sprintf("{ fromTo: %v, captured: %v }", move.fromTo, move.captured)
}

func MakeMove(from *Coord, to *Coord) *Move {
  return &Move{&FromTo{from, to}, nil}
}
