package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"scattergories-backend/config"
	"scattergories-backend/internal/client/controllers/responses"
	"scattergories-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllGameRoomsShouldReturnAllRooms(t *testing.T) {
	ResetDatabase()

	// Create some users to be hosts of game rooms
	host1 := models.User{Name: "host1"}
	host2 := models.User{Name: "host2"}
	config.DB.Create(&host1)
	config.DB.Create(&host2)

	// Create some game rooms
	gameRoom1 := models.GameRoom{RoomCode: "room1", IsPrivate: false, HostID: &host1.ID}
	gameRoom2 := models.GameRoom{RoomCode: "room2", IsPrivate: true, HostID: &host2.ID, Passcode: "secret"}
	config.DB.Create(&gameRoom1)
	config.DB.Create(&gameRoom2)

	req, _ := http.NewRequest(http.MethodGet, "/game-rooms", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var rooms []responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &rooms)
	assert.NoError(t, err)
	assert.Len(t, rooms, 2)

	assert.Equal(t, "room1", rooms[0].RoomCode)
	assert.Equal(t, "room2", rooms[1].RoomCode)
	assert.Equal(t, gameRoom1.HostID, rooms[0].HostID)
	assert.Equal(t, host1.Name, *rooms[0].HostName)
	assert.Equal(t, gameRoom2.HostID, rooms[1].HostID)
	assert.Equal(t, host2.Name, *rooms[1].HostName)
	assert.False(t, rooms[0].IsPrivate)
	assert.True(t, rooms[1].IsPrivate)
}

func TestGetAllGameRoomsShouldReturnNoRoom(t *testing.T) {
	ResetDatabase()

	req, _ := http.NewRequest(http.MethodGet, "/game-rooms", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var rooms []responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &rooms)
	assert.NoError(t, err)
	assert.Len(t, rooms, 0)
}

func TestGetGameRoomShouldReturnRoom(t *testing.T) {
	ResetDatabase()

	// Create a test user
	user := models.User{Name: "user"}
	config.DB.Create(&user)

	// Create a game room with the user being the host
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false, HostID: &user.ID}
	config.DB.Create(&gameRoom)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/game-rooms/%d", gameRoom.ID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedRoom responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedRoom)
	assert.NoError(t, err)
	assert.Equal(t, gameRoom.RoomCode, returnedRoom.RoomCode)
	assert.Equal(t, gameRoom.IsPrivate, returnedRoom.IsPrivate)
	assert.Equal(t, gameRoom.HostID, returnedRoom.HostID)
	assert.Equal(t, "user", *returnedRoom.HostName)
	assert.False(t, returnedRoom.IsPrivate)
}

func TestGetGameRoomWithoutHostShouldReturnRoomWithoutHost(t *testing.T) {
	ResetDatabase()

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/game-rooms/%d", gameRoom.ID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedRoom responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedRoom)
	assert.NoError(t, err)
	assert.Equal(t, gameRoom.RoomCode, returnedRoom.RoomCode)
	assert.Equal(t, gameRoom.IsPrivate, returnedRoom.IsPrivate)
	assert.Nil(t, returnedRoom.HostID)
	assert.Nil(t, returnedRoom.HostName)
	assert.False(t, returnedRoom.IsPrivate)
}

func TestGetGameRoomShouldReturnRoomNotFound(t *testing.T) {
	ResetDatabase()

	nonExistentRoomID := uint(999)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/game-rooms/%d", nonExistentRoomID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrGameRoomNotFound, response["error"])
}

func TestGetGameRoomShouldReturnIDInvalid(t *testing.T) {
	ResetDatabase()

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/game-rooms/%s", "invalidIDFormat"), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidRoomID, response["error"])
}

func TestCreateGameRoomShouldCreateRoom(t *testing.T) {
	ResetDatabase()

	// Create a test user
	host := models.User{Name: "hostuser"}
	config.DB.Create(&host)

	// Create a game room with the user being the host
	createPayload := map[string]interface{}{
		"host_id":    host.ID,
		"is_private": false,
	}
	createJSON, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedRoom responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedRoom)
	assert.NoError(t, err)
	assert.Equal(t, host.ID, *returnedRoom.HostID)
	assert.Equal(t, host.Name, *returnedRoom.HostName)
	assert.False(t, returnedRoom.IsPrivate)
	assert.NotEmpty(t, returnedRoom.RoomCode)

	// Verify that the host user's GameRoomID is set
	var updatedHost models.User
	config.DB.First(&updatedHost, host.ID)
	assert.Equal(t, returnedRoom.ID, *updatedHost.GameRoomID)
}

func TestCreatePrivateGameRoomShouldCreateRoom(t *testing.T) {
	ResetDatabase()

	// Create a test user
	host := models.User{Name: "hostuser"}
	config.DB.Create(&host)

	// Create a game room with the user being the host
	createPayload := map[string]interface{}{
		"host_id":    host.ID,
		"is_private": true,
		"passcode":   "secret",
	}
	createJSON, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedRoom responses.GameRoomResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedRoom)
	assert.NoError(t, err)
	assert.Equal(t, host.ID, *returnedRoom.HostID)
	assert.Equal(t, host.Name, *returnedRoom.HostName)
	assert.True(t, returnedRoom.IsPrivate)
	assert.NotEmpty(t, returnedRoom.RoomCode)
}

func TestCreateGameRoomShouldFailValidationGivenNoHost(t *testing.T) {
	ResetDatabase()

	// Create a game room with the user being the host
	createPayload := map[string]interface{}{
		"is_private": false,
	}
	createJSON, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'GameRoomRequest.HostID' Error:Field validation for 'HostID' failed on the 'required' tag", response["error"])
}

func TestCreateGameRoomShouldFailValidationGivenPublicRoomNoPasscode(t *testing.T) {
	ResetDatabase()

	// Create a test user
	host := models.User{Name: "hostuser"}
	config.DB.Create(&host)

	// Create a game room with the user being the host
	createPayload := map[string]interface{}{
		"host_id":    host.ID,
		"is_private": true,
	}
	createJSON, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'GameRoomRequest.Passcode' Error:Field validation for 'Passcode' failed on the 'passcode_required_if_private' tag", response["error"])
}

func TestCreateGameRoomShouldReturnIDInvalid(t *testing.T) {
	ResetDatabase()

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/game-rooms/%s", "invalidIDFormat"), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidRoomID, response["error"])
}

func TestCreateGameRoomShouldFailGivenNonExistentHostUser(t *testing.T) {
	ResetDatabase()

	// Create game room request with a non-existent host ID
	createPayload := map[string]interface{}{
		"host_id":    9999, // Non-existent user ID
		"is_private": true,
		"passcode":   "secret",
	}
	createJSON, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrHostNotFound, response["error"])
}

func TestCreateGameRoomShouldFailGivenDuplicateHost(t *testing.T) {
	ResetDatabase()

	// Create a test user to act as the host
	host := models.User{Name: "hostuser"}
	config.DB.Create(&host)

	// Create an existing game room with the same host
	existingGameRoom := models.GameRoom{RoomCode: "existingroom", IsPrivate: false, HostID: &host.ID}
	config.DB.Create(&existingGameRoom)

	// Create game room request with the same host ID
	createPayload := map[string]interface{}{
		"host_id":    host.ID,
		"is_private": false,
	}
	createJSON, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/game-rooms", bytes.NewBuffer(createJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserIsAlreadyHostOfAnotherRoom, response["error"])
}

func TestJoinGameRoomShouldJoinGameRoom(t *testing.T) {
	ResetDatabase()

	// Create test user
	user := models.User{Name: "user"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	// Join game room request
	joinPayload := map[string]interface{}{
		"user_id": user.ID,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID), bytes.NewBuffer(joinJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User joined game room", response["message"])

	// Verify that the user's GameRoomID is set
	var updatedUser models.User
	config.DB.First(&updatedUser, user.ID)
	assert.Equal(t, gameRoom.ID, *updatedUser.GameRoomID)
}

func TestJoinGameRoomInBetweenGamesShouldJoinGameRoom(t *testing.T) {
	ResetDatabase()

	// Create test user
	user := models.User{Name: "user"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	// Create an completed game in the room
	activeGame := models.Game{GameRoomID: gameRoom.ID, Status: models.GameStatusCompleted}
	config.DB.Create(&activeGame)

	// Join game room request
	joinPayload := map[string]interface{}{
		"user_id": user.ID,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID), bytes.NewBuffer(joinJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User joined game room", response["message"])

	// Verify that the user's GameRoomID is set
	var updatedUser models.User
	config.DB.First(&updatedUser, user.ID)
	assert.Equal(t, gameRoom.ID, *updatedUser.GameRoomID)
}

func TestJoinGameRoomShouldReturnIDInvalid(t *testing.T) {
	ResetDatabase()

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%s/join", "invalidIDFormat"), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidRoomID, response["error"])
}

func TestJoinGameRoomShouldFailValidationGivenNoUser(t *testing.T) {
	ResetDatabase()

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	joinPayload := map[string]interface{}{}
	joinJSON, _ := json.Marshal(joinPayload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID), bytes.NewBuffer(joinJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'JoinLeaveRequest.UserID' Error:Field validation for 'UserID' failed on the 'required' tag", response["error"])
}

func TestJoinGameRoomShouldReturnRoomNotFound(t *testing.T) {
	ResetDatabase()

	// Create a test user
	user := models.User{Name: "testuser"}
	config.DB.Create(&user)

	// Join game room request with non-existent room ID
	joinPayload := map[string]interface{}{
		"user_id": user.ID,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	nonExistentRoomID := uint(999)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/join", nonExistentRoomID), bytes.NewBuffer(joinJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrGameRoomNotFound, response["error"])
}

func TestJoinGameRoomShouldReturnActiveGameExists(t *testing.T) {
	ResetDatabase()

	// Create a test user
	user := models.User{Name: "testuser"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false, HostID: &user.ID}
	config.DB.Create(&gameRoom)

	// Create an active game in the room
	activeGame := models.Game{GameRoomID: gameRoom.ID, Status: models.GameStatusOngoing}
	config.DB.Create(&activeGame)

	// Join game room request
	joinPayload := map[string]interface{}{
		"user_id": user.ID,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID), bytes.NewBuffer(joinJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrActiveGameExists, response["error"])
}

func TestJoinGameRoomShouldReturnUserNotFound(t *testing.T) {
	ResetDatabase()

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	// Join game room request with non-existent user ID
	joinPayload := map[string]interface{}{
		"user_id": 9999,
	}
	joinJSON, _ := json.Marshal(joinPayload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/join", gameRoom.ID), bytes.NewBuffer(joinJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotFound, response["error"])
}

func TestLeaveGameRoomShouldLeaveRoom(t *testing.T) {
	ResetDatabase()

	// Create test users
	user := models.User{Name: "user"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	// Assign user to the game room
	user.GameRoomID = &gameRoom.ID
	config.DB.Save(&user)

	// Leave game room request
	leavePayload := map[string]interface{}{
		"user_id": user.ID,
	}
	leaveJSON, _ := json.Marshal(leavePayload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/leave", gameRoom.ID), bytes.NewBuffer(leaveJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User left game room", response["message"])

	// Verify that the user's GameRoomID is unset
	var updatedUser models.User
	config.DB.First(&updatedUser, user.ID)
	assert.Nil(t, updatedUser.GameRoomID)
}

func TestLeaveGameRoomShouldReturnIDInvalid(t *testing.T) {
	ResetDatabase()

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%s/leave", "invalidIDFormat"), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidRoomID, response["error"])
}

func TestLeaveGameRoomShouldFailValidationGivenNoUser(t *testing.T) {
	ResetDatabase()

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	joinPayload := map[string]interface{}{}
	joinJSON, _ := json.Marshal(joinPayload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/leave", gameRoom.ID), bytes.NewBuffer(joinJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'JoinLeaveRequest.UserID' Error:Field validation for 'UserID' failed on the 'required' tag", response["error"])
}

func TestLeaveGameRoomShouldReturnRoomNotFound(t *testing.T) {
	ResetDatabase()

	// Create a test user
	user := models.User{Name: "testuser"}
	config.DB.Create(&user)

	// Leave game room request with non-existent room ID
	leavePayload := map[string]interface{}{
		"user_id": user.ID,
	}
	leaveJSON, _ := json.Marshal(leavePayload)

	nonExistentRoomID := uint(999)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/leave", nonExistentRoomID), bytes.NewBuffer(leaveJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrGameRoomNotFound, response["error"])
}

func TestLeaveGameRoomShouldReturnUserNotFound(t *testing.T) {
	ResetDatabase()

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	// Leave game room request with non-existent user ID
	leavePayload := map[string]interface{}{
		"user_id": 9999,
	}
	leaveJSON, _ := json.Marshal(leavePayload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/leave", gameRoom.ID), bytes.NewBuffer(leaveJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotFound, response["error"])
}

func TestLeaveGameRoomUserNotInSpecifiedRoom(t *testing.T) {
	ResetDatabase()

	// Create test users
	user := models.User{Name: "user"}
	config.DB.Create(&user)

	// Create a game room
	gameRoom := models.GameRoom{RoomCode: "room1", IsPrivate: false}
	config.DB.Create(&gameRoom)

	// Leave game room request with user who is not in the room
	leavePayload := map[string]interface{}{
		"user_id": user.ID,
	}
	leaveJSON, _ := json.Marshal(leavePayload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/game-rooms/%d/leave", gameRoom.ID), bytes.NewBuffer(leaveJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotInSpecifiedRoom, response["error"])
}
