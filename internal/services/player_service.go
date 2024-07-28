package services

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
)

type PlayerService interface {
	GetPlayersByGameID(gameID uint) ([]*domain.Player, error)
	GetPlayerByUserIDAndGameID(userID uint, gameID uint) (*domain.Player, error)
	CreatePlayersInGame(game *domain.Game) error
	CreatePlayer(userID uint, gameID uint) error
}

type PlayerServiceImpl struct {
	playerRepository repositories.PlayerRepository
	userService      UserService
}

func NewPlayerService(playerRepository repositories.PlayerRepository, userService UserService) PlayerService {
	return &PlayerServiceImpl{playerRepository: playerRepository, userService: userService}
}

func (s *PlayerServiceImpl) GetPlayersByGameID(gameID uint) ([]*domain.Player, error) {
	return s.playerRepository.GetPlayersByGameID(gameID)
}

func (s *PlayerServiceImpl) GetPlayerByUserIDAndGameID(userID uint, gameID uint) (*domain.Player, error) {
	return s.playerRepository.GetPlayerByUserIDGameID(userID, gameID)
}

func (s *PlayerServiceImpl) CreatePlayersInGame(game *domain.Game) error {
	users, err := s.userService.GetUsersByGameRoomID(game.GameRoomID)
	if err != nil {
		return err
	}

	for _, user := range users {
		if err := s.CreatePlayer(user.ID, game.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *PlayerServiceImpl) CreatePlayer(userID uint, gameID uint) error {
	gamePlayer := &domain.Player{
		UserID: userID,
		GameID: gameID,
		Score:  0,
	}
	if err := s.playerRepository.CreatePlayer(gamePlayer); err != nil {
		return err
	}
	return nil
}
