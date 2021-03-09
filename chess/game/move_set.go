package game

import "fmt"
import "strconv"
import "strings"

func MakeCoordSet() *Set {
  return MakeSet(func (element interface{}) string {
      return coordKey(element.(*Coord))
    })
}

func MakeMoveSet() *Set {
  return MakeSet(func (element interface{}) string {
      return moveKey(element.(*Move))
    })
}

func coordKey(coord *Coord) string {
  if coord == nil {
    return "nil"
  }
  return strings.Join(coordStrs(coord), "|")
}

func moveKey(move *Move) string {
  if move == nil {
    return "nil"
  }
  return strings.Join(moveStrs(move), "|")
}

func moveStrs(move *Move) []string {
  if move == nil {
    return []string{"nil"}
  }
  return append(coordStrs(move.from), coordStrs(move.to)...)
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
