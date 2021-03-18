package ai

import "jsdu/chess/game"
import "testing"

/*
func TestGetMoveBug(t *testing.T) {
  chessGame := game.MakeGame()
  game.MakeMoves(
      chessGame,
      []string{"b1c3", "a7a5", "a1b1", "a5a4", "b1a1", "a4a3", "a1b1", "a3b2"})
  player := MakeAiPlayer(game.White, chessGame)

  move := player.GetMove()

  if !chessGame.MakeMove(move) {
    t.Errorf("game:\n%v\nAI move: %v is illegal", chessGame, move)
  }
}
*/

func BenchmarkGetMove(b *testing.B) {
  chessGame := game.MakeGame()
  whitePlayer := MakeAiPlayer(game.White, chessGame)
  b.ResetTimer()

  for i := 0; i < b.N; i++ {
    whitePlayer.GetMove()
  }
}
