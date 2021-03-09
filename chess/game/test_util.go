package game

import "fmt"

func makeMoves(game *Game, moves []string) {
  for _, strMove := range moves {
    move := ParseFromTo(strMove)
    if ok := game.MakeMove(move); ok != nil {
      fmt.Printf("%v\nmove: %v\n", game, move)
      panic(fmt.Sprintf("game:\n%v\nerror: %v", game, ok.Error()))
    }
  }
}
