package game

import "reflect"

type Event struct {
  moves []*FromTo
  captured *Captured
  isCheck bool
}

func moveEvent(move *FromTo) (*Event, bool) {
  return &Event{[]*FromTo{move}, nil, false}, true
}

func captureEvent(move* FromTo, piece *Piece) (*Event, bool) {
  return &Event{[]*FromTo{move}, &Captured{piece, move.to}, false}, true
}

func moveOrCaptureEvent(move* FromTo, piece *Piece) (*Event, bool) {
  if piece == nil {
    return moveEvent(move)
  }
  return captureEvent(move, piece)
}

func enPassantEvent(move *FromTo, piece *Piece, coord* Coord) (*Event, bool) {
  return &Event{[]*FromTo{move}, &Captured{piece, coord}, false}, true
}

func castleEvent(kingMove *FromTo, rookMove *FromTo) (*Event, bool) {
  return &Event{[]*FromTo{kingMove, rookMove}, nil, false}, true
}

func badEvent() (*Event, bool) {
  return nil, false
}

func InterpretMove(move *FromTo, game *Game) (*Event, bool) {
  piece := game.board.Get(move.from)
  if piece == nil || piece.color != game.turn {
    return nil, false
  }
  event, ok := interpretSimple(piece, move, game)
  if !ok && piece.name == 'k' {
    event, ok = interpretCastle(piece, move, game)
  }
  if !ok {
    return badEvent()
  }
  return checkForCheck(event, game)
}

// Doesn't check for castle or check or player turn.
func interpretSimple(piece *Piece, move *FromTo, game *Game) (*Event, bool) {
  switch piece.name {
    case 'p': return interpretPawn(piece, move, game)
    case 'n': return interpretKnight(piece, move, game)
    case 'b': return interpretBishop(piece, move, game)
    case 'r': return interpretRook(piece, move, game)
    case 'q': return interpretQueen(piece, move, game)
    case 'k': return interpretKing(piece, move, game)
  }
  panic(piece)
}

// Assumes everything is ok before the check
func checkForCheck(event *Event, game *Game) (*Event, bool) {
  // Make move
  if event.captured != nil {
    game.board.Set(event.captured.coord, nil)
  }
  for _, move := range event.moves {
    game.board.MovePiece(move)
  }
  // Check for checks
  kingInCheck, otherKingInCheck := identifyChecks(game)
  // Undo move
  for _, move := range event.moves {
    game.board.MovePiece(move.Reverse())
  }
  if event.captured != nil {
    game.board.Set(event.captured.coord, event.captured.piece)
  }
  // Handle checks
  if kingInCheck {
    return nil, false
  }
  event.isCheck = otherKingInCheck
  return event, true
}

// Returns (current king in check, other king in check)
func identifyChecks(game *Game) (bool, bool) {
  king, otherKing := findKings(game)
  return hasThreat(game.getNextTurn(), king, game),
      hasThreat(game.turn, otherKing, game)
}

func findKings(game *Game) (*Coord, *Coord) {
  var king *Coord
  var otherKing *Coord
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      coord := &Coord{row, col}
      if piece := game.board.Get(coord); piece != nil && piece.name == 'k' {
        if piece.color == game.turn {
          king = coord
        } else {
          otherKing = coord
        }
        if king != nil && otherKing != nil {
          return king, otherKing
        }
      }
    }
  }
  // Non-king games
  return king, otherKing
}

func interpretPawn(piece *Piece, move *FromTo, game *Game) (*Event, bool) {
  if !isPawnForward(move, game) {
    return badEvent()
  }
  rowDiff, colDiff := move.Diff()
  // 2 forward
  if rowDiff == 2 {
    if colDiff == 0  &&
        game.board.Get(getPawnForward(piece.color, move.from, 1)) == nil &&
        game.board.Get(move.to) == nil && isPawnStart(move.from, piece.color) {
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
  if toPiece != nil && toPiece.color != piece.color {
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

func interpretKnight(piece *Piece, move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if (rowDiff + colDiff) != 3 || abs(rowDiff - colDiff) != 1 {
    return badEvent()
  }
  toPiece := game.board.Get(move.to)
  if toPiece == nil {
    return moveEvent(move)
  }
  if toPiece.color == piece.color {
    return badEvent()
  }
  return captureEvent(move, toPiece)
}

func interpretBishop(piece *Piece, move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff != colDiff {
    return badEvent()
  }
  return interpretStraight(piece, move, game)
}

func interpretRook(piece *Piece, move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff != 0 && colDiff != 0 {
    return badEvent()
  }
  return interpretStraight(piece, move, game)
}

func interpretQueen(piece *Piece, move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff != colDiff && rowDiff != 0 && colDiff != 0 {
    return badEvent()
  }
  return interpretStraight(piece, move, game)
}

func interpretKing(piece *Piece, move *FromTo, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff > 1 || colDiff > 1 {
    return badEvent()
  }
  // Normal move
  toPiece := game.board.Get(move.to)
  if colDiff != 1 && rowDiff != 1 ||
      (toPiece != nil && toPiece.color == piece.color) {
    return badEvent()
  }
  return moveOrCaptureEvent(move, toPiece)
}

func interpretCastle(piece *Piece, move *FromTo, game *Game) (*Event, bool) {
  if move.from.col != 4 || (piece.color == Black && move.from.row != 7) ||
      (piece.color == White && move.from.row != 0) {
    return badEvent()
  }
  rowDiff, colDiff := move.Diff()
  if  colDiff != 2 || rowDiff != 0 {
    return badEvent()
  }
  rookMove := castleRookMove(move.to)
  if !emptyBetween(MakeFromTo(move.from, rookMove.from), game) ||
      hasMoved(move.from, game) || hasMoved(rookMove.from, game) ||
      hasThreatsBetween(move.from, rookMove.from, game) ||
      hasThreat(piece.color.Other(), move.from, game) {
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

func interpretStraight(piece *Piece, move *FromTo, game *Game) (*Event, bool) {
  toPiece := game.board.Get(move.to)
  if toPiece != nil && toPiece.color == piece.color {
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

func hasThreatsBetween(from *Coord, to *Coord, game *Game) bool {
  incRow := getInc(from.row, to.row)
  incCol := getInc(from.col, to.col)
  row := from.row + incRow
  col := from.col + incCol
  other := game.board.Get(from).color.Other()
  for ; row != to.row || col != to.col; {
    if hasThreat(other, &Coord{row, col}, game) {
      return true
    }
    row += incRow
    col += incCol
  }
  return false
}

// Returns true if to is threatened by a piece of color. Doesn't account for
// moving into check.
func hasThreat(color Color, to *Coord, game *Game) bool {
  if to == nil {
    return false
  }
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      from := &Coord{row, col}
      if piece := game.board.Get(from); piece != nil && piece.color == color {
        _, ok := interpretSimple(piece, &FromTo{from, to}, game)
        if ok {
          return true
        }
      }
    }
  }
  return false
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
  if game.board.Get(move.from).color == Black {
    return move.to.row < move.from.row
  }
  return move.to.row > move.from.row
}

func pawnForward(from *Coord, game *Game) *Coord {
  if game.board.Get(from).color == Black {
    return &Coord{from.row - 1, from.col}
  }
  return &Coord{from.row + 1, from.col}
}
