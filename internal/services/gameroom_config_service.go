package services

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/pkg/utils"
	"strings"
)

var (
	defaultTimeLimit       = 100
	defaultNumberOfPrompts = 10
)

type GameConfigService interface {
	UpdateGameConfig(request *domain.GameRoomConfig) (*domain.GameRoomConfig, error)
	GetGameRoomConfigByRoomID(roomID uint) (*domain.GameRoomConfig, error)
	CreateDefaultGameRoomConfig(gameRoomID uint) error
}

type GameConfigServiceImpl struct {
	gameRoomConfigRepository repositories.GameRoomConfigRepository
}

func NewGameConfigService(gameRoomConfigRepository repositories.GameRoomConfigRepository) GameConfigService {
	return &GameConfigServiceImpl{gameRoomConfigRepository: gameRoomConfigRepository}
}

func (s *GameConfigServiceImpl) UpdateGameConfig(request *domain.GameRoomConfig) (*domain.GameRoomConfig, error) {
	// Fetch game room config
	gameRoomConfig, err := s.GetGameRoomConfigByRoomID(request.GameRoomID)
	if err != nil {
		return nil, err
	}

	// Update and save game room config
	gameRoomConfig.TimeLimit = request.TimeLimit
	gameRoomConfig.NumberOfPrompts = request.NumberOfPrompts
	gameRoomConfig.Letter = strings.ToUpper(request.Letter)

	if err := s.gameRoomConfigRepository.UpdateGameRoomConfig(gameRoomConfig); err != nil {
		return nil, err
	}

	return gameRoomConfig, nil
}

func (s *GameConfigServiceImpl) GetGameRoomConfigByRoomID(roomID uint) (*domain.GameRoomConfig, error) {
	return s.gameRoomConfigRepository.GetGameRoomConfigByRoomID(roomID)
}

func (s *GameConfigServiceImpl) CreateDefaultGameRoomConfig(gameRoomID uint) error {
	gameRoomConfig := &domain.GameRoomConfig{
		GameRoomID:      gameRoomID,
		TimeLimit:       defaultTimeLimit,
		NumberOfPrompts: defaultNumberOfPrompts,
		Letter:          utils.GetRandomLetter(),
	}

	if err := s.gameRoomConfigRepository.CreateGameRoomConfig(gameRoomConfig); err != nil {
		return err
	}
	return nil
}
