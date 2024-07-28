package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"

	"time"

	"gorm.io/gorm"
)

type GameService interface {
	StartGame(roomID uint) (*domain.Game, *domain.GameRoomConfig, []*domain.GamePrompt, error)
	EndGame(roomID uint, gameID uint) (*domain.Game, []*domain.Player, error)
	UpdateGame(game *domain.Game) error
	GetGameByID(gameID uint) (*domain.Game, error)
	VerifyNoActiveGameInRoom(roomID uint) error
	GetOngoingGameInRoom(roomID uint) (*domain.Game, error)
}

type GameServiceImpl struct {
	gameRepository    repositories.GameRepository
	playerService     PlayerService
	gamePromptService GamePromptService
	gameConfigService GameConfigService
}

func NewGameService(gameRepository repositories.GameRepository, playerService PlayerService, gamePromptService GamePromptService, gameConfigService GameConfigService) GameService {
	return &GameServiceImpl{
		gameRepository:    gameRepository,
		playerService:     playerService,
		gamePromptService: gamePromptService,
		gameConfigService: gameConfigService,
	}
}

func (s *GameServiceImpl) StartGame(roomID uint) (*domain.Game, *domain.GameRoomConfig, []*domain.GamePrompt, error) {
	// Verify no game at the Ongoing or Voting stage
	if err := s.VerifyNoActiveGameInRoom(roomID); err != nil {
		return nil, nil, nil, err
	}

	// Create a new game with the status set to Ongoing
	game := &domain.Game{
		GameRoomID: roomID,
		Status:     domain.GameStatusOngoing,
		StartTime:  time.Now(),
	}
	if err := s.gameRepository.CreateGame(game); err != nil {
		return nil, nil, nil, err
	}

	// Find all users in the GameRoom and create Player entries for the new game
	if err := s.playerService.CreatePlayersInGame(game); err != nil {
		return nil, nil, nil, err
	}

	// Load GameRoomConfig
	gameRoomConfig, err := s.gameConfigService.GetGameRoomConfigByRoomID(roomID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create and load default game prompts
	if err := s.gamePromptService.createGamePrompts(game.ID, gameRoomConfig.NumberOfPrompts); err != nil {
		return nil, nil, nil, err
	}

	gamePrompts, err := s.gamePromptService.getGamePromptsByGameID(game.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Return StartGameReponse
	return game, gameRoomConfig, gamePrompts, nil
}

func (s *GameServiceImpl) EndGame(roomID uint, gameID uint) (*domain.Game, []*domain.Player, error) {
	// Find the game, set status to completed, and update the end time
	game, err := s.GetGameByID(gameID)
	if err != nil {
		return nil, nil, err
	}
	game.Status = domain.GameStatusCompleted
	game.EndTime = time.Now()
	s.UpdateGame(game)

	// Calculate final scores
	players, err := s.playerService.GetPlayersByGameID(gameID)
	if err != nil {
		return nil, nil, err
	}

	return game, players, nil
}

func (s *GameServiceImpl) GetGameByID(gameID uint) (*domain.Game, error) {
	return s.gameRepository.GetGameByID(gameID)
}

func (s *GameServiceImpl) UpdateGame(game *domain.Game) error {
	return s.gameRepository.UpdateGame(game)
}

func (s *GameServiceImpl) VerifyNoActiveGameInRoom(roomID uint) error {
	_, err := s.gameRepository.GetGameByRoomIDAndStatus(roomID, string(domain.GameStatusOngoing), string(domain.GameStatusVoting))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // No active games found
		}
		return err
	}
	return common.ErrActiveGameExists
}

func (s *GameServiceImpl) GetOngoingGameInRoom(roomID uint) (*domain.Game, error) {
	game, err := s.gameRepository.GetGameByRoomIDAndStatus(roomID, string(domain.GameStatusOngoing))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrNoOngoingGameInRoom // No ongoing games found
		}
		return nil, err
	}
	return game, nil
}
