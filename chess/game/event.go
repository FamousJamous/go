package game

import (
  "fmt"
  "strings"
)

type Event struct {
  moves []*Move
  captured *Captured
  isCheck bool
}

func (event *Event) String() string {
  if event == nil {
    return "nil"
  }
  builder := &strings.Builder{}
  builder.WriteString("{moves: ")
  for i, move := range event.moves {
    if i != 0 {
      builder.WriteByte(' ')
    }
    builder.WriteString(move.String())
  }
  builder.WriteString(
      fmt.Sprintf(", captured: %v, isCheck: %v}", event.captured,
                  event.isCheck))
  return builder.String()
}

type Captured struct {
  piece *Piece
  coord *Coord
}

func (captured *Captured) String() string {
  if captured == nil {
    return "nil"
  }
  return fmt.Sprintf("{ piece: %v, coord: %v }", captured.piece, captured.coord)
}

func (event *Event) apply(board *Board) {
  if event.captured != nil {
    board.Set(event.captured.coord, nil)
  }
  for _, move := range event.moves {
    board.MovePiece(move)
  }
}

func (event *Event) undo(board *Board) {
  for _, move := range event.moves {
    board.MovePiece(move.Reverse())
  }
  if event.captured != nil {
    board.Set(event.captured.coord, event.captured.piece)
  }
}
