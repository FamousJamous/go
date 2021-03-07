package game

import (
  "reflect"
  "testing"
  "unicode"
)

func checkLegalFailure(
    t *testing.T, game *Game, fromTo *FromTo, got error,
    want error) {
  if got != nil && got.Error() != want.Error() {
    t.Errorf("game:\n%v\nfromTo: %v\nGot: %v\nWant: %v",
             game, fromTo, got, want)
  }
}

func TestCheckLegal_NoPiece(t *testing.T) {
  game := MakeGame()
  fromTo := &FromTo{&Coord{2, 2}, &Coord{3, 3}}

  cOk := CheckLegal(fromTo, game)

  gameError := &GameError{"Illegal move"}
  checkLegalFailure(t, game, fromTo, cOk, gameError)
}

func checkInterpretMove(
    t *testing.T, game *Game, from *Coord, moves []*Move) {
  legal_to_set := MakeCoordSet()
  for _, move := range moves {
    if _, ok := InterpretMove(move.fromTo, game); !ok {
      t.Errorf("game:\n%v\nmove not ok: %v", game, move)
    }
    legal_to_set.Insert(move.fromTo.to)
  }
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      to := &Coord{row, col}
      if !legal_to_set.Contains(to) {
        move := MakeFromTo(from, to)
        _, ok := InterpretMove(move, game)
        if ok {
          t.Errorf("game:\n%v\nmove ok: %v", game, move)
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
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_PawnCapture(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"e2e4", "d7d5"})
  from := ParseCoord("e4")

  moves := LegalMovesFrom(from, game)

  forward1 := ParseMove("e4e5")
  capture := &Move{
      ParseFromTo("e4d5"), &Captured{&Piece{'p', Black}, ParseCoord("d5")}}
  want := []*Move{forward1, capture}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_EnPassant(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"e2e4", "d7d5", "e4e5", "f7f5"})
  from := ParseCoord("e5")

  moves := LegalMovesFrom(from, game)

  forward1 := ParseMove("e5e6")
  enPassant := &Move{
    ParseFromTo("e5f6"), &Captured{&Piece{'p', Black}, ParseCoord("f5")}}
  want := []*Move{forward1, enPassant}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_Knight(t *testing.T) {
  game := MakeGame()
  from := ParseCoord("b1")

  moves := LegalMovesFrom(from, game)

  right := ParseMove("b1c3")
  left := ParseMove("b1a3")
  want := []*Move{right, left}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_KnightCapture(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"b1c3", "d7d5"})
  from := ParseCoord("c3")

  moves := LegalMovesFrom(from, game)

  capture := &Move{
      ParseFromTo("c3d5"), &Captured{&Piece{'p', Black}, ParseCoord("d5")}}
  want := []*Move{
      ParseMove("c3e4"), ParseMove("c3a4"), capture, ParseMove("c3b5"),
      ParseMove("c3b1")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_BishopNoMoves(t *testing.T) {
  game := MakeGame()
  from := ParseCoord("c1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_BishopCapture(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"e2e4", "e7e5", "f1c4", "c7c5"})
  from := ParseCoord("c4")

  moves := LegalMovesFrom(from, game)

  want := []*Move{
    // Up right
    ParseMove("c4d5"),
    ParseMove("c4e6"),
    &Move{ParseFromTo("c4f7"), &Captured{&Piece{'p', Black}, ParseCoord("f7")}},
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
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_Rook(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"a2a4", "e7e5", "a1a3", "f8d6", "a3e3", "c7c5"})
  from := ParseCoord("e3")

  moves := LegalMovesFrom(from, game)

  want := []*Move{
    // Up
    ParseMove("e3e4"),
    &Move{ParseFromTo("e3e5"), &Captured{&Piece{'p', Black}, ParseCoord("e5")}},
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
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_Queen(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"e2e4", "e7e5", "d1g4", "e8e7"})
  from := ParseCoord("g4")

  moves := LegalMovesFrom(from, game)

  want := []*Move{
    // Up right
    ParseMove("g4h5"),
    // Down right
    ParseMove("g4h3"),
    // Up left
    ParseMove("g4f5"),
    ParseMove("g4e6"),
    &Move{ParseFromTo("g4d7"), &Captured{&Piece{'p', Black}, ParseCoord("d7")}},
    // Down left
    ParseMove("g4f3"),
    ParseMove("g4e2"),
    ParseMove("g4d1"),
    // Up
    ParseMove("g4g5"),
    ParseMove("g4g6"),
    &Move{ParseFromTo("g4g7"), &Captured{&Piece{'p', Black}, ParseCoord("g7")}},
    // Down
    ParseMove("g4g3"),
    // Left
    ParseMove("g4f4"),
    // Right
    ParseMove("g4h4")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_KingSideCastle(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"e2e4", "e7e5", "f1c4", "d7d5", "g1f3", "c7c5"})
  from := ParseCoord("e1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{ParseMove("e1e2"), ParseMove("e1f1"), ParseMove("e1g1")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_QueenSideCastle(t *testing.T) {
  game := MakeGame()
  makeMoves(
    game,
    []string{"d2d4", "e7e5", "d1d3", "d7d5", "c1f4", "c7c5", "b1a3", "b7b5"})
  from := ParseCoord("e1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{ParseMove("e1d2"), ParseMove("e1d1"), ParseMove("e1c1")}
  checkLegalMovesFrom(t, game, from, moves, want)
  checkInterpretMove(t, game, from, moves)
}

func TestLegalMovesFrom_BlockKingSideCastle(t *testing.T) {
  game := MakeGame()
  makeMoves(game, []string{"e2e4", "g8f6", "f1c4", "f6g4", "g1f3", "g4e3"})
  from := ParseCoord("e1")

  moves := LegalMovesFrom(from, game)

  want := []*Move{ParseMove("e1e2"), ParseMove("e1f1")}
  checkLegalMovesFrom(t, game, from, moves, want)
}

func loadGame(boardStr string) *Game {
  board := EmptyBoard()
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      board.Set(&Coord{row, col}, byteToPiece(boardStr[row * 8 + col]))
    }
  }
  return LoadGame(White, board, MakeHistory())
}

func byteToPiece(b byte) *Piece {
  if b == ' ' {
    return nil
  }
  if unicode.IsUpper(rune(b)) {
    return &Piece{b, Black}
  }
  return &Piece{b, White}
}

/*
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
    "        ")  // 8
  kingFrom := ParseCoord("e4")
  rookFrom := ParseCoord("e3")

  kingMoves := LegalMovesFrom(kingFrom, game)
  rookMoves := LegalMovesFrom(rookFrom, game)

  kingWant := []*Move{ParseMove("e4f5"), ParseMove("e4f3")}
  rookWant := []*Move{
    ParseMove("e3e2"),
    &Move{
      ParseFromTo("e3e1"), &Captured{&Piece{'r', Black}, ParseCoord("e1")}}}
  checkLegalMovesFrom(t, game, kingFrom, kingMoves, kingWant)
  checkLegalMovesFrom(t, game, rookFrom, rookMoves, rookWant)
}
*/

// check
