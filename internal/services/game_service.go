package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/pkg/utils"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameService interface {
	StartGame(roomID uuid.UUID) (*domain.Game, *domain.GameRoomConfig, []*domain.GamePrompt, error)
	EndGame(roomID uuid.UUID, gameID uuid.UUID) (*domain.Game, []*domain.Player, error)
	UpdateGame(game *domain.Game) error
	GetGameByID(gameID uuid.UUID) (*domain.Game, error)
	VerifyNoActiveGameInRoom(roomID uuid.UUID) error
	GetOngoingGameInRoom(roomID uuid.UUID) (*domain.Game, error)
}

type GameServiceImpl struct {
	db                *gorm.DB
	gameRepository    repositories.GameRepository
	playerService     PlayerService
	gamePromptService GamePromptService
	gameConfigService GameConfigService
}

func NewGameService(db *gorm.DB, gameRepository repositories.GameRepository, playerService PlayerService, gamePromptService GamePromptService, gameConfigService GameConfigService) GameService {
	return &GameServiceImpl{
		db:                db,
		gameRepository:    gameRepository,
		playerService:     playerService,
		gamePromptService: gamePromptService,
		gameConfigService: gameConfigService,
	}
}

func (s *GameServiceImpl) StartGame(roomID uuid.UUID) (*domain.Game, *domain.GameRoomConfig, []*domain.GamePrompt, error) {
	var game *domain.Game
	var gameRoomConfig *domain.GameRoomConfig
	var gamePrompts []*domain.GamePrompt

	err := utils.WithTransaction(s.db, func(tx *gorm.DB) error {

		// Verify no game at the Ongoing or Voting stage
		if err := s.VerifyNoActiveGameInRoom(roomID); err != nil {
			return err
		}

		// Create a new game with the status set to Ongoing
		game = &domain.Game{
			GameRoomID: roomID,
			Status:     domain.GameStatusOngoing,
			StartTime:  time.Now(),
		}
		if err := s.gameRepository.CreateGame(game); err != nil {
			return err
		}

		// Find all users in the GameRoom and create Player entries for the new game
		if err := s.playerService.CreatePlayersInGame(game); err != nil {
			return err
		}

		// Load GameRoomConfig
		var err error
		gameRoomConfig, err = s.gameConfigService.GetGameRoomConfigByRoomID(roomID)
		if err != nil {
			return err
		}

		// Create and load default game prompts
		if err := s.gamePromptService.createGamePrompts(game.ID, gameRoomConfig.NumberOfPrompts); err != nil {
			return err
		}

		gamePrompts, err = s.gamePromptService.getGamePromptsByGameID(game.ID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}
	return game, gameRoomConfig, gamePrompts, nil
}

func (s *GameServiceImpl) EndGame(roomID uuid.UUID, gameID uuid.UUID) (*domain.Game, []*domain.Player, error) {
	var game *domain.Game
	var players []*domain.Player

	err := utils.WithTransaction(s.db, func(tx *gorm.DB) error {

		// Find the game, set status to completed, and update the end time
		var err error
		game, err = s.GetGameByID(gameID)
		if err != nil {
			return err
		}
		game.Status = domain.GameStatusCompleted
		game.EndTime = time.Now()
		s.UpdateGame(game)

		// Calculate final scores
		players, err = s.playerService.GetPlayersByGameID(gameID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return game, players, nil
}

func (s *GameServiceImpl) GetGameByID(gameID uuid.UUID) (*domain.Game, error) {
	return s.gameRepository.GetGameByID(gameID)
}

func (s *GameServiceImpl) UpdateGame(game *domain.Game) error {
	return s.gameRepository.UpdateGame(game)
}

func (s *GameServiceImpl) VerifyNoActiveGameInRoom(roomID uuid.UUID) error {
	_, err := s.gameRepository.GetGameByRoomIDAndStatus(roomID, string(domain.GameStatusOngoing), string(domain.GameStatusVoting))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // No active games found
		}
		return err
	}
	return common.ErrActiveGameExists
}

func (s *GameServiceImpl) GetOngoingGameInRoom(roomID uuid.UUID) (*domain.Game, error) {
	game, err := s.gameRepository.GetGameByRoomIDAndStatus(roomID, string(domain.GameStatusOngoing))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrNoOngoingGameInRoom // No ongoing games found
		}
		return nil, err
	}
	return game, nil
}
