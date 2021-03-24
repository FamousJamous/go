package ai

import (
  "fmt"
  "jsdu/chess/game"
  "minimax"
)

type AiGame struct {
  chessGame *game.Game
}

func MakeAiGame(chessGame *game.Game) *AiGame {
  return &AiGame{chessGame}
}

func (aiGame *AiGame) String() string {
  return aiGame.chessGame.String()
}

func (aiGame *AiGame) GetAllMoves() []minimax.MiniMaxMove {
  chessMoves := aiGame.chessGame.GetAllMoves()
  moves := make([]minimax.MiniMaxMove, 0, len(chessMoves))
  for _, chessMove := range chessMoves {
    moves = append(moves, chessMove)
    if !chessMove.InRange() {
      fmt.Printf("wtf:\n%v\nmove: %v\n", aiGame.chessGame, chessMove)
    }
  }
  return moves
}

func (aiGame *AiGame) GetScore() minimax.Score {
  /*
  state := aiGame.chessGame.GetState()
  switch state {
    case game.WhiteWins: return minimax.MaxScore
    case game.BlackWins: return minimax.MinScore
    case game.Draw: return 0
  }
  */
  board := aiGame.chessGame.GetBoard()
  return minimax.Score(
      board.GetPoints(game.White) - board.GetPoints(game.Black))
}

func (aiGame *AiGame) MakeMove(move minimax.MiniMaxMove) {
  ok := aiGame.chessGame.MakeMove(move.(*game.Move))
  if !ok {
    panic(aiGame)
  }
}

func (aiGame *AiGame) UndoMove() {
  if ok := aiGame.chessGame.UndoMove(); !ok {
    panic(aiGame)
  }
}

func (aiGame *AiGame) StringKey() string {
  return aiGame.chessGame.GetBoard().StringKey()
}

type AiPlayer struct {
  aiGame *AiGame
  state *minimax.MiniMaxState
}

func MakeAiPlayer(
  color game.Color, chessGame *game.Game, depth int,
) game.Player {
  aiGame := &AiGame{chessGame}
  return &AiPlayer{aiGame, minimax.MakeState(aiGame, color == game.Black, depth)}
}

func (player *AiPlayer) GetMove() *game.Move {
  move := player.state.GetMove()
  fmt.Printf("game score: %v\n", player.aiGame.GetScore())
  fmt.Printf("chose move: %v\n", move)
  return move.(*game.Move)
}
