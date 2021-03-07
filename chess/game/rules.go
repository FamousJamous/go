package game

import "reflect"

type Event struct {
  moves []*FromTo
  captured *Captured
  // TODO isCheck
}

func moveEvent(move *FromTo) (*Event, bool) {
  return &Event{[]*FromTo{move}, nil}, true
}

func captureEvent(move* FromTo, piece *Piece) (*Event, bool) {
  return &Event{[]*FromTo{move}, &Captured{piece, move.to}}, true
}

func moveOrCaptureEvent(move* FromTo, piece *Piece) (*Event, bool) {
  if piece == nil {
    return moveEvent(move)
  }
  return captureEvent(move, piece)
}

func enPassantEvent(move *FromTo, piece *Piece, coord* Coord) (*Event, bool) {
  return &Event{[]*FromTo{move}, &Captured{piece, coord}}, true
}

func castleEvent(kingMove *FromTo, rookMove *FromTo) (*Event, bool) {
  return &Event{[]*FromTo{kingMove, rookMove}, nil}, true
}

func badEvent() (*Event, bool) {
  return nil, false
}

func InterpretMove(move *FromTo, game *Game) (*Event, bool) {
  piece := game.board.Get(move.from)
  if piece == nil || piece.color != game.turn {
    return nil, false
  }
  switch piece.name {
    case 'p': return interpretPawn(move, game)
    case 'n': return interpretKnight(move, game)
    case 'b': return interpretBishop(move, game)
    case 'r': return interpretRook(move, game)
    case 'q': return interpretQueen(move, game)
    case 'k': return interpretKing(move, game)
  }
  panic(piece)
}

func interpretPawn(move *FromTo, game *Game) (*Event, bool) {
  if !isPawnForward(move, game) {
    return badEvent()
  }
  rowDiff, colDiff := move.Diff()
  // 2 forward
  if rowDiff == 2 {
    if colDiff == 0  &&
        game.board.Get(getPawnForward(game.turn, move.from, 1)) == nil &&
        game.board.Get(move.to) == nil && isPawnStart(move.from, game.turn) {
      return moveEvent(move)
    }
    return badEvent()
  }
  if rowDiff != 1 {
    return badEvent()
  }
  // 1 forward
  if colDiff == 0 {
    if game.board.Get(move.to) == nil {
      return moveEvent(move)
    }
    return badEvent()
  }
  if colDiff != 1 {
    return badEvent()
  }
  // capture
  toPiece := game.board.Get(move.to)
  if toPiece != nil && toPiece.color != game.turn {
    return captureEvent(move, toPiece)
  }
  // en passant
  if toPiece != nil {
    return badEvent()
  }
  lastMove := game.history.GetLastMove()
  if lastMove == nil {
    return badEvent()
  }
  lastFromTo := lastMove.fromTo
  lastPiece := game.board.Get(lastFromTo.to)
  if lastPiece.name != 'p' {
    return badEvent()
  }
  toRowDiff, _ := lastFromTo.Diff()
  if toRowDiff == 2 && lastFromTo.to.row == move.from.row &&
      lastFromTo.to.col == move.to.col {
    return enPassantEvent(move, lastPiece, lastFromTo.to)
  }
  return badEvent()
}

func interpretKnight(move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if (rowDiff + colDiff) != 3 || abs(rowDiff - colDiff) != 1 {
    return badEvent()
  }
  toPiece := game.board.Get(move.to)
  if toPiece == nil {
    return moveEvent(move)
  }
  if toPiece.color == game.turn {
    return badEvent()
  }
  return captureEvent(move, toPiece)
}

func interpretBishop(move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff != colDiff {
    return badEvent()
  }
  return interpretStraight(move, game)
}

func interpretRook(move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff != 0 && colDiff != 0 {
    return badEvent()
  }
  return interpretStraight(move, game)
}

func interpretQueen(move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff != colDiff && rowDiff != 0 && colDiff != 0 {
    return badEvent()
  }
  return interpretStraight(move, game)
}

func interpretKing(move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff > 1 || colDiff > 2 {
    return badEvent()
  }
  // Normal move
  toPiece := game.board.Get(move.to)
  if colDiff == 1 || rowDiff == 1 {
    if toPiece != nil && toPiece.color == game.turn {
      return badEvent()
    }
    return moveOrCaptureEvent(move, toPiece)
  }
  // Castle
  if move.from.col != 4 || (move.from.row != 0 && move.from.row != 7) ||
      colDiff != 2 || rowDiff != 0 || toPiece != nil {
    return badEvent()
  }
  rookMove := castleRookMove(move.to)
  if !emptyBetween(MakeFromTo(move.from, rookMove.from), game) ||
      hasMoved(move.from, game) || hasMoved(rookMove.from, game) {
    return badEvent()
  }
  return castleEvent(move, rookMove)
}

func hasMoved(coord *Coord, game *Game) bool {
  for _, move := range game.history.AllMoves() {
    if reflect.DeepEqual(move.fromTo.from, coord) {
      return true
    }
  }
  return false
}

func castleRookMove(kingTo *Coord) *FromTo {
  row := kingTo.row
  if kingTo.col == 2 {
    return MakeFromTo(&Coord{row, 0}, &Coord{row, 3})
  }
  // kingTo.col == 6
  return MakeFromTo(&Coord{row, 7}, &Coord{row, 5})
}

func interpretStraight(move *FromTo, game *Game) (*Event, bool) {
  toPiece := game.board.Get(move.to)
  if toPiece != nil && toPiece.color == game.turn {
    return badEvent()
  }
  if !emptyBetween(move, game) {
    return badEvent()
  }
  return moveOrCaptureEvent(move, toPiece)
}

// Expects diagonal or straight line between
func emptyBetween(move *FromTo, game *Game) bool {
  incRow := getInc(move.from.row, move.to.row)
  incCol := getInc(move.from.col, move.to.col)
  row := move.from.row + incRow
  col := move.from.col + incCol
  for ; row != move.to.row || col != move.to.col; {
    if game.board.Get(&Coord{row, col}) != nil {
      return false
    }
    row += incRow
    col += incCol
  }
  return true
}

func getInc(from int, to int) int {
  if to > from {
    return 1
  } else if to < from {
    return -1
  }
  return 0
}

func isPawnForward(move *FromTo, game *Game) bool {
  if game.turn == Black {
    return move.to.row < move.from.row
  }
  return move.to.row > move.from.row
}

func pawnForward(from *Coord, game *Game) *Coord {
  if game.turn == Black {
    return &Coord{from.row - 1, from.col}
  }
  return &Coord{from.row + 1, from.col}
}
