package game

import "strings"

type History struct {
  events []*Event
  from_to_count map[int]int
}

func MakeHistory() *History {
  return &History{make([]*Event, 0, 30), make(map[int]int)}
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
    history.from_to_count[move.from.toKey()]++
  }
}

func (history *History) GetLastEvent() *Event {
  if len(history.events) == 0 {
    return nil
  }
  return history.events[len(history.events) - 1]
}

func (history *History) HasMoved(from *Coord) bool {
  return history.from_to_count[from.toKey()] > 0
}

func (history *History) UndoMove(board* Board) bool {
  // Get last event
  event := history.GetLastEvent()
  if event == nil {
    return false
  }
  for _, move := range event.moves {
    history.from_to_count[move.from.toKey()]--
  }
  // Undo last event
  event.undo(board)
  // Remove last event
  history.events = history.events[:len(history.events) - 1]
  return true
}

func (history *History) String() string {
  builder := strings.Builder{}
  for i := range history.events {
    if i != 0 {
      builder.WriteString(" ")
    }
    builder.WriteString(history.events[i].moves[0].String())
  }
  return builder.String()
}
