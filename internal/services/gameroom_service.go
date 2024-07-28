package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/pkg/utils"

	"github.com/google/uuid"
)

type GameRoomService interface {
	CreateGameRoom(hostID uuid.UUID, isPrivate bool, passcode string) (*domain.GameRoom, error)
	GetAllGameRooms() ([]*domain.GameRoom, error)
	GetGameRoomByID(roomID uuid.UUID) (*domain.GameRoom, error)
	DeleteGameRoomByID(roomID uuid.UUID) error
	UpdateHost(roomID uuid.UUID, newHostID uuid.UUID) (*domain.GameRoom, error)
	VerifyGameRoomHost(roomID uuid.UUID, userID uuid.UUID, errorMessage error) error
	VerifyHostNotInOtherRoom(hostID uuid.UUID) error
}

type GameRoomServiceImpl struct {
	gameRoomRepository repositories.GameRoomRepository
	userService        UserService
	gameConfigService  GameConfigService
}

func NewGameRoomService(gameRoomRepository repositories.GameRoomRepository, userService UserService, gameConfigService GameConfigService) GameRoomService {
	return &GameRoomServiceImpl{
		gameRoomRepository: gameRoomRepository,
		userService:        userService,
		gameConfigService:  gameConfigService,
	}
}

func (s *GameRoomServiceImpl) CreateGameRoom(hostID uuid.UUID, isPrivate bool, passcode string) (*domain.GameRoom, error) {
	// Verify that the host user exists
	host, err := s.userService.GetUserByID(hostID)
	if err != nil {
		return nil, err
	}

	// Verify that the host user is not a host in another game room
	if err := s.VerifyHostNotInOtherRoom(hostID); err != nil {
		return nil, err
	}

	// Create Game Room in the database
	gameRoom := &domain.GameRoom{
		RoomCode:  utils.GenerateRoomCode(),
		HostID:    hostID,
		IsPrivate: isPrivate,
		Passcode:  passcode,
	}
	if err := s.gameRoomRepository.CreateGameRoom(gameRoom); err != nil {
		return nil, err
	}

	// Update the user table with the associated game room id
	host.GameRoomID = &gameRoom.ID
	if err := s.userService.UpdateUser(host); err != nil {
		return nil, err
	}

	// Create default GameRoomConfig for the new GameRoom
	if err := s.gameConfigService.CreateDefaultGameRoomConfig(gameRoom.ID); err != nil {
		return nil, err
	}

	// Reload the game room with the assoicated host
	gameRoomResponse, err := s.GetGameRoomByID(gameRoom.ID)
	if err != nil {
		return nil, err
	}
	return gameRoomResponse, nil
}

func (s *GameRoomServiceImpl) GetAllGameRooms() ([]*domain.GameRoom, error) {
	return s.gameRoomRepository.GetAllGameRooms()
}

func (s *GameRoomServiceImpl) GetGameRoomByID(roomID uuid.UUID) (*domain.GameRoom, error) {
	return s.gameRoomRepository.GetGameRoomByID(roomID)
}

func (s *GameRoomServiceImpl) DeleteGameRoomByID(roomID uuid.UUID) error {
	result := s.gameRoomRepository.DeleteGameRoomByID(roomID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *GameRoomServiceImpl) UpdateHost(roomID uuid.UUID, newHostID uuid.UUID) (*domain.GameRoom, error) {
	// Get the game room
	gameRoom, err := s.GetGameRoomByID(roomID)
	if err != nil {
		return nil, err
	}

	// Verify that the new host user exists
	if _, err := s.userService.GetUserByID(newHostID); err != nil {
		return nil, err
	}

	// Verify that the host user is not a host in another game room
	if err := s.VerifyHostNotInOtherRoom(newHostID); err != nil {
		return nil, err
	}

	// Update host
	gameRoom.HostID = newHostID
	s.gameRoomRepository.UpdateGameRoom(gameRoom)

	// Reload the game room with the assoicated host
	gameRoomResponse, err := s.GetGameRoomByID(gameRoom.ID)
	if err != nil {
		return nil, err
	}
	return gameRoomResponse, nil
}

func (s *GameRoomServiceImpl) VerifyGameRoomHost(roomID uuid.UUID, userID uuid.UUID, errorMessage error) error {
	gameRoom, err := s.GetGameRoomByID(roomID)
	if err != nil {
		return err
	}

	if gameRoom.HostID != userID {
		return errorMessage
	}
	return nil
}

func (s *GameRoomServiceImpl) VerifyHostNotInOtherRoom(hostID uuid.UUID) error {
	_, err := s.gameRoomRepository.GetGameRoomGivenHost(hostID)
	if err != nil {
		if err == common.ErrGameRoomWithGivenHostNotFound {
			return nil
		}
		return err
	}
	return common.ErrUserIsAlreadyHostOfAnotherRoom
}
