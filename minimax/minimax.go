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
  MakeMove(move MiniMaxMove)
  UndoMove()
  String() string
  StringKey() string
}

type MiniMaxState struct {
  game MiniMaxGame
  minimizeStart bool
  maxDepth int
  visited map[string]Score
}

func MakeState(game MiniMaxGame, minimize bool, maxDepth int) *MiniMaxState {
  return &MiniMaxState{game, minimize, maxDepth, make(map[string]Score)}
}

func MiniMax(game MiniMaxGame, minimize bool, maxDepth int) MiniMaxMove {
  state := MakeState(game, minimize, maxDepth)
  move, _ := state.run(minimize, 1)
  return move
}

func (state *MiniMaxState) GetMove() MiniMaxMove {
  state.visited = make(map[string]Score)
  move, _ := state.run(state.minimizeStart, 1)
  return move
}

func (state *MiniMaxState) run(minimize bool, depth int) (MiniMaxMove, Score) {
  moves := state.game.GetAllMoves()
  if len(moves) == 0 {
    return nil, state.game.GetScore()
  }
  bestMove, bestScore :=
      state.tryMove(minimize, depth, moves[0])
  for i, n := 1, len(moves); i < n; i++ {
    move, score := state.tryMove(minimize, depth, moves[i])
    if (minimize && score < bestScore) || (!minimize && score > bestScore) {
      bestMove = move
      bestScore = score
    }
  }
  return bestMove, bestScore
}

func (state *MiniMaxState) tryMove(
    minimize bool, depth int, move MiniMaxMove) (MiniMaxMove, Score) {
  state.game.MakeMove(move)
  key := state.game.StringKey()
  if score, ok := state.visited[key]; ok {
    // Already visited
    state.game.UndoMove()
    return move, score
  }
  score := state.game.GetScore()
  // Mark visited before recursive calls
  state.visited[key] = score
  if score != MaxScore && score != MinScore && depth != state.maxDepth {
    _, score = state.run(!minimize, depth + 1)
    // Update score if we made recursive calls
    state.visited[key] = score
  }
  state.game.UndoMove()
  return move, score
}
