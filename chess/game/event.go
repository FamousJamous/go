package game

type Event struct {
  moves []*FromTo
  captured *Captured
  isCheck bool
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
