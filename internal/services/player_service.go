package services

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
)

func getPlayersByGameID(gameID uint) ([]*domain.Player, error) {
	return repositories.GetPlayersByGameID(gameID)
}

func getPlayerByUserIDAndGameID(userID uint, gameID uint) (*domain.Player, error) {
	return repositories.GetPlayerByUserIDGameID(userID, gameID)
}

func createPlayersInGame(users []*domain.User, gameID uint) error {
	for _, user := range users {
		if err := createPlayer(user.ID, gameID); err != nil {
			return err
		}
	}
	return nil
}

func createPlayer(userID uint, gameID uint) error {
	gamePlayer := &domain.Player{
		UserID: userID,
		GameID: gameID,
		Score:  0,
	}
	if err := repositories.CreatePlayer(gamePlayer); err != nil {
		return err
	}
	return nil
}
