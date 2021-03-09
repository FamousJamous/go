package game

import (
  "fmt"
  "strings"
)

type Game struct {
  turn Color
  board *Board
  history *History
}

func LoadGame(turn Color, board *Board, history *History) *Game {
  return &Game{turn, board, history}
}

func MakeGame() *Game {
  return &Game{White, MakeBoard(), MakeHistory()}
}

func (game *Game) MakeMove(move *FromTo) error {
  event, ok := InterpretMove(move, game)
  if !ok {
    return &GameError{fmt.Sprintf("Illegal move %v", move)}
  }
  // Handle en passant
  event.apply(game.board)
  game.history.AddEvent(event)
  // Switch turns
  return game.switchTurns()
}

func (game *Game) UndoMove() error {
  if ok := game.history.UndoMove(game.board); ok != nil {
    return ok
  }
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

func ParseFromTo(str string) *FromTo {
  return &FromTo{ParseCoord(str), ParseCoord(str[2:])}
}

func ParseMove(str string) *Move {
  return &Move{ParseFromTo(str), nil}
}

type State int

const (
  BlackWins State = iota
  WhiteWins = iota
  Draw = iota
  NotOver = iota
)

func (game *Game) GetState() State {
  return NotOver
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


