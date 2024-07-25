package services

import (
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"
)

func getPlayersByGameID(gameID uint) ([]*models.Player, error) {
	return repositories.GetPlayersByGameID(gameID)
}

func getPlayerByUserIDAndGameID(userID uint, gameID uint) (*models.Player, error) {
	return repositories.GetPlayerByUserIDGameID(userID, gameID)
}

func createPlayersInGame(users []*models.User, gameID uint) error {
	for _, user := range users {
		if err := createPlayer(user.ID, gameID); err != nil {
			return err
		}
	}
	return nil
}

func createPlayer(userID uint, gameID uint) error {
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
