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

func (game *Game) MakeMove(fromTo *FromTo) error {
  if ok := CheckLegal(fromTo, game); ok != nil {
    return ok
  }
  // Handle en passant
  capturedPiece := game.board.Get(fromTo.to)
  if ok := game.board.MovePiece(fromTo); ok != nil {
    // Should've been validated already
    panic(ok.Error())
  }
  if isCastle(fromTo, game.board) {
    if ok := game.board.MovePiece(getRookCastleMove(fromTo)); ok != nil {
      // Should've been validated already
      panic(ok.Error())
    }
  }
  game.history.AddMove(fromTo, capturedPiece)
  return game.switchTurns()
}

func (game *Game) UndoMove() error {
  if ok := game.history.UndoMove(game.board); ok != nil {
    return ok
  }
  return game.switchTurns()
}

func (game *Game) switchTurns() error {
  if game.turn == White {
    game.turn = Black
  } else {
    game.turn = White
  }
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


