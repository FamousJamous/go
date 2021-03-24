package main

import (
  "fmt"
  "minimax"
  "strings"
)

type Color int

const (
  kX Color = iota
  kO = iota
)

func (color Color) String() string {
  if color == kX {
    return "x"
  }
  return "o"
}

type Piece byte

func (piece Piece) String() string {
  return fmt.Sprintf("%c", piece)
}

func MakePiece(color Color) Piece {
  if color == kX {
    return 'x'
  }
  return 'o'
}

type Coord struct {
  row int
  col int
}

func (coord *Coord) String() string {
  return fmt.Sprintf("[%v, %v]", coord.row, coord.col)
}

type Move struct {
  coord Coord
  color Color
}

func (move *Move) String() string {
  return fmt.Sprint(move.coord, ", ", move.color)
}

type GameError struct {
  message string
}

func (gameError *GameError) Error() string {
  return gameError.message
}

type Board struct {
  impl [3][3]Piece
}

func MakeBoard() *Board {
  board := &Board{}
  for row := 0; row < 3; row++ {
    for col := 0; col < 3; col++ {
      board.Set(row, col, ' ')
    }
  }
  return board
}

func (board* Board) Set(row int, col int, piece Piece) {
  board.impl[row][col] = piece
}

func (board* Board) Get(row int, col int) Piece {
  return board.impl[row][col]
}

func (board* Board) String() string {
  builder := &strings.Builder{}
  printCols(builder)
  for row := 0; row < 3; row++ {
    printLine(builder)
    builder.WriteString(fmt.Sprint(" ", row, " "))
    for col := 0; col < 3; col++ {
      builder.WriteString(fmt.Sprintf("+ %c ", board.Get(row, col)))
    }
    builder.WriteString(fmt.Sprintln("+", row))
  }
  printLine(builder)
  printCols(builder)
  return builder.String()
}

func printCols(builder *strings.Builder) {
  builder.WriteString("  ")
  for col := 0; col < 3; col++ {
    builder.WriteString(fmt.Sprint("   ", col))
  }
  builder.WriteByte('\n')
}

func printLine(builder *strings.Builder) {
  builder.WriteString("   ")
  for col := 0; col < 3; col++ {
    builder.WriteString("+---")
  }
  builder.WriteString("+\n");
}

func (board *Board) StringKey() string {
  builder := &strings.Builder{}
  for row := 0; row < 3; row++ {
    for col := 0; col < 3; col++ {
      builder.WriteString(board.Get(row, col).String())
    }
  }
  return builder.String()
}

type State int

const (
  kNotOver State = iota
  kXWins = iota
  kOWins = iota
  kDraw = iota
)

func (state State) String() string {
  switch state {
    case kNotOver: return "Not over"
    case kXWins: return "x wins!"
    case kOWins: return "o wins!"
    case kDraw: return "draw!"
    default: panic("Illegal state")
  }
}

type Player interface {
  GetMove() (*Move, error)
}

type Game struct {
  board *Board
  turn Color
  cachedState State
  moveHistory []*Move
}

func MakeGame() *Game {
  game := &Game{}
  game.board = MakeBoard()
  game.turn = kX
  game.moveHistory = make([]*Move, 0, 9)
  return game
}

func (game *Game) GetBoard() *Board {
  return game.board
}

func (game *Game) GetTurn() Color {
  return game.turn
}

func (game *Game) String() string {
  return game.board.String()
  builder := &strings.Builder{}
  builder.WriteString(game.GetBoard().String())
  builder.WriteByte('\n')
  return builder.String()
}

func outOfRange(val int) bool {
  return val < 0 || val > 2
}

func isIllegal(game *Game, move* Move) bool {
  coord := move.coord
  board := game.GetBoard()
  return outOfRange(coord.row) || outOfRange(coord.col) ||
      game.GetTurn() != move.color || board.Get(coord.row, coord.col) != ' '
}

func (game *Game) MakeMove(move* Move) error {
  if isIllegal(game, move) {
    return &GameError{fmt.Sprint("Illegal move: ", move)}
  }
  game.board.Set(move.coord.row, move.coord.col, MakePiece(move.color))
  if game.turn == kX {
    game.turn = kO
  } else {
    game.turn = kX
  }
  game.moveHistory = append(game.moveHistory, move)
  return nil
}

func PieceToWinner(piece Piece) State {
  switch piece {
    case 'x': return kXWins
    case 'o': return kOWins
    default: panic("Tried to get winner from invalid piece")
  }
}

func (game *Game) UndoMove() {
  lastIndex := len(game.moveHistory) - 1
  coord := game.moveHistory[lastIndex].coord
  game.GetBoard().Set(coord.row, coord.col, ' ')
  if game.turn == kX {
    game.turn = kO
  } else {
    game.turn = kX
  }
  game.moveHistory = game.moveHistory[:lastIndex]
}

func (game *Game) GetState() State {
  board := game.board
  for row := 0; row < 3; row++ {
    if board.Get(row, 0) != ' ' && board.Get(row, 0) == board.Get(row, 1) &&
        board.Get(row, 1) == board.Get(row, 2) {
      return PieceToWinner(board.Get(row, 0))
    }
  }
  for col := 0; col < 3; col++ {
    if board.Get(0, col) != ' ' && board.Get(0, col) == board.Get(1, col) &&
        board.Get(1, col) == board.Get(2, col) {
      return PieceToWinner(board.Get(0, col))
    }
  }
  if board.Get(0, 0) != ' ' && board.Get(0, 0) == board.Get(1, 1) &&
      board.Get(1, 1) == board.Get(2, 2) {
    return PieceToWinner(board.Get(0, 0))
  }
  if board.Get(2, 0) != ' ' && board.Get(2, 0) == board.Get(1, 1) &&
      board.Get(1, 1) == board.Get(0, 2) {
    return PieceToWinner(board.Get(2, 0))
  }
  for row := 0; row < 3; row++ {
    for col := 0; col < 3; col++ {
      if board.Get(row, col) == ' ' {
        return kNotOver
      }
    }
  }
  return kDraw
}

type HumanPlayer struct {
  color Color
}

func MakeHumanPlayer(color Color) *HumanPlayer {
  return &HumanPlayer{color}
}

func (player* HumanPlayer) GetMove() (*Move, error) {
  fmt.Printf("%v's turn. Enter move: ", player.color.String())
  var line string
  fmt.Scanln(&line)
  if len(line) != 2 {
    return nil, &GameError{"Invalid move. Expected format <row><col>"}
  }
  return &Move{Coord{int(line[0] - '0'), int(line[1] - '0')}, player.color},
         nil
}

type AiPlayer struct {
  state *minimax.MiniMaxState
}

func MakeAiPlayer(color Color, game *Game) *AiPlayer {
  return &AiPlayer{minimax.MakeState(&AiGame{game}, color == kO, 9)}
}

func (player *AiPlayer) GetMove() (*Move, error) {
  return player.state.GetMove().(*Move), nil
}

type AiGame struct {
  game *Game
}

func (aiGame *AiGame) String() string {
  return aiGame.game.String()
}

func (aiGame *AiGame) StringKey() string {
  return aiGame.game.GetBoard().StringKey()
}

func (aiGame *AiGame) GetAllMoves() []minimax.MiniMaxMove {
  game := aiGame.game
  board := game.GetBoard()
  moves := make([]minimax.MiniMaxMove, 0, 9)
  for row := 0; row < 3; row++ {
    for col := 0; col < 3; col++ {
      if board.Get(row, col) == ' ' {
        moves = append(moves, &Move{Coord{row, col}, game.GetTurn()})
      }
    }
  }
  return moves
}

func (aiGame *AiGame) GetScore() minimax.Score {
  switch aiGame.game.GetState() {
    case kXWins: return minimax.MaxScore
    case kOWins: return minimax.MinScore
    default: return minimax.Score(0)
  }
}

func (aiGame *AiGame) MakeMove(aiMove minimax.MiniMaxMove) {
  move := aiMove.(*Move)
  if ok := aiGame.game.MakeMove(move); ok != nil {
    panic(ok.Error())
  }
}

func (aiGame *AiGame) UndoMove() {
  aiGame.game.UndoMove()
}

type PlayerManager struct {
  xPlayer Player
  oPlayer Player
  game *Game
}

func MakePlayerManager(
    xPlayer Player, oPlayer Player, game *Game) *PlayerManager {
  return &PlayerManager{xPlayer, oPlayer, game}
}

func (manager *PlayerManager) GetCurrentPlayer() Player {
  if manager.game.GetTurn() == kX {
    return manager.xPlayer
  }
  return manager.oPlayer
}

func main() {
  game := MakeGame()
  manager := MakePlayerManager(
      MakeAiPlayer(kX, game), MakeAiPlayer(kO, game), game)
  state := kNotOver
  for ; state == kNotOver; state = game.GetState() {
    fmt.Print(game)
    move, getMoveOk := manager.GetCurrentPlayer().GetMove()
    if getMoveOk != nil {
      fmt.Println(getMoveOk.Error())
      continue
    }
    makeMoveOk := game.MakeMove(move)
    if makeMoveOk != nil {
      fmt.Println(makeMoveOk.Error())
    }
  }
  fmt.Println(game)
  fmt.Println(state)
}
