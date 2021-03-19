package game

import (
  "fmt"
  //"math/rand"
  "reflect"
  "testing"
)

func checkGetPieces(t *testing.T, board *Board) {
  gotBlack := board.GetPieces(Black)
  gotWhite := board.GetPieces(White)

  wantBlack := make(map[int]*Piece)
  wantWhite := make(map[int]*Piece)
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      coord := &Coord{row, col}
      if piece := board.Get(coord);
          piece != nil {
        if piece.color == Black {
          wantBlack[coord.toKey()] = piece
        } else {
          wantWhite[coord.toKey()] = piece
        }
      }
    }
  }
  if !reflect.DeepEqual(gotBlack, wantBlack) ||
      !reflect.DeepEqual(gotWhite, wantWhite){
    t.Errorf("board:\n%v\ngot black: %v\nwant black: %v\nblack diff: %v\n" +
        "got white: %v\nwant white: %v\nwhite diff: %v",
        board, gotBlack, wantBlack, diffMaps(gotBlack, wantBlack),
        gotWhite, wantWhite, diffMaps(gotWhite, wantWhite))
  }
}

func diffMaps(got map[int]*Piece, want map[int]*Piece) string {
  for coord, gotPiece := range got {
    if wantPiece, ok := want[coord];
        !ok || !reflect.DeepEqual(gotPiece, wantPiece) {
      return fmt.Sprintf("mismatch - got %v, '%v', want '%v'",
          keyToCoord(coord), gotPiece, wantPiece)
    }
  }
  for coord, wantPiece := range want {
    if gotPiece, ok := got[coord]; !ok || !reflect.DeepEqual(gotPiece, wantPiece) {
      return fmt.Sprintf("mismatch - want %v, '%v', got '%v'",
          keyToCoord(coord), wantPiece, gotPiece)
    }
  }
  return "no diff"
}

func TestGetPieces(t *testing.T) {
  board := MakeBoard()

  checkGetPieces(t, board)
}

func makeBoardMoves(board *Board, moveStrs []string) {
  for _, moveStr := range moveStrs {
    move := ParseMove(moveStr)
    if !move.apply(board) {
      panic(fmt.Sprintf("%v\nfailed to apply move: %v", board, move))
    }
  }
}

func TestGetPiecesAfterMoves(t *testing.T) {
  board := MakeBoard()
  makeBoardMoves(
      board, []string{"e2e4", "g8f6", "f1c4", "f6g4", "g1f3", "g4e3"})

  checkGetPieces(t, board)
}

func TestGetPiecesBug(t *testing.T) {
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
  game.GetState()

  checkGetPieces(t, game.board)
}

func TestGetPiecesOtherBug(t *testing.T) {
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
  LegalMovesFrom(ParseCoord("e4"), game)
  LegalMovesFrom(ParseCoord("e3"), game)
  InterpretMove(ParseMove("e3e1"), game)

  checkGetPieces(t, game.board)
}

func TestGetPoints(t *testing.T) {
  board := MakeBoard()

  checkGetPoints(t, board)
}

func TestGetPoints_AfterMoves(t *testing.T) {
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

  checkGetPoints(t, game.board)
}

func checkGetPointsColor(t *testing.T, board *Board, color Color) {
  got := board.GetPoints(color)

  want := 0
  for _, piece := range board.GetPieces(color) {
    want += piece.GetPoints()
  }
  if want != got {
    t.Errorf("board:\n%v\ngot: %v\nwant: %v", board, got, want)
  }
}

func checkGetPoints(t *testing.T, board *Board) {
  checkGetPointsColor(t, board, Black)
  checkGetPointsColor(t, board, White)
}

/*
func TestGetPoints_RandomMovesAndUndo(t *testing.T) {
  game := MakeGame()

  for i := 0; i < 20; i++ {
    moves := game.GetAllMoves()
    move := moves[rand.Intn(len(moves))]
    if !game.MakeMove(move) {
      panic(fmt.Sprintf("game:\n%v\nmove: %v", game, move))
    }
    checkGetPoints(t, game.board)
  }
}
*/
