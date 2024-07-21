package services

import (
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"
)

func CreatePlayersInGame(users []*models.User, gameID uint) error {
	for _, user := range users {
		if err := CreatePlayer(user.ID, gameID); err != nil {
			return err
		}
	}
	return nil
}

func CreatePlayer(userID uint, gameID uint) error {
	gamePlayer := &models.Player{
		UserID: userID,
		GameID: gameID,
		Score:  0,
	}
	if err := repositories.CreatePlayer(gamePlayer); err != nil {
		return err
	}
	return nil
}
