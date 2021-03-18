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

func getPieceScore(from *game.Coord, chessGame *game.Game) int {
  piece := chessGame.GetBoard().Get(from)
  if piece == nil {
    return 0
  }
  score := getPieceScoreColorless(piece)
  if piece.GetColor() == game.Black {
    return -score
  }
  return score
}

func getPieceScoreColorless(piece *game.Piece) int {
  switch piece.GetName() {
    case 'p': return 10
    case 'r': return 50
    case 'n': return 30
    case 'b': return 30
    case 'q': return 90
    case 'k': return 0
    default:
      panic(fmt.Sprintf("Unexpected piece: %c", piece.GetName()))
      return 0
  }
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
  color game.Color
}

func MakeAiPlayer(color game.Color, chessGame *game.Game) game.Player {
  return &AiPlayer{&AiGame{chessGame}, color}
}

func (player *AiPlayer) GetMove() *game.Move {
  move, _ := minimax.MiniMax(player.aiGame, player.color == game.Black, 1, 4)
  fmt.Printf("chose move: %v\n", move)
  return move.(*game.Move)
}
