package game

import (
  "fmt"
  "strings"
)

type History struct {
  events []*Event
  toToCount map[int]int
  whiteCaptures []*Piece
  blackCaptures []*Piece
}

func MakeHistory() *History {
  return &History{
    make([]*Event, 0, 30), make(map[int]int), make([]*Piece, 0, 15),
    make([]*Piece, 0, 15)}
}

func (history *History) AllEvents() []*Event {
  return history.events
}

func (history *History) IsEmpty() bool {
  return len(history.events) == 0
}

func (history *History) AddEvent(event *Event) {
  history.events = append(history.events, event)
  for _, move := range event.moves {
    history.toToCount[move.to.toKey()]++
  }
  if event.captured != nil {
    piece := event.captured.piece
    if piece.color == Black {
      history.whiteCaptures = append(history.whiteCaptures, piece)
    } else {
      history.blackCaptures = append(history.blackCaptures, piece)
    }
  }
}

func (history *History) GetLastEvent() *Event {
  if len(history.events) == 0 {
    return nil
  }
  return history.events[len(history.events) - 1]
}

func (history *History) HasMoved(to *Coord) bool {
  return history.toToCount[to.toKey()] > 0
}

func (history *History) UndoMove(board* Board) bool {
  // Get last event
  event := history.GetLastEvent()
  if event == nil {
    return false
  }
  if event.captured != nil {
    piece := event.captured.piece
    if piece.color == Black {
      history.whiteCaptures =
        history.whiteCaptures[:len(history.whiteCaptures) - 1]
    } else {
      history.blackCaptures =
        history.blackCaptures[:len(history.blackCaptures) - 1]
    }
  }
  for _, move := range event.moves {
    history.toToCount[move.to.toKey()]--
  }
  // Undo last event
  event.undo(board)
  // Remove last event
  history.events = history.events[:len(history.events) - 1]
  return true
}

func (history *History) String() string {
  builder := &strings.Builder{}
  for i := range history.events {
    if i != 0 {
      builder.WriteString(" ")
    }
    builder.WriteString(history.events[i].moves[0].String())
  }
  builder.WriteByte('\n')
  printCaptures(builder, White, history.whiteCaptures)
  builder.WriteByte('\n')
  printCaptures(builder, Black, history.blackCaptures)
  return builder.String()
}

func printCaptures(builder *strings.Builder, color Color, captures []*Piece) {
  builder.WriteString(fmt.Sprint(color, " captures:"))
  for _, piece := range captures {
    builder.WriteString(fmt.Sprint(" ", piece))
  }
}
