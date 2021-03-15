package game

import "fmt"
import "reflect"
import "testing"

func testGetPieces(t *testing.T, board *Board) {
  gotBlack := board.GetPieces(Black)
  gotWhite := board.GetPieces(White)

  wantBlack := make(map[*Piece]*Coord)
  wantWhite := make(map[*Piece]*Coord)
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      coord := &Coord{row, col}
      if piece := board.Get(coord);
          piece != nil {
        if piece.color == Black {
          wantBlack[piece] = coord
        } else {
          wantWhite[piece] = coord
        }
      }
    }
  }
  if !reflect.DeepEqual(gotBlack, wantBlack) ||
      !reflect.DeepEqual(gotWhite, wantWhite){
    t.Errorf("board:\n%v\ngot black: %v\nwant black: %v\nblack diff: %v\n" +
        " got white: %v\nwant white: %v\nwhite diff: %v",
        board, gotBlack, wantBlack, diffMaps(gotBlack, wantBlack),
        gotWhite, wantWhite, diffMaps(gotWhite, wantWhite))
  }
}

func diffMaps(lhs map[*Piece]*Coord, rhs map[*Piece]*Coord) string {
  for piece, lhsCoord := range lhs {
    if rhsCoord, ok := rhs[piece];
        !ok || !reflect.DeepEqual(lhsCoord, rhsCoord) {
      return fmt.Sprintf("mismatch - lhs %v, %v, rhs %v", piece, lhsCoord, rhsCoord)
    }
  }
  for piece, rhsCoord := range rhs {
    if lhsCoord, ok := lhs[piece]; !ok {
      return fmt.Sprintf("mismatch - rhs %v, %v, rhs %v", piece, rhsCoord, lhsCoord)
    }
  }
  return "no diff"
}

func TestGetPieces(t *testing.T) {
  board := MakeBoard()

  testGetPieces(t, board)
}

func makeBoardMoves(board *Board, moveStrs []string) {
  for _, moveStr := range moveStrs {
    if moveError := board.MovePiece(ParseMove(moveStr)); moveError != nil {
      panic(moveError.Error())
    }
  }
}

func TestGetPiecesAfterMoves(t *testing.T) {
  board := MakeBoard()
  makeBoardMoves(
      board, []string{"e2e4", "g8f6", "f1c4", "f6g4", "g1f3", "g4e3"})

  testGetPieces(t, board)
}
