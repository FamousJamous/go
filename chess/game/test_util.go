package game

import (
  "fmt"
  "reflect"
  "testing"
)

func makeMoves(game *Game, moves []string) {
  for _, strMove := range moves {
    move := ParseMove(strMove)
    if ok := game.MakeMove(move); ok != nil {
      fmt.Printf("%v\nmove: %v\n", game, move)
      panic(fmt.Sprintf("game:\n%v\nerror: %v", game, ok.Error()))
    }
  }
}

func checkPiece(t *testing.T, game *Game, coordStr string, want *Piece) {
  got := game.board.Get(ParseCoord(coordStr))
  if !reflect.DeepEqual(got, want) {
    t.Errorf("game:\n%v\ngot: %v has %v\nwant: %v", game, coordStr, got, want)
  }
}
