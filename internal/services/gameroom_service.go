package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/pkg/utils"
	"time"
)

func CreateGameRoom(hostID uint, isPrivate bool, passcode string) (*domain.GameRoom, error) {
	// Verify that the host user exists
	host, err := GetUserByID(hostID)
	if err != nil {
		return nil, err
	}

	// Verify that the host user is not a host in another game room
	if err := verifyHostNotInOtherRoom(hostID); err != nil {
		return nil, err
	}

	// Create Game Room in the database
	gameRoom := &domain.GameRoom{
		RoomCode:  utils.GenerateRoomCode(),
		HostID:    hostID,
		IsPrivate: isPrivate,
		Passcode:  passcode,
	}
	if err := repositories.CreateGameRoom(gameRoom); err != nil {
		return nil, err
	}

	// Update the user table with the associated game room id
	host.GameRoomID = &gameRoom.ID
	if err := updateUser(host); err != nil {
		return nil, err
	}

	// Create default GameRoomConfig for the new GameRoom
	if err := createDefaultGameRoomConfig(gameRoom.ID); err != nil {
		return nil, err
	}

	// Reload the game room with the assoicated host
	gameRoomResponse, err := GetGameRoomByID(gameRoom.ID)
	if err != nil {
		return nil, err
	}
	return gameRoomResponse, nil
}

func GetAllGameRooms() ([]*domain.GameRoom, error) {
	return repositories.GetAllGameRooms()
}

func GetGameRoomByID(roomID uint) (*domain.GameRoom, error) {
	return repositories.GetGameRoomByID(roomID)
}

func DeleteGameRoomByID(roomID uint) error {
	result := repositories.DeleteGameRoomByID(roomID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func LoadDataForRoom(roomID uint) (*domain.Game, []*domain.Answer, error) {
	// Get the Ongoing game
	game, err := getOngoingGameInRoom(roomID)
	if err != nil {
		return nil, nil, err
	}

	// Set game status to Voting stage and update endtime
	game.Status = domain.GameStatusVoting
	game.EndTime = time.Now()
	if err := updateGame(game); err != nil {
		return nil, nil, err
	}

	// Load answers with related Player and GamePrompt (including Prompt)
	answers, err := getAnswersByGameID(game.ID)
	if err != nil {
		return nil, nil, err
	}

	return game, answers, nil
}

func updateHost(roomID uint, newHostID uint) (*domain.GameRoom, error) {
	// Get the game room
	gameRoom, err := GetGameRoomByID(roomID)
	if err != nil {
		return nil, err
	}

	// Verify that the new host user exists
	if _, err := GetUserByID(newHostID); err != nil {
		return nil, err
	}

	// Verify that the host user is not a host in another game room
	if err := verifyHostNotInOtherRoom(newHostID); err != nil {
		return nil, err
	}

	// Update host
	gameRoom.HostID = newHostID
	repositories.UpdateGameRoom(gameRoom)

	// Reload the game room with the assoicated host
	gameRoomResponse, err := GetGameRoomByID(gameRoom.ID)
	if err != nil {
		return nil, err
	}
	return gameRoomResponse, nil
}

func verifyGameRoomHost(roomID uint, userID uint, errorMessage error) error {
	gameRoom, err := GetGameRoomByID(roomID)
	if err != nil {
		return err
	}

	if gameRoom.HostID != userID {
		return errorMessage
	}
	return nil
}

func verifyHostNotInOtherRoom(hostID uint) error {
	_, err := repositories.GetGameRoomGivenHost(hostID)
	if err != nil {
		if err == common.ErrGameRoomWithGivenHostNotFound {
			return nil
		}
		return err
	}
	return common.ErrUserIsAlreadyHostOfAnotherRoom
}
