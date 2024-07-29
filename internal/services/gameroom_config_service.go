package services

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/pkg/utils"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	defaultTimeLimit       = 100
	defaultNumberOfPrompts = 10
)

type GameConfigService interface {
	UpdateGameConfig(request *domain.GameRoomConfig) (*domain.GameRoomConfig, error)
	GetGameRoomConfigByRoomID(roomID uuid.UUID) (*domain.GameRoomConfig, error)
	CreateDefaultGameRoomConfig(gameRoomID uuid.UUID) error
}

type GameConfigServiceImpl struct {
	db                       *gorm.DB
	gameRoomConfigRepository repositories.GameRoomConfigRepository
}

func NewGameConfigService(db *gorm.DB, gameRoomConfigRepository repositories.GameRoomConfigRepository) GameConfigService {
	return &GameConfigServiceImpl{db: db, gameRoomConfigRepository: gameRoomConfigRepository}
}

func (s *GameConfigServiceImpl) UpdateGameConfig(request *domain.GameRoomConfig) (*domain.GameRoomConfig, error) {
	var gameRoomConfig *domain.GameRoomConfig

	err := utils.WithTransaction(s.db, func(tx *gorm.DB) error {

		// Fetch game room config
		gameRoomConfig, err := s.GetGameRoomConfigByRoomID(request.GameRoomID)
		if err != nil {
			return err
		}

		// Update and save game room config
		gameRoomConfig.TimeLimit = request.TimeLimit
		gameRoomConfig.NumberOfPrompts = request.NumberOfPrompts
		gameRoomConfig.Letter = strings.ToUpper(request.Letter)

		if err := s.gameRoomConfigRepository.UpdateGameRoomConfig(gameRoomConfig); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return gameRoomConfig, nil
}

func (s *GameConfigServiceImpl) GetGameRoomConfigByRoomID(roomID uuid.UUID) (*domain.GameRoomConfig, error) {
	return s.gameRoomConfigRepository.GetGameRoomConfigByRoomID(roomID)
}

func (s *GameConfigServiceImpl) CreateDefaultGameRoomConfig(gameRoomID uuid.UUID) error {
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
