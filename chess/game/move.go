package game

import "fmt"

type Move struct {
  from *Coord
  to *Coord
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
  return &Move{move.to, move.from}
}
