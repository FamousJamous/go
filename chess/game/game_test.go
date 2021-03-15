package game

import "testing"

func TestMakeMove_KingSideCastle(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"e2e4", "e7e5", "f1c4", "d7d5", "g1f3", "c7c5"})

  game.MakeMove(ParseMove("e1g1"))

  checkPiece(t, game, "f1", &Piece{'r', White})
  checkPiece(t, game, "g1", &Piece{'k', White})
}

func TestMakeMove_QueenSideCastle(t *testing.T) {
  game := MakeGame()
  makeMoves(
    game,
    []string{"d2d4", "e7e5", "d1d3", "d7d5", "c1f4", "c7c5", "b1a3", "b7b5"})

  game.MakeMove(ParseMove("e1c1"))

  checkPiece(t, game, "c1", &Piece{'k', White})
  checkPiece(t, game, "d1", &Piece{'r', White})
}

func TestUndoMove_KingSideCastle(t *testing.T) {
  game := MakeGame()
  makeMoves(
    game, []string{"e2e4", "e7e5", "f1c4", "d7d5", "g1f3", "c7c5", "e1g1"})

  game.UndoMove()

  checkPiece(t, game, "h1", &Piece{'r', White})
  checkPiece(t, game, "e1", &Piece{'k', White})
}

func TestUndoMove_QueenSideCastle(t *testing.T) {
  game := MakeGame()
  makeMoves(
    game,
    []string{
      "d2d4", "e7e5", "d1d3", "d7d5", "c1f4", "c7c5", "b1a3", "b7b5", "e1c1"})

  game.UndoMove()

  checkPiece(t, game, "e1", &Piece{'k', White})
  checkPiece(t, game, "a1", &Piece{'r', White})
}

func TestUndoMove_EnPassant(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"e2e4", "h7h5", "e4e5", "d7d5", "e5d6"})
  checkPiece(t, game, "d6", &Piece{'p', White})
  checkPiece(t, game, "d5", nil)

  game.UndoMove()

  checkPiece(t, game, "e5", &Piece{'p', White})
  checkPiece(t, game, "d5", &Piece{'p', Black})
}

func TestGetState_NoMovesDraw(t *testing.T) {
  game := loadGame(
  // abcdehfh
    "    k   " + // 1
    "    P   " + // 2
    "    Q   " + // 3
    "        " + // 4
    "        " + // 5
    "        " + // 6
    "        " + // 7
    "        ")  // 8

  checkState(t, game, Draw)
}

func TestGetState_WhiteWins(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "rnb k nr" + // 1
    "pppp ppp" + // 2
    "     q  " + // 3
    "  b p   " + // 4
    "    P   " + // 5
    "  NP    " + // 6
    "PPP  PPP" + // 7
    "R BQKBNR")  // 8
  game.MakeMove(ParseMove("f3f7"))

  checkState(t, game, WhiteWins)
}

func TestGetState_LoneKingDraw(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "        " + // 1
    "        " + // 2
    "        " + // 3
    "   k    " + // 4
    "        " + // 5
    "   K    " + // 6
    "        " + // 7
    "        ")  // 8

  checkState(t, game, Draw)
}

func TestGetState_KingAndBishopDraw(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "        " + // 1
    "        " + // 2
    "  b     " + // 3
    "   k    " + // 4
    "        " + // 5
    "   K    " + // 6
    "  B     " + // 7
    "        ")  // 8

  checkState(t, game, Draw)
}

func TestGetState_KingAndKnightDraw(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "        " + // 1
    "        " + // 2
    "    n   " + // 3
    "   k    " + // 4
    "        " + // 5
    "   K    " + // 6
    "    N   " + // 7
    "        ")  // 8

  checkState(t, game, Draw)
}

func TestGetState_ComboDraw(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "        " + // 1
    "        " + // 2
    "  n     " + // 3
    "   k    " + // 4
    "        " + // 5
    "   K    " + // 6
    "  B     " + // 7
    "        ")  // 8

  checkState(t, game, Draw)
}

func TestGetState_ThreeFoldRepetition(t *testing.T) {
  game := MakeGame()
  makeMoves(
    game,
    []string{"b1c3", "b8c6", "c3b1", "c6b8", "b1c3", "b8c6", "c3b1", "c6b8"})

  checkState(t, game, Draw)
}

func TestGetState_WhiteInCheck(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "        " + // 1
    "        " + // 2
    "  n     " + // 3
    "   k    " + // 4
    "   R    " + // 5
    "   K    " + // 6
    "  B     " + // 7
    "        ")  // 8

  checkState(t, game, WhiteInCheck)
}

func TestGetState_BlackInCheck(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "        " + // 1
    "        " + // 2
    "  n     " + // 3
    "   k    " + // 4
    "   r    " + // 5
    "   K    " + // 6
    "  B     " + // 7
    "        ")  // 8

  checkState(t, game, BlackInCheck)
}

func TestGetState_GameCrashed(t *testing.T) {
  game := MakeGame()
  makeMoves(
      game,
      []string{"b1c3", "a7a6", "c3d5", "a6a5", "a1b1", "a5a4", "b1a1", "a4a3",
      "b2a3", "b7b6", "a1b1", "c7c6", "d5b6", "a8a7", "b6c8", "a7a8", "c8e7",
      "d8e7", "e2e3", "c6c5", "c1b2", "e7e4", "b1a1", "e4f5", "a1b1", "f5g6",
      "b1a1", "f7f6", "a1b1", "f6f5", "b1a1", "h7h6", "a1b1", "h6h5", "b1a1",
      "h5h4", "a1b1", "h4h3", "g1h3", "c5c4", "b1a1", "c4c3", "d2c3", "g6h7",
      "a1b1", "g7g6", "d1d4", "b8c6", "h3g5", "c6d4", "g5h7", "d4c2", "e1d2",
      "c2e3", "h7f8", "e3f1", "b1f1", "e8f8", "f1e1", "f5f4", "e1e2", "g6g5",
      "b2c1", "g5g4", "d2d3", "f4f3", "g2f3", "g4f3", "e2e3", "g8f6", "e3f3",
      "a8a6", "d3c2", "a6b6", "c1d2", "b6a6", "c3c4", "a6b6", "c2c3", "b6a6",
      "a3a4", "a6b6", "a2a3", "b6a6", "h1g1", "a6b6", "g1g2", "b6a6", "d2e3",
      "a6c6", "h2h3", "h8h7", "g2g3", "h7h8", "c3b4", "h8h7", "e3d4", "c6a6",
      "f3f4", "a6c6", "f2f3", "c6a6", "h3h4", "a6c6", "g3g4", "c6a6", "a4a5",
      "a6c6", "a3a4", "c6a6", "b4b5", "a6c6", "c4c5", "h7h8", "d4e5", "f8e7",
      "e5f6", "c6f6", "f4f6", "e7f6", "f3f4", "f6f5", "g4g5", "f5e4", "c5c6",
      "d7d6", "g5g6", "d6d5", "f4f5", "e4e5", "h4h5", "d5d4", "a5a6", "d4d3",
      "a4a5", "d3d2"})

  game.GetState()
}

func BenchmarkGetState(b *testing.B) {
  game := MakeGame()
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    game.GetState()
  }
}

func BenchmarkGetStateAfterSomeMoves(b *testing.B) {
  game := loadGame(
  // abcdefgh
    "rnb k nr" + // 1
    "pppp ppp" + // 2
    "     q  " + // 3
    "  b p   " + // 4
    "    P   " + // 5
    "  NP    " + // 6
    "PPP  PPP" + // 7
    "R BQKBNR")  // 8
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    game.GetState()
  }
}

func BenchmarkGetAllMoves(b *testing.B) {
  game := loadGame(
  // abcdefgh
    "rnb k nr" + // 1
    "pppp ppp" + // 2
    "     q  " + // 3
    "  b p   " + // 4
    "    P   " + // 5
    "  NP    " + // 6
    "PPP  PPP" + // 7
    "R BQKBNR")  // 8
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    game.GetAllMoves()
  }
}
