package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"scattergories-backend/config"
	"scattergories-backend/internal/api/handlers/responses"
	"scattergories-backend/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllGameRoomsShouldReturnAllRooms(t *testing.T) {
	resetDatabase()

	// Create some users to be hosts of game rooms
	host1 := domain.User{Name: "host1"}
	host2 := domain.User{Name: "host2"}
	config.DB.Create(&host1)
	config.DB.Create(&host2)

	// Create some game rooms
	gameRoom1 := domain.GameRoom{RoomCode: "room1", IsPrivate: false, HostID: host1.ID}
	gameRoom2 := domain.GameRoom{RoomCode: "room2", IsPrivate: true, HostID: host2.ID, Passcode: "secret"}
	config.DB.Create(&gameRoom1)
	config.DB.Create(&gameRoom2)

	resp := makeAuthenticatedRequest(http.MethodGet, "/game-rooms", nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var rooms []responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &rooms)
	assert.NoError(t, err)
	assert.Len(t, rooms, 2)

	assert.Equal(t, "room1", rooms[0].RoomCode)
	assert.Equal(t, "room2", rooms[1].RoomCode)
	assert.Equal(t, gameRoom1.HostID, rooms[0].HostID)
	assert.Equal(t, host1.Name, rooms[0].HostName)
	assert.Equal(t, gameRoom2.HostID, rooms[1].HostID)
	assert.Equal(t, host2.Name, rooms[1].HostName)
	assert.False(t, rooms[0].IsPrivate)
	assert.True(t, rooms[1].IsPrivate)
}

func TestGetAllGameRoomsShouldReturnNoRoom(t *testing.T) {
	resetDatabase()

	resp := makeAuthenticatedRequest(http.MethodGet, "/game-rooms", nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var rooms []responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &rooms)
	assert.NoError(t, err)
	assert.Len(t, rooms, 0)
}

func TestGetGameRoomShouldReturnRoom(t *testing.T) {
	resetDatabase()

	// Create a test user
	user := domain.User{Name: "user", Type: domain.UserTypeGuest}
	config.DB.Create(&user)

	// Create a game room with the user being the host
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false, HostID: user.ID}
	config.DB.Create(&gameRoom)
	user.GameRoomID = &gameRoom.ID
	config.DB.Save(&user)

	url := fmt.Sprintf("/game-rooms/%d", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodGet, url, nil, user.ID, string(user.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedRoom responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedRoom)
	assert.NoError(t, err)
	assert.Equal(t, gameRoom.RoomCode, returnedRoom.RoomCode)
	assert.Equal(t, gameRoom.IsPrivate, returnedRoom.IsPrivate)
	assert.Equal(t, gameRoom.HostID, returnedRoom.HostID)
	assert.Equal(t, "user", returnedRoom.HostName)
	assert.False(t, returnedRoom.IsPrivate)
}

func TestGetGameRoomShouldReturnRoomNotFound(t *testing.T) {
	resetDatabase()

	nonExistentRoomID := uint(999)

	url := fmt.Sprintf("/game-rooms/%d", nonExistentRoomID)
	resp := makeAuthenticatedRequest(http.MethodGet, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusUnauthorized, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrGameRoomNotFound, response["error"])
}

func TestGetGameRoomShouldReturnIDInvalid(t *testing.T) {
	resetDatabase()

	url := fmt.Sprintf("/game-rooms/%s", "invalidIDFormat")
	resp := makeAuthenticatedRequest(http.MethodGet, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidRoomID, response["error"])
}

func TestCreateGameRoomShouldCreateRoom(t *testing.T) {
	resetDatabase()

	// Create a test user
	host := domain.User{Name: "hostuser"}
	config.DB.Create(&host)

	// Create a game room with the user being the host
	createPayload := map[string]interface{}{
		"host_id":    host.ID,
		"is_private": false,
	}
	createJSON, _ := json.Marshal(createPayload)

	resp := makeAuthenticatedRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusCreated, resp.Code)

	var returnedRoom responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedRoom)
	assert.NoError(t, err)
	assert.Equal(t, host.ID, returnedRoom.HostID)
	assert.Equal(t, host.Name, returnedRoom.HostName)
	assert.False(t, returnedRoom.IsPrivate)
	assert.NotEmpty(t, returnedRoom.RoomCode)

	// Verify that the host user's GameRoomID is set
	var updatedHost domain.User
	config.DB.First(&updatedHost, host.ID)
	assert.Equal(t, returnedRoom.ID, *updatedHost.GameRoomID)
}

func TestCreatePrivateGameRoomShouldCreateRoom(t *testing.T) {
	resetDatabase()

	// Create a test user
	host := domain.User{Name: "hostuser"}
	config.DB.Create(&host)

	// Create a game room with the user being the host
	createPayload := map[string]interface{}{
		"host_id":    host.ID,
		"is_private": true,
		"passcode":   "secret",
	}
	createJSON, _ := json.Marshal(createPayload)

	resp := makeAuthenticatedRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusCreated, resp.Code)

	var returnedRoom responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedRoom)
	assert.NoError(t, err)
	assert.Equal(t, host.ID, returnedRoom.HostID)
	assert.Equal(t, host.Name, returnedRoom.HostName)
	assert.True(t, returnedRoom.IsPrivate)
	assert.NotEmpty(t, returnedRoom.RoomCode)
}

func TestCreateGameRoomShouldFailValidationGivenNoHost(t *testing.T) {
	resetDatabase()

	// Create a game room with the user being the host
	createPayload := map[string]interface{}{
		"is_private": false,
	}
	createJSON, _ := json.Marshal(createPayload)

	resp := makeAuthenticatedRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'GameRoomRequest.HostID' Error:Field validation for 'HostID' failed on the 'required' tag", response["error"])
}

func TestCreateGameRoomShouldFailValidationGivenPublicRoomNoPasscode(t *testing.T) {
	resetDatabase()

	// Create a test user
	host := domain.User{Name: "hostuser"}
	config.DB.Create(&host)

	// Create a game room with the user being the host
	createPayload := map[string]interface{}{
		"host_id":    host.ID,
		"is_private": true,
	}
	createJSON, _ := json.Marshal(createPayload)

	resp := makeAuthenticatedRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'GameRoomRequest.Passcode' Error:Field validation for 'Passcode' failed on the 'passcode_required_if_private' tag", response["error"])
}

func TestCreateGameRoomShouldFailGivenNonExistentHostUser(t *testing.T) {
	resetDatabase()

	// Create game room request with a non-existent host ID
	createPayload := map[string]interface{}{
		"host_id":    9999, // Non-existent user ID
		"is_private": true,
		"passcode":   "secret",
	}
	createJSON, _ := json.Marshal(createPayload)

	resp := makeAuthenticatedRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotFound, response["error"])
}

func TestCreateGameRoomShouldFailGivenDuplicateHost(t *testing.T) {
	resetDatabase()

	// Create a test user to act as the host
	host := domain.User{Name: "hostuser"}
	config.DB.Create(&host)

	// Create an existing game room with the same host
	existingGameRoom := domain.GameRoom{RoomCode: "existingroom", IsPrivate: false, HostID: host.ID}
	config.DB.Create(&existingGameRoom)

	// Create game room request with the same host ID
	createPayload := map[string]interface{}{
		"host_id":    host.ID,
		"is_private": false,
	}
	createJSON, _ := json.Marshal(createPayload)

	resp := makeAuthenticatedRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusConflict, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserIsAlreadyHostOfAnotherRoom, response["error"])
}

func TestJoinGameRoomShouldJoinGameRoom(t *testing.T) {
	resetDatabase()

	// Create test user
	user := domain.User{Name: "user"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false, HostID: testUser.ID}
	config.DB.Create(&gameRoom)

	// Join game room request
	joinPayload := map[string]interface{}{
		"user_id": user.ID,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	url := fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(joinJSON), user.ID, string(user.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User joined game room", response["message"])

	// Verify that the user's GameRoomID is set
	var updatedUser domain.User
	config.DB.First(&updatedUser, user.ID)
	assert.Equal(t, gameRoom.ID, *updatedUser.GameRoomID)
}

func TestJoinGameRoomInBetweenGamesShouldJoinGameRoom(t *testing.T) {
	createTestUser()
	resetDatabase()

	// Create test user
	user := domain.User{Name: "user"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false, HostID: testUser.ID}
	if err := config.DB.Create(&gameRoom).Error; err != nil {
		t.Fatal(err)
	}

	// Create an completed game in the room
	activeGame := domain.Game{GameRoomID: gameRoom.ID, Status: domain.GameStatusCompleted}
	config.DB.Create(&activeGame)

	// Join game room request
	joinPayload := map[string]interface{}{
		"user_id": user.ID,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	url := fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(joinJSON), user.ID, string(user.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User joined game room", response["message"])

	// Verify that the user's GameRoomID is set
	var updatedUser domain.User
	config.DB.First(&updatedUser, user.ID)
	assert.Equal(t, gameRoom.ID, *updatedUser.GameRoomID)
}

func TestJoinGameRoomShouldReturnIDInvalid(t *testing.T) {
	resetDatabase()

	url := fmt.Sprintf("/game-rooms/%s/join", "invalidIDFormat")
	resp := makeAuthenticatedRequest(http.MethodPut, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidRoomID, response["error"])
}

func TestJoinGameRoomShouldFailValidationGivenNoUser(t *testing.T) {
	resetDatabase()

	// Create a game room
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	joinPayload := map[string]interface{}{}
	joinJSON, _ := json.Marshal(joinPayload)

	url := fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(joinJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'JoinLeaveRequest.UserID' Error:Field validation for 'UserID' failed on the 'required' tag", response["error"])
}

func TestJoinGameRoomShouldReturnRoomNotFound(t *testing.T) {
	resetDatabase()

	// Create a test user
	user := domain.User{Name: "testuser"}
	config.DB.Create(&user)

	// Join game room request with non-existent room ID
	joinPayload := map[string]interface{}{
		"user_id": user.ID,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	nonExistentRoomID := uint(999)
	url := fmt.Sprintf("/game-rooms/%d/join", nonExistentRoomID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(joinJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrGameRoomNotFound, response["error"])
}

func TestJoinGameRoomShouldReturnActiveGameExists(t *testing.T) {
	resetDatabase()

	// Create a test user
	user := domain.User{Name: "testuser"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false, HostID: user.ID}
	config.DB.Create(&gameRoom)

	// Create an active game in the room
	activeGame := domain.Game{GameRoomID: gameRoom.ID, Status: domain.GameStatusOngoing}
	config.DB.Create(&activeGame)

	// Join game room request
	joinPayload := map[string]interface{}{
		"user_id": user.ID,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	url := fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(joinJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusConflict, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrActiveGameExists, response["error"])
}

func TestJoinGameRoomShouldReturnUserNotFound(t *testing.T) {
	resetDatabase()

	// Create a game room
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false, HostID: testUser.ID}
	config.DB.Create(&gameRoom)

	// Join game room request with non-existent user ID
	joinPayload := map[string]interface{}{
		"user_id": 9999,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	url := fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(joinJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotFound, response["error"])
}

func TestLeaveGameRoomShouldLeaveRoom(t *testing.T) {
	resetDatabase()

	// Create test users
	user := domain.User{Name: "user"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	// Assign user to the game room
	user.GameRoomID = &gameRoom.ID
	config.DB.Save(&user)

	// Leave game room request
	leavePayload := map[string]interface{}{
		"user_id": user.ID,
	}
	leaveJSON, _ := json.Marshal(leavePayload)

	url := fmt.Sprintf("/game-rooms/%d/leave", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(leaveJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User left game room", response["message"])

	// Verify that the user's GameRoomID is unset
	var updatedUser domain.User
	config.DB.First(&updatedUser, user.ID)
	assert.Nil(t, updatedUser.GameRoomID)
}

func TestLeaveGameRoomShouldReturnIDInvalid(t *testing.T) {
	resetDatabase()

	url := fmt.Sprintf("/game-rooms/%s/leave", "invalidIDFormat")
	resp := makeAuthenticatedRequest(http.MethodPut, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidRoomID, response["error"])
}

func TestLeaveGameRoomShouldFailValidationGivenNoUser(t *testing.T) {
	resetDatabase()

	// Create a game room
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	leavePayload := map[string]interface{}{}
	leaveJSON, _ := json.Marshal(leavePayload)

	url := fmt.Sprintf("/game-rooms/%d/leave", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(leaveJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'JoinLeaveRequest.UserID' Error:Field validation for 'UserID' failed on the 'required' tag", response["error"])
}

func TestLeaveGameRoomShouldReturnRoomNotFound(t *testing.T) {
	resetDatabase()

	// Create a test user
	user := domain.User{Name: "testuser"}
	config.DB.Create(&user)

	// Leave game room request with non-existent room ID
	leavePayload := map[string]interface{}{
		"user_id": user.ID,
	}
	leaveJSON, _ := json.Marshal(leavePayload)

	nonExistentRoomID := uint(999)
	url := fmt.Sprintf("/game-rooms/%d/leave", nonExistentRoomID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(leaveJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrGameRoomNotFound, response["error"])
}

func TestLeaveGameRoomShouldReturnUserNotFound(t *testing.T) {
	resetDatabase()

	// Create a game room
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	// Leave game room request with non-existent user ID
	leavePayload := map[string]interface{}{
		"user_id": 9999,
	}
	leaveJSON, _ := json.Marshal(leavePayload)

	url := fmt.Sprintf("/game-rooms/%d/leave", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(leaveJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotFound, response["error"])
}

func TestLeaveGameRoomUserNotInSpecifiedRoom(t *testing.T) {
	resetDatabase()

	// Create test users
	user := domain.User{Name: "user"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := domain.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	// Leave game room request with user who is not in the room
	leavePayload := map[string]interface{}{
		"user_id": user.ID,
	}
	leaveJSON, _ := json.Marshal(leavePayload)

	url := fmt.Sprintf("/game-rooms/%d/leave", gameRoom.ID)
	resp := makeAuthenticatedRequest(http.MethodPut, url, bytes.NewBuffer(leaveJSON), testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotInSpecifiedRoom, response["error"])
}

// TestJoinGameRoom - reject if game room full
// TestLeaveGameRoom - assign new host if host left
// TestLeaveGameRoom - remove game room if all user left
// TestDeleteGameRoom
