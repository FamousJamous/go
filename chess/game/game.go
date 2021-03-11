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
  game.boardCounts[board.stringKey()]++
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
  game := &Game{White, MakeBoard(), MakeHistory(), make(map[string]int)}
  game.boardCounts[game.board.stringKey()]++
  return game
}

func (game *Game) MakeMove(move *Move) error {
  event, ok := InterpretMove(move, game)
  if !ok {
    return &GameError{fmt.Sprintf("Illegal move %v", move)}
  }
  // Handle en passant
  event.apply(game.board)
  game.history.AddEvent(event)
  // Count board state
  game.boardCounts[game.board.stringKey()]++
  // Switch turns
  return game.switchTurns()
}

func (game *Game) UndoMove() error {
  if ok := game.history.UndoMove(game.board); ok != nil {
    return ok
  }
  game.boardCounts[game.board.stringKey()]--
  return game.switchTurns()
}

func (game *Game) getNextTurn() Color {
  return game.turn.Other()
}

func (game *Game) switchTurns() error {
  game.turn = game.getNextTurn()
  return nil
}

func (game *Game) Turn() Color {
  return game.turn
}

func ParseCoord(str string) *Coord {
  return &Coord{int(str[1] - '1'), int(str[0] - 'a')}
}

func ParseMove(str string) *Move {
  return &Move{ParseCoord(str), ParseCoord(str[2:])}
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

func (game *Game) GetAllMoves() []*Move {
  moves := make([]*Move, 0, 64)
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      moves = append(moves, LegalMovesFrom(&Coord{row, col}, game)...)
    }
  }
  return moves
}

func (game *Game) GetState() State {
  if insufficientMaterial(White, game) && insufficientMaterial(Black, game) {
    return Draw
  }
  if game.boardCounts[game.board.stringKey()] > 2 {
    return Draw
  }
  moves := game.GetAllMoves()
  kingInCheck, otherKingInCheck := identifyChecks(game)
  if len(moves) == 0 {
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
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      if piece := game.board.Get(&Coord{row, col}); piece != nil &&
          piece.color == color {
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
  builder.WriteString(fmt.Sprintf("\n%v's turn", game.turn))
  return builder.String()
}


