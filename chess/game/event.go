package game

import (
  "fmt"
  "strings"
)

type Event struct {
  moves []*Move
  captured *Captured
  promoteTo *Piece
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
  if event.captured != nil {
    builder.WriteString(fmt.Sprintf(", captured: %v", event.captured))
  }
  builder.WriteByte('}')
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
    if !move.InRange() || !move.apply(board) {
      panic(fmt.Sprintf("board:\n%v\nmove out of range: %v", board, move))
    }
  }
}

func (event *Event) undo(board *Board) {
  for i := len(event.moves) - 1; i >= 0; i-- {
    if move := event.moves[i]; !move.undo(board) {
      panic(fmt.Sprintf("board:\n%v\nmove out of range: %v", board, move))
    }
  }
  if event.captured != nil {
    board.Set(event.captured.coord, event.captured.piece)
  }
}
