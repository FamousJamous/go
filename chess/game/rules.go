package game

func moveEvent(move *Move) (*Event, bool) {
  return &Event{[]*Move{move}, nil, nil}, true
}

func promoEvent(event *Event, pawn *Piece) (*Event, bool) {
  move := event.moves[0]
  if !validPromoteTo(move) {
    return badEvent()
  }
  event.promoteTo = &Piece{move.promoteTo, pawn.color}
  return event, true
}

func captureEvent(move* Move, capturedPiece *Piece) (*Event, bool) {
  return &Event{[]*Move{move}, &Captured{capturedPiece, move.to}, nil}, true
}

func moveOrCaptureEvent(move* Move, piece *Piece) (*Event, bool) {
  if piece == nil {
    return moveEvent(move)
  }
  return captureEvent(move, piece)
}

func enPassantEvent(move *Move, piece *Piece, coord* Coord) (*Event, bool) {
  return &Event{[]*Move{move}, &Captured{piece, coord}, nil}, true
}

func castleEvent(kingMove *Move, rookMove *Move) (*Event, bool) {
  return &Event{[]*Move{kingMove, rookMove}, nil, nil}, true
}

func badEvent() (*Event, bool) {
  return nil, false
}

func InterpretMove(move *Move, game *Game) (*Event, bool) {
  piece := game.board.Get(move.from)
  if piece == nil || piece.color != game.turn {
    return badEvent()
  }
  event, ok := interpretSimple(piece, move, game)
  if !ok && piece.name == 'k' {
    event, ok = interpretCastle(piece, move, game)
  }
  if !ok {
    return badEvent()
  }
  if piece.name == 'p' && isPawnPromoRow(piece, move.to) {
    event, ok = promoEvent(event, piece)
    if !ok {
      return badEvent()
    }
  }
  return checkForCheck(event, game)
}

// Doesn't check for castle, check, player turn, or promo.
func interpretSimple(piece *Piece, move *Move, game *Game) (*Event, bool) {
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
  event.apply(game.board)
  // Check for checks
  kingInCheck, _ := identifyChecks(game)
  // Undo move
  event.undo(game.board)
  // Handle checks
  if kingInCheck {
    return badEvent()
  }
  return event, true
}

// Returns (current king in check, other king in check)
func identifyChecks(game *Game) (bool, bool) {
  king, otherKing := game.board.getKingPositions(game.turn)
  return hasThreat(game.getNextTurn(), king, game),
      hasThreat(game.turn, otherKing, game)
}

func interpretPawn(piece *Piece, move *Move, game *Game) (*Event, bool) {
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
    if game.board.Get(move.to) != nil {
      return badEvent()
    }
    // Normal move
    return moveEvent(move)
  }
  // capture
  if colDiff != 1 {
    return badEvent()
  }
  toPiece := game.board.Get(move.to)
  if toPiece != nil && toPiece.color != piece.color {
    // Capture
    return captureEvent(move, toPiece)
  }
  // en passant
  if toPiece != nil {
    return badEvent()
  }
  lastEvent := game.history.GetLastEvent()
  if lastEvent == nil {
    return badEvent()
  }
  lastMove := lastEvent.moves[0]
  lastPiece := game.board.Get(lastMove.to)
  if lastPiece.name != 'p' {
    return badEvent()
  }
  toRowDiff, _ := lastMove.Diff()
  if toRowDiff == 2 && lastMove.to.row == move.from.row &&
      lastMove.to.col == move.to.col {
    return enPassantEvent(move, lastPiece, lastMove.to)
  }
  return badEvent()
}

// Assume move is legal in all ways except check.
func validPromoteTo(move *Move) bool {
  switch move.promoteTo {
    case 'q': return true
    case 'b': return true
    case 'n': return true
    case 'r': return true
    default: return false
  }
}

// Assumes pawn is a pawn
func isPawnPromoRow(pawn *Piece, coord *Coord) bool {
  if pawn.color == Black {
    return coord.row == 0
  }
  return coord.row == 7
}

func interpretKnight(piece *Piece, move *Move, game *Game) (*Event, bool) {
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

func interpretBishop(piece *Piece, move *Move, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff != colDiff {
    return badEvent()
  }
  return interpretStraight(piece, move, game)
}

func interpretRook(piece *Piece, move *Move, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff != 0 && colDiff != 0 {
    return badEvent()
  }
  return interpretStraight(piece, move, game)
}

func interpretQueen(piece *Piece, move *Move, game *Game) (*Event, bool) {
  rowDiff, colDiff := move.Diff()
  if rowDiff != colDiff && rowDiff != 0 && colDiff != 0 {
    return badEvent()
  }
  return interpretStraight(piece, move, game)
}

func interpretKing(piece *Piece, move *Move, game *Game) (*Event, bool) {
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

func interpretCastle(piece *Piece, move *Move, game *Game) (*Event, bool) {
  if move.from.col != 4 || (piece.color == Black && move.from.row != 7) ||
      (piece.color == White && move.from.row != 0) {
    return badEvent()
  }
  rowDiff, colDiff := move.Diff()
  if  colDiff != 2 || rowDiff != 0 {
    return badEvent()
  }
  rookMove := castleRookMove(move.to)
  if !hasRook(rookMove.from, piece.color, game) ||
      !emptyBetween(move.from, rookMove.from, game) ||
      hasMoved(move.from, game) || hasMoved(rookMove.from, game) ||
      hasThreatsBetween(move.from, rookMove.from, game) ||
      hasThreat(piece.color.Other(), move.from, game) {
    return badEvent()
  }
  return castleEvent(move, rookMove)
}

func hasRook(coord *Coord, color Color, game *Game) bool {
  piece := game.board.Get(coord)
  return piece != nil && piece.name == 'r' && piece.color == color
}

func hasMoved(coord *Coord, game *Game) bool {
  return game.history.HasMoved(coord)
}

func castleRookMove(kingTo *Coord) *Move {
  row := kingTo.row
  if kingTo.col == 2 {
    return MakeMove(&Coord{row, 0}, &Coord{row, 3})
  }
  // kingTo.col == 6
  return MakeMove(&Coord{row, 7}, &Coord{row, 5})
}

func interpretStraight(piece *Piece, move *Move, game *Game) (*Event, bool) {
  toPiece := game.board.Get(move.to)
  if toPiece != nil && toPiece.color == piece.color {
    return badEvent()
  }
  if !emptyBetween(move.from, move.to, game) {
    return badEvent()
  }
  return moveOrCaptureEvent(move, toPiece)
}

// Expects diagonal or straight line between
func emptyBetween(from *Coord, to *Coord, game *Game) bool {
  incRow := getInc(from.row, to.row)
  incCol := getInc(from.col, to.col)
  row := from.row + incRow
  col := from.col + incCol
  for ; row != to.row || col != to.col; {
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
  for from, piece := range game.board.GetPieces(color) {
    if _, ok := interpretSimple(piece, MakeMove(keyToCoord(from), to), game);
        ok {
      return true
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

func isPawnForward(move *Move, game *Game) bool {
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
