package minimax

type Score int

const (
  MaxScore = Score(int(kMaxUint >> 1))
  MinScore = Score(-MaxScore - 1)
)

const kMaxUint = ^uint(0)

type MiniMaxMove interface { }

type MiniMaxGame interface {
  GetAllMoves() []MiniMaxMove
  GetScore() Score
  MakeMove(move MiniMaxMove) error
  UndoMove()
  String() string
}

func MiniMax(game MiniMaxGame,
             minimize bool,
             depth int,
             maxDepth int) (MiniMaxMove, Score) {
  moves := game.GetAllMoves()
  if len(moves) == 0 {
    return nil, game.GetScore()
  }
  bestMove, bestScore := tryMove(game, minimize, depth, maxDepth, moves[0])
  for i, n := 1, len(moves); i < n; i++ {
    move, score := tryMove(game, minimize, depth, maxDepth, moves[i])
    if (minimize && score < bestScore) || (!minimize && score > bestScore) {
      bestMove = move
      bestScore = score
    }
  }
  return bestMove, bestScore
}

func tryMove(game MiniMaxGame,
             minimize bool,
             depth int,
             maxDepth int,
             move MiniMaxMove) (MiniMaxMove, Score) {
  if ok := game.MakeMove(move); ok != nil {
    panic(ok.Error())
  }
  score := game.GetScore()
  if score != MaxScore && score != MinScore && depth != maxDepth {
    _, score = MiniMax(game, !minimize, depth + 1, maxDepth)
  }
  game.UndoMove()
  return move, score
}
