package game

func LegalMovesFrom(from *Coord, game *Game) []*Move {
  piece := game.board.Get(from)
  if piece == nil || piece.color != game.turn {
    return []*Move{}
  }
  return legalMovesForPiece(piece, from, game)
}

func legalMovesForPiece(piece *Piece, from *Coord, game *Game) []*Move {
  moves := simpleMovesForPiece(piece, from, game)
  if piece.name == 'k' {
    moves = append(moves, castleKingMovesFrom(from, game)...)
  }
  return moves
}

// Omits castling
func simpleMovesForPiece(piece *Piece, from *Coord, game *Game) []*Move {
  switch piece.name {
    case 'p': return legalPawnMovesFrom(from, game)
    case 'n': return legalKnightMovesFrom(from, game)
    case 'b': return legalBishopMovesFrom(from, game)
    case 'r': return legalRookMovesFrom(from, game)
    case 'q': return legalQueenMovesFrom(from, game)
    case 'k': return simpleKingMovesFrom(from, game)
    default: return make([]*Move, 0)
  }
}

func appendPawnIf(
  toIsEmpty bool, pawn *Piece, from *Coord, to *Coord, moves []*Move,
  game *Game,
) []*Move {
  if (game.board.Get(to) == nil) != toIsEmpty {
    return moves
  }
  if isPawnPromoRow(pawn, to) {
    promo := MakePromo(from, to, 'q')
    if _, ok := InterpretMove(promo, game); !ok {
      return moves
    }
    moves = append(moves, promo)
    // Assume the rest are also ok.
    moves = append(moves, MakePromo(from, to, 'b'))
    moves = append(moves, MakePromo(from, to, 'n'))
    moves = append(moves, MakePromo(from, to, 'r'))
    return moves
  }
  move := MakeMove(from, to)
  if _, ok := InterpretMove(move, game); !ok {
    return moves
  }
  return append(moves, move)
}

func legalPawnMovesFrom(from *Coord, game *Game) []*Move {
  pawn := game.board.Get(from)
  moves := make([]*Move, 0, 3)
  // Forward 1
  forward1 := getPawnForward(pawn.color, from, 1)
  moves = appendPawnIf(/*toIsEmpty=*/true, pawn, from, forward1, moves, game)
  // Forward 2
  if game.board.Get(forward1) == nil && isPawnStart(from, pawn.color) {
    moves = appendIfEmpty(
      from, getPawnForward(pawn.color, from, 2), game, moves)
  }
  // Left capture
  moves = appendPawnIf(/*toIsEmpty=*/false, pawn, from,
                       pawnDiagonal(pawn.color, from, -1), moves, game)
  // Right capture
  moves = appendPawnIf(/*toIsEmpty=*/false, pawn, from,
                       pawnDiagonal(pawn.color, from, 1), moves, game)
  // En passant
  moves = appendIfEnPassant(from, game, moves)
  return moves
}

func legalKnightMovesFrom(from *Coord, game *Game) []*Move {
  moves := make([]*Move, 0, 8)
  moves = appendIfEmptyOrCapture(from, &Coord{from.row + 1, from.col + 2}, game, moves)
  moves = appendIfEmptyOrCapture(from, &Coord{from.row + 1, from.col - 2}, game, moves)
  moves = appendIfEmptyOrCapture(from, &Coord{from.row - 1, from.col + 2}, game, moves)
  moves = appendIfEmptyOrCapture(from, &Coord{from.row - 1, from.col - 2}, game, moves)
  moves = appendIfEmptyOrCapture(from, &Coord{from.row + 2, from.col + 1}, game, moves)
  moves = appendIfEmptyOrCapture(from, &Coord{from.row + 2, from.col - 1}, game, moves)
  moves = appendIfEmptyOrCapture(from, &Coord{from.row - 2, from.col + 1}, game, moves)
  return appendIfEmptyOrCapture(from, &Coord{from.row - 2, from.col - 1}, game, moves)
}

func legalBishopMovesFrom(from *Coord, game *Game) []*Move {
  moves := make([]*Move, 0, 13)
  // Up right
  moves = appendMovesInRange(
    from, game, moves,
    func(row int) int { return row + 1},
    func(col int) int { return col + 1})
  // Down right
  moves = appendMovesInRange(
    from, game, moves,
    func(row int) int { return row - 1},
    func(col int) int { return col + 1})
  // Up left
  moves = appendMovesInRange(
    from, game, moves,
    func(row int) int { return row + 1},
    func(col int) int { return col - 1})
  // Down left
  moves = appendMovesInRange(
    from, game, moves,
    func(row int) int { return row - 1},
    func(col int) int { return col - 1})
  return moves
}

func legalRookMovesFrom(from *Coord, game *Game) []*Move {
  moves := make([]*Move, 0, 14)
  // Up
  moves = appendMovesInRange(
    from, game, moves,
    func(row int) int { return row + 1},
    func(col int) int { return col})
  // Down
  moves = appendMovesInRange(
    from, game, moves,
    func(row int) int { return row - 1},
    func(col int) int { return col})
  // Left
  moves = appendMovesInRange(
    from, game, moves,
    func(row int) int { return row},
    func(col int) int { return col - 1})
  // Right
  moves = appendMovesInRange(
    from, game, moves,
    func(row int) int { return row},
    func(col int) int { return col + 1})
  return moves
}

func legalQueenMovesFrom(from *Coord, game *Game) []*Move {
  return append(
      legalBishopMovesFrom(from, game), legalRookMovesFrom(from, game)...)
}

func simpleKingMovesFrom(from *Coord, game *Game) []*Move {
  moves := make([]*Move, 0, 8)
  moves = appendIfEmptyOrCapture(from, &Coord{from.row + 1, from.col - 1}, game, moves);
  moves = appendIfEmptyOrCapture(from, &Coord{from.row + 1, from.col}, game, moves);
  moves = appendIfEmptyOrCapture(from, &Coord{from.row + 1, from.col + 1}, game, moves);
  moves = appendIfEmptyOrCapture(from, &Coord{from.row - 1, from.col - 1}, game, moves);
  moves = appendIfEmptyOrCapture(from, &Coord{from.row - 1, from.col}, game, moves);
  moves = appendIfEmptyOrCapture(from, &Coord{from.row - 1, from.col + 1}, game, moves);
  moves = appendIfEmptyOrCapture(from, &Coord{from.row, from.col - 1}, game, moves);
  moves = appendIfEmptyOrCapture(from, &Coord{from.row, from.col + 1}, game, moves);
  return moves
}

func appendIfNotCheck(moves []*Move, move *Move, game *Game) []*Move {
  if _, ok := InterpretMove(move, game); ok {
    return append(moves, move)
  }
  return moves
}

func castleKingMovesFrom(from *Coord, game *Game) []*Move {
  piece := game.board.Get(from)
  moves := make([]*Move, 0, 2)
  // Left 2
  left := MakeMove(from, &Coord{from.row, from.col - 2})
  if _, ok := interpretCastle(piece, left, game); ok {
    moves = appendIfNotCheck(moves, left, game)
  }
  // Right 2
  right := MakeMove(from, &Coord{from.row, from.col + 2})
  if _, ok := interpretCastle(piece, right, game); ok {
    moves = appendIfNotCheck(moves, right, game)
  }
  return moves
}

// Keeps appending moves in a direction until it goes off the board or we run
// into a piece
func appendMovesInRange(
    from *Coord, game *Game, moves []*Move, nextRow func(int) int,
    nextCol func(int) int) []*Move {
  piece := game.board.Get(from)
  to := &Coord{from.row, from.col}
  for i := 0; i < 7; i++ {
    to = &Coord{nextRow(to.row), nextCol(to.col)}
    if !to.InRange() {
      // Out of bounds
      break
    }
    captured := game.board.Get(to)
    if captured != nil {
      // Capture
      if captured.color != piece.color {
        moves = appendIfNotCheck(moves, MakeMove(from, to), game)
      }
      // Ran into a piece
      break
    }
    moves = appendIfNotCheck(moves, MakeMove(from, to), game)
  }
  return moves
}

func appendIfEmptyOrCapture(from *Coord, to *Coord, game *Game, moves []*Move) []*Move {
  if captured := game.board.Get(to); captured != nil {
    // Captured
    if game.board.Get(from).color != captured.color {
      return appendIfNotCheck(moves, MakeMove(from, to), game)
    }
    // Can't capture own piece
    return moves
  }
  if to.InRange() {
    // Move to empty space
    return appendIfNotCheck(moves, MakeMove(from, to), game)
  }
  return moves
}

func appendIfEnPassant(from *Coord, game *Game, moves []*Move) []*Move {
  lastEvent := game.history.GetLastEvent()
  if lastEvent == nil {
    return moves
  }
  lastTo := lastEvent.moves[0].to
  lastPiece := game.board.Get(lastTo)
  if lastPiece.name != 'p' ||
      !isPawnStart(lastEvent.moves[0].from, lastPiece.color) ||
      lastTo.row != from.row ||
      abs(lastTo.col - from.col) != 1 {
    return moves
  }
  pawn := game.board.Get(from)
  to := &Coord{getPawnForwardRow(pawn.color, from, 1), lastTo.col}
  if game.board.Get(to) == nil {
    return appendIfNotCheck(moves, MakeMove(from, to), game)
  }
  return moves
}

func appendIfEmpty(from *Coord, to *Coord, game *Game, moves []*Move) []*Move {
  if game.board.Get(to) == nil {
    return appendIfNotCheck(moves, MakeMove(from, to), game)
  }
  return moves
}

func pawnDiagonal(color Color, from *Coord, colShift int) *Coord {
  return &Coord{getPawnForwardRow(color, from, 1), from.col + colShift}
}

func getPawnForward(color Color, from *Coord, numRows int) *Coord {
  return &Coord{getPawnForwardRow(color, from, numRows), from.col};
}

func getPawnForwardRow(color Color, from *Coord, numRows int) int {
  if color == Black {
    return from.row - numRows
  }
  return from.row + numRows
}

func isPawnStart(coord *Coord, turn Color) bool {
  if turn == Black {
    return coord.row == 6
  }
  return coord.row == 1
}

// Assumes move is legal
func isCastle(move *Move, board *Board) bool {
  piece := board.Get(move.to)
  return piece != nil && piece.name == 'k' &&
      abs(move.from.col - move.to.col) > 1
}

// Assumes move is a castle move
func getRookCastleMove(move *Move) *Move {
  row := move.from.row
  if move.to.col == 6 {
    return MakeMove(&Coord{row, 7}, &Coord{row, 5})
  }
  return MakeMove(&Coord{row, 0}, &Coord{row, 3})
}
