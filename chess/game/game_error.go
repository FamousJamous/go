package game

type GameError struct {
  message string
}

func (e *GameError) Error() string {
  return e.message
}
