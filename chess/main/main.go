package main

import (
  "ai"
  "fmt"
  "jsdu/chess/game"
  "time"
)

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

func (player *HumanPlayer) GetMove() *game.Move {
  fmt.Println(player.chessGame)
  fmt.Println("Enter move (<a-h><1-8><a-h><1-8>): ")
  var line string
  fmt.Scanln(&line)
  if len(line) != 4 {
    fmt.Println("Invalid move")
    return player.GetMove()
  }
  move := game.ParseMove(line)
  if !move.InRange() {
    fmt.Println("Invalid move")
    return player.GetMove()
  }
  return move
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
      ai.MakeAiPlayer(game.White, chessGame, 5),
      ai.MakeAiPlayer(game.Black, chessGame, 5),
      chessGame}
  lastTime := time.Now()
  for state := chessGame.GetState(); !state.IsOver();
      state = chessGame.GetState() {
    move := manager.GetCurrentPlayer().GetMove()
    if ok := chessGame.MakeMove(move); !ok {
      fmt.Printf("failed to make move: %v\n", move)
    }
    fmt.Println(chessGame)
    currentTime := time.Now()
    fmt.Printf(
      "white points: %v, black points %v, time since last event: %v\n",
      chessGame.GetBoard().GetPoints(game.White),
      chessGame.GetBoard().GetPoints(game.Black),
      currentTime.Sub(lastTime))
    lastTime = currentTime
  }
}
