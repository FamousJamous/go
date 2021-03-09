package game

import (
  "reflect"
)

func CheckLegal(fromTo *FromTo, game *Game) error {
  moves := LegalMovesFrom(fromTo.from, game)
  for _, move := range moves {
    if reflect.DeepEqual(move.fromTo, fromTo) {
      return nil
    }
  }
  return &GameError{"Illegal move"}
}

func LegalMovesFrom(from *Coord, game *Game) []*Move {
  piece := game.board.Get(from)
  if piece == nil || piece.color != game.turn {
    return make([]*Move, 0)
  }
  return legalMovesForPiece(piece, from, game)
}

func threatsForColor(color Color, game *Game) *Set {
  threats := MakeCoordSet()
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      from := &Coord{row, col}
      if piece := game.board.Get(from); piece != nil && piece.color == color {
        for _, move := range simpleMovesForPiece(piece, from, game) {
          threats.Insert(move.fromTo.to)
        }
      }
    }
  }
  return threats
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

func legalPawnMovesFrom(from *Coord, game *Game) []*Move {
  pawn := game.board.Get(from)
  moves := make([]*Move, 0, 3)
  // Forward 1
  forward1 := getPawnForward(pawn.color, from, 1)
  moves = appendIfEmpty(from, forward1, game, moves)
  // Forward 2
  if game.board.Get(forward1) == nil && isPawnStart(from, pawn.color) {
    moves = appendIfEmpty(
      from, getPawnForward(pawn.color, from, 2), game, moves)
  }
  // Left capture
  moves = appendIfCapture(
      from, getPawnDiagonal(pawn.color, from, -1), game, moves)
  // Right capture
  moves = appendIfCapture(
      from, getPawnDiagonal(pawn.color, from, 1), game, moves)
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

func castleKingMovesFrom(from *Coord, game *Game) []*Move {
  moves := make([]*Move, 0, 2)
  kingStart := getKingStart(from, game)
  if reflect.DeepEqual(*from, *kingStart) && noMovesFrom(kingStart, game) {
    // Check for threats between rook and king
    row := from.row
    threats := getOpponentThreats(game)
    // Left
    if rangeIsEmptyAndHasNoThreats(row, 1, 4, game, threats) &&
        noMovesFrom(&Coord{row, 0}, game) {
      // Rook implied by 2 space king move
      moves = append(moves, MakeMove(from, &Coord{row, 2}))
    }
    if rangeIsEmptyAndHasNoThreats(row, 5, 7, game, threats) &&
        noMovesFrom(&Coord{row, 7}, game) {
      // Rook implied by 2 space king move
      moves = append(moves, MakeMove(from, &Coord{row, 6}))
    }
  }
  return moves
}

func getKingStart(from *Coord, game *Game) *Coord {
  if game.board.Get(from).color == Black {
    return &Coord{7, 4}
  } else {
    return &Coord{0, 4}
  }
}

func getOpponentThreats(game *Game) *Set {
  if game.turn == Black {
    return threatsForColor(White, game)
  }
  return threatsForColor(Black, game)
}

func rangeIsEmptyAndHasNoThreats(
    row int, colStart int, colEnd int, game *Game, threats *Set) bool {
  for col := colStart; col < colEnd; col++ {
    coord := &Coord{row, col}
    if game.board.Get(coord) != nil || threats.Contains(coord) {
      return false
    }
  }
  return true
}

func noMovesFrom(from *Coord, game *Game) bool {
  for _, event := range game.history.AllEvents() {
    if reflect.DeepEqual(event.moves[0].from, from) {
      return false
    }
  }
  return true
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
        moves = append(moves, &Move{&FromTo{from, to}, &Captured{captured, to}})
      }
      // Ran into a piece
      break
    }
    moves = append(moves, MakeMove(from, to))
  }
  return moves
}

func appendIfEmptyOrCapture(from *Coord, to *Coord, game *Game, moves []*Move) []*Move {
  if captured := game.board.Get(to); captured != nil {
    // Captured
    if game.board.Get(from).color != captured.color {
      return append(moves, &Move{&FromTo{from, to}, &Captured{captured, to}})
    }
    // Can't capture own piece
    return moves
  }
  if to.InRange() {
    // Move to empty space
    return append(moves, MakeMove(from, to))
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
    return append(moves, &Move{&FromTo{from, to}, &Captured{lastPiece, lastTo}})
  }
  return moves
}

func appendIfEmpty(from *Coord, to *Coord, game *Game, moves []*Move) []*Move {
  if game.board.Get(to) == nil {
    return append(moves, MakeMove(from, to))
  }
  return moves
}

func appendIfCapture(
    from *Coord, to *Coord, game *Game, moves []*Move) []*Move {
  piece := game.board.Get(from)
  captured := game.board.Get(to)
  if captured == nil || captured.color == piece.color {
    return moves
  }
  return append(moves, &Move{&FromTo{from, to}, &Captured{captured, to}})
}

func getPawnDiagonal(color Color, from *Coord, colShift int) *Coord {
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

// Assumes fromTo is legal
func isCastle(fromTo *FromTo, board *Board) bool {
  piece := board.Get(fromTo.to)
  return piece != nil && piece.name == 'k' &&
      abs(fromTo.from.col - fromTo.to.col) > 1
}

// Assumes fromTo is a castle move
func getRookCastleMove(fromTo *FromTo) *FromTo {
  row := fromTo.from.row
  if fromTo.to.col == 6 {
    return &FromTo{&Coord{row, 7}, &Coord{row, 5}}
  } else {
    return &FromTo{&Coord{row, 0}, &Coord{row, 3}}
  }
}
