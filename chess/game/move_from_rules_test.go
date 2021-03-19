package game

import (
  "reflect"
  "testing"
)

func checkLegalFailure(
    t *testing.T, game *Game, move *Move, got bool,
    want bool) {
  if got != want {
    t.Errorf("game:\n%v\nmove: %v\nGot: %v\nWant: %v",
             game, move, got, want)
  }
}

func checkEvent(
    t *testing.T, game *Game, move *Move, got *Event, want *Event) {
  if !reflect.DeepEqual(got, want) {
    t.Errorf("game:\n%v\nmove: %v\ngot: %v\nwant: %v",
             game, move, got, want)
  }
}

func TestCheckLegal_NoPiece(t *testing.T) {
  game := MakeGame()
  move := MakeMove(&Coord{2, 2}, &Coord{3, 3})

  _, ok := InterpretMove(move, game)

  checkLegalFailure(t, game, move, ok, false)
}

func checkInterpretMove(
    t *testing.T, game *Game, from *Coord, moves []*Move) {
  legal_to_set := MakeCoordSet()
  for _, move := range moves {
    if _, ok := InterpretMove(move, game); !ok {
      t.Errorf("game:\n%v\nwant ok; but move not ok: %v", game, move)
    }
    legal_to_set.Insert(move.to)
  }
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      to := &Coord{row, col}
      if !legal_to_set.Contains(to) {
        move := MakeMove(from, to)
        _, ok := InterpretMove(move, game)
        if ok {
          t.Errorf("game:\n%v\nwant !ok; move ok: %v", game, move)
        }
      }
    }
  }
}

func checkLegalMovesFrom(
    t *testing.T, game *Game, from *Coord, got []*Move,
    want []*Move) {
  if !reflect.DeepEqual(got, want) {
    t.Errorf("game:\n%v\nfrom: %v\nGot: %v\nWant: %v",
             game, from, got, want)
  }
}

func TestLegalMovesFrom_PawnForward(t *testing.T) {
  game := MakeGame()
  from := ParseCoord("e2")

  moves := LegalMovesFrom(from, game)

  forward1 := ParseMove("e2e3")
  forward2 := ParseMove("e2e4")
  want := []*Move{forward1, forward2}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
}

func TestLegalMovesFrom_PawnCapture(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"e2e4", "d7d5"})
  from := ParseCoord("e4")

  moves := LegalMovesFrom(from, game)

  forward1 := ParseMove("e4e5")
  capture := ParseMove("e4d5")
  want := []*Move{forward1, capture}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
}

func TestLegalMovesFrom_EnPassantLegal(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"e2e4", "d7d5", "e4e5", "f7f5"})
  from := ParseCoord("e5")

  moves := LegalMovesFrom(from, game)
  enPassant := ParseMove("e5f6")
  enPassantEvent, _ := InterpretMove(enPassant, game)

  want := []*Move{ParseMove("e5e6"), enPassant}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
  checkEvent(
      t, game, enPassant, enPassantEvent,
      &Event{[]*Move{enPassant},
             &Captured{&Piece{'p', Black}, ParseCoord("f5")}, nil})
}

func TestLegalMovesFrom_EnPassantHappened(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"e2e4", "d7d5", "e4e5", "f7f5", "e5f6"})

  checkPiece(t, game, "f6", &Piece{'p', White})
  checkPiece(t, game, "f5", nil)
}

func TestLegalMovesFrom_EnPassantUndo(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"e2e4", "d7d5", "e4e5", "f7f5", "e5f6"})

  game.UndoMove()

  checkPiece(t, game, "e5", &Piece{'p', White})
  checkPiece(t, game, "f5", &Piece{'p', Black})
}

func TestLegalMovesFrom_Knight(t *testing.T) {
  game := MakeGame()
  from := ParseCoord("b1")

  moves := LegalMovesFrom(from, game)

  right := ParseMove("b1c3")
  left := ParseMove("b1a3")
  want := []*Move{right, left}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
}

func TestLegalMovesFrom_KnightCapture(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"b1c3", "d7d5"})
  from := ParseCoord("c3")

  moves := LegalMovesFrom(from, game)
  capture := ParseMove("c3d5")
  captureEvent, _ := InterpretMove(capture, game)

  want := []*Move{
      ParseMove("c3e4"), ParseMove("c3a4"), ParseMove("c3d5"),
      ParseMove("c3b5"), ParseMove("c3b1")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
  checkEvent(
      t, game, capture, captureEvent,
      &Event{[]*Move{ParseMove("c3d5")},
             &Captured{&Piece{'p', Black}, ParseCoord("d5")}, nil})
}

func TestLegalMovesFrom_BishopNoMoves(t *testing.T) {
  game := MakeGame()
  from := ParseCoord("c1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
}

func TestLegalMovesFrom_BishopCapture(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"e2e4", "e7e5", "f1c4", "c7c5"})
  from := ParseCoord("c4")

  moves := LegalMovesFrom(from, game)
  capture := ParseMove("c4f7")
  captureEvent, _ := InterpretMove(capture, game)

  want := []*Move{
    // Up right
    ParseMove("c4d5"),
    ParseMove("c4e6"),
    ParseMove("c4f7"),
    // Down right
    ParseMove("c4d3"),
    ParseMove("c4e2"),
    ParseMove("c4f1"),
    // Up left
    ParseMove("c4b5"),
    ParseMove("c4a6"),
    // Down left
    ParseMove("c4b3")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
  checkEvent(t, game, capture, captureEvent,
             &Event{[]*Move{ParseMove("c4f7")},
                    &Captured{&Piece{'p', Black}, ParseCoord("f7")}, nil})
}

func TestLegalMovesFrom_Rook(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"a2a4", "e7e5", "a1a3", "f8d6", "a3e3", "c7c5"})
  from := ParseCoord("e3")

  moves := LegalMovesFrom(from, game)
  capture := ParseMove("e3e5")
  captureEvent, _ := InterpretMove(capture, game)

  want := []*Move{
    // Up
    ParseMove("e3e4"),
    ParseMove("e3e5"),
    // Down
    // Left
    ParseMove("e3d3"),
    ParseMove("e3c3"),
    ParseMove("e3b3"),
    ParseMove("e3a3"),
    // Right
    ParseMove("e3f3"),
    ParseMove("e3g3"),
    ParseMove("e3h3")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
  checkEvent(t, game, capture, captureEvent,
             &Event{[]*Move{ParseMove("e3e5")},
                    &Captured{&Piece{'p', Black}, ParseCoord("e5")}, nil})
}

func TestLegalMovesFrom_Queen(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"e2e4", "e7e5", "d1g4", "e8e7"})
  from := ParseCoord("g4")

  moves := LegalMovesFrom(from, game)
  capture1 := ParseMove("g4d7")
  captureEvent1, _ := InterpretMove(capture1, game)
  capture2 := ParseMove("g4g7")
  captureEvent2, _ := InterpretMove(capture2, game)

  want := []*Move{
    // Up right
    ParseMove("g4h5"),
    // Down right
    ParseMove("g4h3"),
    // Up left
    ParseMove("g4f5"),
    ParseMove("g4e6"),
    ParseMove("g4d7"),
    // Down left
    ParseMove("g4f3"),
    ParseMove("g4e2"),
    ParseMove("g4d1"),
    // Up
    ParseMove("g4g5"),
    ParseMove("g4g6"),
    ParseMove("g4g7"),
    // Down
    ParseMove("g4g3"),
    // Left
    ParseMove("g4f4"),
    // Right
    ParseMove("g4h4")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
  checkEvent(t, game, capture1, captureEvent1,
             &Event{[]*Move{ParseMove("g4d7")},
                    &Captured{&Piece{'p', Black}, ParseCoord("d7")}, nil})
  checkEvent(t, game, capture2, captureEvent2,
             &Event{[]*Move{ParseMove("g4g7")},
                    &Captured{&Piece{'p', Black}, ParseCoord("g7")}, nil})
}

func TestLegalMovesFrom_KingSideCastle(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"e2e4", "e7e5", "f1c4", "d7d5", "g1f3", "c7c5"})
  from := ParseCoord("e1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{ParseMove("e1e2"), ParseMove("e1f1"), ParseMove("e1g1")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
}

func TestLegalMovesFrom_QueenSideCastle(t *testing.T) {
  game := MakeGame()
  MakeMoves(
    game,
    []string{"d2d4", "e7e5", "d1d3", "d7d5", "c1f4", "c7c5", "b1a3", "b7b5"})
  from := ParseCoord("e1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{ParseMove("e1d2"), ParseMove("e1d1"), ParseMove("e1c1")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
}

func TestLegalMovesFrom_ThreatenKingSideCastle(t *testing.T) {
  game := MakeGame()
  MakeMoves(game, []string{"e2e4", "g8f6", "f1c4", "f6g4", "g1f3", "g4e3"})
  from := ParseCoord("e1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{ParseMove("e1e2")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
}

func TestLegalMovesFrom_CheckKingSideCastle(t *testing.T) {
  game := loadGame(
  // abcdehfh
    "rnbqk  r" + // 1
    "pppp ppp" + // 2
    "        " + // 3
    "    R   " + // 4
    "       P" + // 5
    "        " + // 6
    "PPPPPPP " + // 7
    "RNBQKBN ")  // 8

  from := ParseCoord("e1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{ParseMove("e1f1")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
}

func TestLegalMovesFrom_PawnBug(t *testing.T) {
  game := loadGame(
  // abcdehfh
    " rb kbnr" + // 1
    "ppp  ppp" + // 2
    "  nPp   " + // 3
    "   p    " + // 4
    "P       " + // 5
    "        " + // 6
    " P PPPPP" + // 7
    "RNBQKBNR")  // 8
  game.MakeMove(ParseMove("a2a3"))
  from := ParseCoord("d3")

  moves := LegalMovesFrom(from, game)
  capture := ParseMove("d3c2")
  captureEvent, _ := InterpretMove(capture, game)

  want := []*Move{ParseMove("d3d2"), ParseMove("d3c2")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
  checkEvent(t, game, capture, captureEvent,
             &Event{[]*Move{ParseMove("d3c2")},
                    &Captured{&Piece{'p', White}, ParseCoord("c2")}, nil})
}

func TestLegalMovesFrom_CantMoveIntoCheck(t *testing.T) {
  game := loadGame(
  // abcdehfh
    "    R   " + // 1
    "        " + // 2
    "    r   " + // 3
    "    k   " + // 4
    "        " + // 5
    "   Q    " + // 6
    "        " + // 7
    "   K    ")  // 8
  //kingFrom := ParseCoord("e4")
  rookFrom := ParseCoord("e3")
  //rookCapture := ParseMove("e3e1")
  //rookCaptureEvent, _ := InterpretMove(rookCapture, game)

  //kingMoves := LegalMovesFrom(kingFrom, game)
  rookMoves := LegalMovesFrom(rookFrom, game)

  //kingWant := []*Move{ParseMove("e4f5"), ParseMove("e4f3")}
  rookWant := []*Move{ParseMove("e3e2"), ParseMove("e3e1")}
  //checkInterpretMove(t, game, kingFrom, kingWant)
  //checkInterpretMove(t, game, rookFrom, rookWant)
  /*
  checkEvent(t, game, rookCapture, rookCaptureEvent,
             &Event{[]*Move{ParseMove("e3e1")},
                    &Captured{&Piece{'r', Black}, ParseCoord("e1")}, nil})
  checkLegalMovesFrom(t, game, kingFrom, kingMoves, kingWant)
  */
  checkLegalMovesFrom(t, game, rookFrom, rookMoves, rookWant)
  checkGetPieces(t, game.board)
}

func TestLegalMovesFrom_CantCastleAfterMovingKing(t *testing.T) {
  game := loadGame(
  // abcdehfh
    "rnbq   r" + // 1
    "ppppkppp" + // 2
    "     n  " + // 3
    "  b p   " + // 4
    "        " + // 5
    "        " + // 6
    "PPPPPPPP" + // 7
    "RNBQKBNR")  // 8
  MakeMoves(game, []string{"e2e1", "e7e5"})

  from := ParseCoord("e1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{ParseMove("e1e2"), ParseMove("e1f1")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, want)
}

func TestMakeMove_Promo(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "    k   " + // 1
    "        " + // 2
    "        " + // 3
    "        " + // 4
    "        " + // 5
    "        " + // 6
    " p      " + // 7
    "    K   ")  // 8
  from := ParseCoord("b7")

  got := LegalMovesFrom(from, game)

  want := []*Move{
    ParseMove("b7b8q"), ParseMove("b7b8b"), ParseMove("b7b8n"),
    ParseMove("b7b8r")}
  checkLegalMovesFrom(t, game, from, got, want)
  checkInterpretMove(t, game, ParseCoord("h7"), want)
}

func TestMakeMove_UndoPromo(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "    k   " + // 1
    "        " + // 2
    "        " + // 3
    "        " + // 4
    "        " + // 5
    "        " + // 6
    " p      " + // 7
    "    K   ")  // 8
  game.MakeMove(ParseMove("b7b8q"))

  checkPiece(t, game, "b7", nil)
  checkPiece(t, game, "b8", &Piece{'q', White})

  game.UndoMove()

  checkPiece(t, game, "b7", &Piece{'p', White})
  checkPiece(t, game, "b8", nil)
}

func TestMakeMove_UndoPromoCapture(t *testing.T) {
  game := loadGame(
  // abcdefgh
    "    k   " + // 1
    "        " + // 2
    "        " + // 3
    "        " + // 4
    "        " + // 5
    "        " + // 6
    " p      " + // 7
    "R   K   ")  // 8
  game.MakeMove(ParseMove("b7a8q"))

  checkPiece(t, game, "b7", nil)
  checkPiece(t, game, "a8", &Piece{'q', White})

  game.UndoMove()

  checkPiece(t, game, "b7", &Piece{'p', White})
  checkPiece(t, game, "a8", &Piece{'r', Black})
}
