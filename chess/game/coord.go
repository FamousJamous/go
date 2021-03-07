package game

import "fmt"

type Coord struct {
  row int
  col int
}

func MakeCoord(row int, col int) *Coord {
  return &Coord{row, col}
}

func (coord *Coord) InRange() bool {
  return 0 <= coord.row && coord.row <= 7 && 0 <= coord.col && coord.col <= 7
}

func (coord *Coord) String() string {
  return fmt.Sprintf("%c%c", byte(coord.col + 'a'), byte(coord.row + '1'))
}
