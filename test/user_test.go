package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scattergories-backend/config"
	"scattergories-backend/internal/client/controllers/responses"
	"scattergories-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAllUsersShouldReturnAllUsers(t *testing.T) {
	resetDatabase()

	// Create user1
	user1 := models.User{Name: "user1", Type: models.UserTypeGuest}
	config.DB.Create(&user1)

	// Create a GameRoom with user1 as the host
	gameRoomID := uint(1)
	gameRoom := models.GameRoom{Model: gorm.Model{ID: gameRoomID}, RoomCode: "testroom", IsPrivate: false, HostID: &user1.ID}
	config.DB.Create(&gameRoom)

	// Create user2 associated with the GameRoom
	user2 := models.User{Name: "user2", Type: models.UserTypeGuest, GameRoomID: &gameRoomID}
	config.DB.Create(&user2)

	resp := makeAuthenticatedRequest(http.MethodGet, "/users", nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var users []responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	assert.Equal(t, "user1", users[0].Name)
	assert.Equal(t, "user2", users[1].Name)
	assert.Equal(t, gameRoomID, *users[1].GameRoomID)
}

func TestGetAllUsersShouldReturnNoUser(t *testing.T) {
	resetDatabase()

	resp := makeAuthenticatedRequest(http.MethodGet, "/users", nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var users []responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}

func TestGetUserShouldReturnUser(t *testing.T) {
	resetDatabase()

	// Create a test user
	user := models.User{Name: "testuser"}
	config.DB.Create(&user)

	url := fmt.Sprintf("/users/%d", user.ID)
	resp := makeAuthenticatedRequest(http.MethodGet, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedUser responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedUser)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", returnedUser.Name)
}

func TestGetUserShouldReturnUserNotFound(t *testing.T) {
	resetDatabase()

	nonExistentUserID := uint(999)

	url := fmt.Sprintf("/users/%d", nonExistentUserID)
	resp := makeAuthenticatedRequest(http.MethodGet, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotFound, response["error"])
}

func TestGetUserShouldReturnIDInvalid(t *testing.T) {
	resetDatabase()

	url := fmt.Sprintf("/users/%s", "invalidIDFormat")
	resp := makeAuthenticatedRequest(http.MethodGet, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidUserID, response["error"])
}

func TestCreateGuestUserShouldCreateUser(t *testing.T) {
	resetDatabase()

	resp := makeUnauthenticatedRequest(http.MethodPost, "/guests", nil)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var returnedUser responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedUser)
	assert.NoError(t, err)
	assert.Contains(t, returnedUser.Name, "Guest")

	verifyUserByGet(t, returnedUser)
}

func TestDeleteUserShouldDeleteUser(t *testing.T) {
	resetDatabase()

	// Create a test user
	user := models.User{Name: "testuser"}
	config.DB.Create(&user)

	url := fmt.Sprintf("/users/%d", user.ID)
	resp := makeAuthenticatedRequest(http.MethodDelete, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusNoContent, resp.Code)

	// Check that the user is deleted
	var deletedUser models.User
	err := config.DB.First(&deletedUser, user.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteUserShouldReturnUserNotFound(t *testing.T) {
	resetDatabase()

	nonExistentUserID := uint(999)

	url := fmt.Sprintf("/users/%d", nonExistentUserID)
	resp := makeAuthenticatedRequest(http.MethodDelete, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotFound, response["error"])
}

func TestDeleteUserShouldReturnIDInvalid(t *testing.T) {
	resetDatabase()

	url := fmt.Sprintf("/users/%s", "invalidIDFormat")
	resp := makeAuthenticatedRequest(http.MethodDelete, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidUserID, response["error"])
}

func verifyUserByGet(t *testing.T, returnedUser responses.UserResponse) {
	url := fmt.Sprintf("/users/%d", returnedUser.ID)
	resp := makeAuthenticatedRequest(http.MethodGet, url, nil, testUser.ID, string(testUser.Type))
	assert.Equal(t, http.StatusOK, resp.Code)

	var createdUser responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.Equal(t, returnedUser.ID, createdUser.ID)
	assert.Equal(t, returnedUser.Name, createdUser.Name)
	assert.Equal(t, returnedUser.GameRoomID, createdUser.GameRoomID)
}
