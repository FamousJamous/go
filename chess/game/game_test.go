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
