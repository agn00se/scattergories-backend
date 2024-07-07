package models

type GameStatus string

const (
	GameStatusOngoing   GameStatus = "ongoing"
	GameStatusVoting    GameStatus = "voting"
	GameStatusCompleted GameStatus = "completed"
)
