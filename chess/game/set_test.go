package game

import (
  "reflect"
  "testing"
)

func TestSet(t *testing.T) {
  moves := MakeMoveSet()

  if moves.Size() != 0 {
    t.Errorf("Expected moves to have Size() 0")
  }

  if !moves.IsEmpty() {
    t.Errorf("Expected moves to be IsEmpty()")
  }

  moves.Insert(ParseMove("e2e4"))
  moves.Insert(ParseMove("e2e4"))
  moves.Insert(ParseMove("e2e5"))

  if moves.Size() != 2 {
    t.Errorf("Expected moves to have Size() 2")
  }

  if !reflect.DeepEqual(
      moves.AllElements(),
      []interface{}{ParseMove("e2e4"), ParseMove("e2e5")}) {
    t.Errorf("moves doesn't contain expected elements")
  }

  moves.Erase(ParseMove("e2e4"))
  if moves.Size() != 1 {
    t.Errorf("Expected moves to have 1 element")
  }

  if moves.IsEmpty() {
    t.Errorf("Expected moves to be non-empty")
  }

  moves.Erase(ParseMove("e2e5"))
  if !moves.IsEmpty() {
    t.Errorf("Expected moves to be empty")
  }
}
