module jsdu/chess/main

go 1.16

replace jsdu/chess/game => ../game

replace minimax => ../../minimax

replace ai => ../ai

require (
	ai v0.0.0-00010101000000-000000000000
	jsdu/chess/game v0.0.0-00010101000000-000000000000
)
