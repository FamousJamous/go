package game

type Event struct {
  moves []*Move
  captured *Captured
  isCheck bool
}

func InterpretMove(fromTo *FromTo, game *Game) (*Event, bool) {
  return nil, false
}
