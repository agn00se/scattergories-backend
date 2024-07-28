package services

import (
	"scattergories-backend/internal/common"

	"github.com/google/uuid"
)

var roomCapacity = 6

type GameRoomJoinService interface {
	JoinGameRoom(userID uuid.UUID, roomID uuid.UUID) error
	LeaveGameRoom(userID uuid.UUID, roomID uuid.UUID) error
}

type GameRoomJoinServiceImpl struct {
	gameRoomService GameRoomService
	userService     UserService
	gameService     GameService
}

func NewGameRoomJoinService(gameRoomService GameRoomService, userService UserService, gameService GameService) GameRoomJoinService {
	return &GameRoomJoinServiceImpl{
		gameRoomService: gameRoomService,
		userService:     userService,
		gameService:     gameService,
	}
}

func (s *GameRoomJoinServiceImpl) JoinGameRoom(userID uuid.UUID, roomID uuid.UUID) error {
	// Verify game room exists
	_, err := s.gameRoomService.GetGameRoomByID(roomID)
	if err != nil {
		return err
	}

	// Verify no game at the Ongoing or Voting stage
	if err := s.gameService.VerifyNoActiveGameInRoom(roomID); err != nil {
		return err
	}

	// Verify user exists
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify game room not full
	usersInRoom, err := s.userService.GetUsersByGameRoomID(roomID)
	if err != nil {
		return err
	}
	if len(usersInRoom) >= roomCapacity {
		return common.ErrGameRoomFull
	}

	// Update the associated game room in the user table
	user.GameRoomID = &roomID
	return s.userService.UpdateUser(user)
}

func (s *GameRoomJoinServiceImpl) LeaveGameRoom(userID uuid.UUID, roomID uuid.UUID) error {
	// Verify game room exists
	gameRoom, err := s.gameRoomService.GetGameRoomByID(roomID)
	if err != nil {
		return err
	}

	// Verify user exists
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify user is in the specified game room and remove game room association
	if user.GameRoomID != nil && *user.GameRoomID == roomID {
		user.GameRoomID = nil
	} else {
		return common.ErrUserNotInSpecifiedRoom
	}

	// Update user record
	if err := s.userService.UpdateUser(user); err != nil {
		return err
	}

	// Check if there are any users left in the game room
	usersInRoom, err := s.userService.GetUsersByGameRoomID(roomID)
	if err != nil {
		return err
	}

	// If no users left, delete the game room
	if len(usersInRoom) == 0 {
		if err := s.gameRoomService.DeleteGameRoomByID(roomID); err != nil {
			return err
		}
		return nil
	}

	// Check if the user is the host and assign a new host if needed
	if gameRoom.HostID == userID {
		if _, err := s.gameRoomService.UpdateHost(roomID, usersInRoom[0].ID); err != nil {
			return err
		}
	}

	return nil
}
