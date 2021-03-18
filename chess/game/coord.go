package game

import "fmt"

type Coord struct {
  row int
  col int
}

func keyToCoord(key int) *Coord {
  return &Coord{key / 8, key % 8}
}

func (coord *Coord) toKey() int {
  return coord.row * 8 + coord.col
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
