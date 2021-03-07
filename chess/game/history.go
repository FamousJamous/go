package game

import "strings"

type History struct {
  moves []*Move
}

func MakeHistory() *History {
  return &History{make([]*Move, 0, 30)}
}

func (history *History) AllMoves() []*Move {
  return history.moves
}

func (history *History) IsEmpty() bool {
  return len(history.moves) == 0
}

func (history *History) AddMove(fromTo *FromTo, capturedPiece *Piece) {
  var captured *Captured = nil
  if capturedPiece != nil {
    // Handle en passant
    captured = &Captured{capturedPiece, fromTo.to}
  }
  history.moves = append(history.moves, &Move{fromTo, captured})
}

func (history *History) GetLastMove() (*Move, error) {
  if len(history.moves) == 0 {
    return nil, &GameError{"No previous moves"}
  }
  return history.moves[len(history.moves) - 1], nil
}

func (history *History) UndoMove(board* Board) error {
  moves := history.moves
  if len(moves) == 0 {
    return &GameError{"No previous moves"}
  }
  last_index := len(moves) - 1
  move := moves[last_index]
  fromTo := move.fromTo
  if isCastle(fromTo, board) {
    if ok := board.MovePiece(getRookCastleMove(fromTo).Reverse());
        ok != nil {
      panic(ok.Error())
    }
  }
  if ok := board.MovePiece(fromTo.Reverse()); ok != nil {
    panic(ok.Error())
  }
  if captured := move.captured; captured != nil {
    board.Set(captured.coord, captured.piece)
  }
  history.moves = moves[:last_index]
  return nil
}

func (history *History) String() string {
  builder := strings.Builder{}
  for i := range history.moves {
    if i != 0 {
      builder.WriteString(" ")
    }
    builder.WriteString(history.moves[i].fromTo.String())
  }
  return builder.String()
}
