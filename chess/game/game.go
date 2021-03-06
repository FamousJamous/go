package game

import (
  "fmt"
  "strings"
)

type Game struct {
  turn Color
  board *Board
  history *History
  boardCounts map[string]int
}

func LoadGame(turn Color, board *Board, history *History) *Game {
  game := &Game{turn, board, history, make(map[string]int)}
  game.boardCounts[board.StringKey()]++
  // Create boardCounts
  events := game.history.events
  for i, n := 0, len(events); i < n; i++ {
    game.UndoMove()
  }
  for _, event := range events {
    game.MakeMove(event.moves[0])
  }
  return game
}

func MakeGame() *Game {
  game := &Game{
    White, MakeBoard(), MakeHistory(), make(map[string]int)}
  game.boardCounts[game.board.StringKey()]++
  return game
}

func (game *Game) MakeMove(move *Move) bool {
  event, ok := InterpretMove(move, game)
  if !ok {
    return false
  }
  // Count board state
  game.boardCounts[game.board.StringKey()]++
  // Handle en passant
  event.apply(game.board)
  game.history.AddEvent(event)
  // Switch turns
  return game.switchTurns()
}

func (game *Game) UndoMove() bool {
  if ok := game.history.UndoMove(game.board); !ok {
    return ok
  }
  game.boardCounts[game.board.StringKey()]--
  return game.switchTurns()
}

func (game *Game) getNextTurn() Color {
  return game.turn.Other()
}

func (game *Game) switchTurns() bool {
  game.turn = game.getNextTurn()
  return true
}

func (game *Game) Turn() Color {
  return game.turn
}

func ParseCoord(str string) *Coord {
  return &Coord{int(str[1] - '1'), int(str[0] - 'a')}
}

func ParseMove(str string) *Move {
  if len(str) == 5 {
    return MakePromo(ParseCoord(str), ParseCoord(str[2:]), str[4])
  }
  return MakeMove(ParseCoord(str), ParseCoord(str[2:]))
}

type State int

const (
  BlackWins State = iota
  BlackInCheck = iota
  WhiteWins = iota
  WhiteInCheck = iota
  Draw = iota
  NotOver = iota
)

func (state State) String() string {
  switch (state) {
    case BlackWins: return "black wins"
    case BlackInCheck: return "black in check"
    case WhiteWins: return "white wins"
    case WhiteInCheck: return "white in check"
    case Draw: return "draw"
    case NotOver: return "not over"
  }
  panic(fmt.Sprintf("Unexpected state %d", state))
}

func (state State) IsOver() bool {
  switch (state) {
    case BlackWins: return true
    case WhiteWins: return true
    case Draw: return true
  }
  return false
}

func (game *Game) GetAllMoves() []*Move {
  moves := make([]*Move, 0, 64)
  for fromKey, _ := range game.board.GetPieces(game.turn) {
    moves = append(moves, LegalMovesFrom(keyToCoord(fromKey), game)...)
  }
  return moves
}

func noLegalMoves(game *Game) bool {
  for fromKey, _ := range game.board.GetPieces(game.turn) {
    if len(LegalMovesFrom(keyToCoord(fromKey), game)) > 0 {
      return false
    }
  }
  return true
}

func (game *Game) GetState() State {
  if insufficientMaterial(White, game) && insufficientMaterial(Black, game) {
    return Draw
  }
  if game.boardCounts[game.board.StringKey()] > 2 {
    return Draw
  }
  kingInCheck, otherKingInCheck := identifyChecks(game)
  if noLegalMoves(game) {
    if kingInCheck {
      if game.turn == White {
        return BlackWins
      }
      return WhiteWins
    }
    return Draw
  }
  if kingInCheck {
    return colorInCheck(game.turn)
  } else if otherKingInCheck {
    return colorInCheck(game.getNextTurn())
  }
  return NotOver
}

func colorInCheck(color Color) State {
  if color == White {
    return WhiteInCheck
  }
  return BlackInCheck
}

func insufficientMaterial(color Color, game *Game) bool {
  oneKnight := false
  oneBishop := false
  for _, piece := range game.board.GetPieces(color) {
    switch piece.name {
      case 'n':
        if oneKnight {
          return false
        }
        oneKnight = true
      case 'b':
        if oneBishop {
          return false
        }
        oneBishop = true
      case 'k':
      default:
        return false
    }
  }
  return !(oneKnight && oneBishop)
}

func (game *Game) GetBoard() BoardView {
  return game.board
}

func (game *Game) String() string {
  builder := strings.Builder{}
  builder.WriteString(game.board.String())
  builder.WriteString(game.history.String())
  builder.WriteByte('\n')
  state := game.GetState()
  if state == NotOver {
    builder.WriteString(fmt.Sprintf("%v's turn", game.turn))
  } else {
    builder.WriteString(state.String())
  }
  return builder.String()
}
