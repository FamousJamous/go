module jsdu/chess/main

go 1.16

replace jsdu/chess/game => ../game

require (
	jsdu/chess/game v0.0.0-00010101000000-000000000000 // indirect
	minimax v0.0.0-00010101000000-000000000000 // indirect
)

replace minimax => ../../minimax
