package game

import "strings"

type History struct {
  events []*Event
}

func MakeHistory() *History {
  return &History{make([]*Event, 0, 30)}
}

func (history *History) AllEvents() []*Event {
  return history.events
}

func (history *History) IsEmpty() bool {
  return len(history.events) == 0
}

func (history *History) AddEvent(event *Event) {
  history.events = append(history.events, event)
}

func (history *History) GetLastEvent() *Event {
  if len(history.events) == 0 {
    return nil
  }
  return history.events[len(history.events) - 1]
}

func (history *History) UndoMove(board* Board) error {
  // Get last event
  event := history.GetLastEvent()
  if event == nil {
    return &GameError{"No previous moves"}
  }
  // Undo last event
  event.undo(board)
  // Remove last event
  history.events = history.events[:len(history.events) - 1]
  return nil
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
