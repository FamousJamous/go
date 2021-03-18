module ai

go 1.16

replace jsdu/chess/game => ../game

replace minimax => ../../minimax

require (
	jsdu/chess/game v0.0.0-00010101000000-000000000000
	minimax v0.0.0-00010101000000-000000000000
)
