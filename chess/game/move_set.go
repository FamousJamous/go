package game

import "fmt"
import "strconv"
import "strings"

func MakeMoveSet() *Set {
  return MakeSet(func (element interface{}) string {
      return moveKey(element.(*Move))
    })
}

func MakeCoordSet() *Set {
  return MakeSet(func (element interface{}) string {
      return coordKey(element.(*Coord))
    })
}

func moveKey(move *Move) string {
  if move == nil {
    return "nil"
  }
  return strings.Join(
      append(fromToStrts(move.fromTo), capturedStrs(move.captured)...), "|")
}

func coordKey(coord *Coord) string {
  if coord == nil {
    return "nil"
  }
  return strings.Join(coordStrs(coord), "|")
}

func fromToStrts(fromTo *FromTo) []string {
  if fromTo == nil {
    return []string{"nil"}
  }
  return append(coordStrs(fromTo.from), coordStrs(fromTo.to)...)
}

func coordStrs(coord *Coord) []string {
  if coord == nil {
    return []string{"nil"}
  }
  return []string{strconv.Itoa(coord.row), strconv.Itoa(coord.col)}
}

func capturedStrs(captured *Captured) []string {
  if captured == nil {
    return []string{"nil"}
  }
  return append(pieceStrs(captured.piece), coordStrs(captured.coord)...)
}

func pieceStrs(piece *Piece) []string {
  return []string{fmt.Sprint(piece.name, piece.color)}
}
