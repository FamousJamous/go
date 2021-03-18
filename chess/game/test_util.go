package game

import (
  "fmt"
  "reflect"
  "testing"
  "unicode"
)

func MakeMoves(game *Game, moves []string) {
  for _, strMove := range moves {
    move := ParseMove(strMove)
    if ok := game.MakeMove(move); !ok {
      panic(fmt.Sprintf("game:\n%v\nfailed to make move: %v", game, move))
    }
  }
}

func checkPiece(t *testing.T, game *Game, coordStr string, want *Piece) {
  got := game.board.Get(ParseCoord(coordStr))
  if !reflect.DeepEqual(got, want) {
    t.Errorf("game:\n%v\ngot: %v has %v\nwant: %v", game, coordStr, got, want)
  }
}

func loadBoard(boardStr string) *Board {
  board := EmptyBoard()
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      board.Set(&Coord{row, col}, byteToPiece(boardStr[row * 8 + col]))
    }
  }
  return board
}

func loadGame(boardStr string) *Game {
  return LoadGame(White, loadBoard(boardStr), MakeHistory())
}

func byteToPiece(b byte) *Piece {
  if b == ' ' {
    return nil
  }
  name := byte(unicode.ToLower(rune(b)))
  if unicode.IsUpper(rune(b)) {
    return &Piece{name, Black}
  }
  return &Piece{name, White}
}

func checkState(t *testing.T, game *Game, want State) {
  got := game.GetState()
  if got != want {
    t.Errorf("game:\n%v\ngot state: %v\nwant: %v", game, got, want)
  }
}
