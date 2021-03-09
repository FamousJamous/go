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

// no moves draw
// insufficient materials draw
// check mate
// complex check mate
