package main

import (
  "fmt"
  "jsdu/chess/game"
  "minimax"
)

type AiGame struct {
  chessGame *game.Game
}

func (aiGame *AiGame) String() string {
  return aiGame.chessGame.String()
}

func (aiGame *AiGame) GetAllMoves() []minimax.MiniMaxMove {
  moves := make([]minimax.MiniMaxMove, 0, 32)
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      chessMoves := game.LegalMovesFrom(
          game.MakeCoord(row, col), aiGame.chessGame)
      for _, move := range chessMoves {
        moves = append(moves, move)
      }
    }
  }
  return moves
}

func (aiGame *AiGame) GetScore() minimax.Score {
  score := 0
  for row := 0; row < 8; row++ {
    for col := 0; col < 8; col++ {
      score += getPieceScore(game.MakeCoord(row, col), aiGame.chessGame)
    }
  }
  return minimax.Score(score)
}

func getPieceScore(from *game.Coord, chessGame *game.Game) int {
  piece := chessGame.GetBoard().Get(from)
  if piece == nil {
    return 0
  }
  score := getPieceScoreColorless(piece)
      //len(game.LegalMovesFrom(from, chessGame))
  if piece.GetColor() == game.Black {
    return -score
  }
  return score
}

func getPieceScoreColorless(piece *game.Piece) int {
  switch piece.GetName() {
    case 'p': return 1
    case 'r': return 5
    case 'n': return 3
    case 'b': return 3
    case 'q': return 9
    case 'k': return 100
    default:
      panic(fmt.Sprintf("Unexpected piece: %c", piece.GetName()))
      return 0
  }
}

func (aiGame *AiGame) MakeMove(move minimax.MiniMaxMove) error {
  return aiGame.chessGame.MakeMove(move.(*game.Move).GetFromTo())
}

func (aiGame *AiGame) UndoMove() {
  aiGame.chessGame.UndoMove()
}

type AiPlayer struct {
  aiGame *AiGame
  color game.Color
}

func MakeAiPlayer(color game.Color, chessGame *game.Game) game.Player {
  return &AiPlayer{&AiGame{chessGame}, color}
}

func (player *AiPlayer) GetMove() *game.FromTo {
  move, _ := minimax.MiniMax(player.aiGame, player.color == game.Black, 1, 4)
  fmt.Printf("chose move: %v\n", move)
  return move.(*game.Move).GetFromTo()
}

type PlayError struct {
  message string
}

func (playError *PlayError) Error() string {
  return playError.message
}

type HumanPlayer struct {
  chessGame *game.Game
  color game.Color
}

func MakeHumanPlayer(color game.Color, chessGame *game.Game) game.Player {
  return &HumanPlayer{chessGame, color}
}

func (player *HumanPlayer) GetMove() *game.FromTo {
  fmt.Println(player.chessGame)
  fmt.Println("Enter move (<a-h><1-8><a-h><1-8>): ")
  var line string
  fmt.Scanln(&line)
  if len(line) != 4 {
    fmt.Println("Invalid move")
    return player.GetMove()
  }
  fromTo := game.ParseFromTo(line)
  if !fromTo.InRange() {
    fmt.Println("Invalid move")
    return player.GetMove()
  }
  return fromTo
}

type PlayerManager struct {
  whitePlayer game.Player
  blackPlayer game.Player
  chessGame *game.Game
}

func (manager *PlayerManager) GetCurrentPlayer() game.Player {
  if manager.chessGame.Turn() == game.Black {
    return manager.blackPlayer
  }
  return manager.whitePlayer
}

func main() {
  chessGame := game.MakeGame()
  manager := &PlayerManager{
      MakeAiPlayer(game.White, chessGame),
      MakeAiPlayer(game.Black, chessGame),
      chessGame}
  var state game.State = game.NotOver
  for state = chessGame.GetState(); state == game.NotOver; {
    fromTo := manager.GetCurrentPlayer().GetMove()
    if ok := chessGame.MakeMove(fromTo); ok != nil {
      fmt.Println(ok.Error())
    }
    fmt.Println(chessGame)
  }
}
