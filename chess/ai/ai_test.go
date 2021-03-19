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

Bug with depth 4
b2b3 d7d5 g1h3 c8h3 g2h3 c7c5 c2c4 h7h6 e2e4 d5c4 f1c4 a7a5 c1a3 f7f5 c4g8 h8g8 a3c5 b7b6 c5e7 d8e7 d1f3 e7e4 f3e4 f5e4 a2a4 e8e7 b1c3
*/

func BenchmarkGetMove(b *testing.B) {
  chessGame := game.MakeGame()
  whitePlayer := MakeAiPlayer(game.White, chessGame, 5)
  b.ResetTimer()

  for i := 0; i < b.N; i++ {
    whitePlayer.GetMove()
  }
}
