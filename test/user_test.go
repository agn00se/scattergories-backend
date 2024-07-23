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
	"gorm.io/gorm"
)

func TestGetAllUsersShouldReturnAllUsers(t *testing.T) {
	ResetDatabase()

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

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

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
	ResetDatabase()

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var users []responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}

func TestGetUserShouldReturnUser(t *testing.T) {
	ResetDatabase()

	// Create a test user
	user := models.User{Name: "testuser"}
	config.DB.Create(&user)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", user.ID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedUser responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedUser)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", returnedUser.Name)
}

func TestGetUserShouldReturnUserNotFound(t *testing.T) {
	ResetDatabase()

	nonExistentUserID := uint(999)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", nonExistentUserID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotFound, response["error"])
}

func TestGetUserShouldReturnIDInvalid(t *testing.T) {
	ResetDatabase()

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", "invalidIDFormat"), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidUserID, response["error"])
}

func TestCreateUserShouldCreateUser(t *testing.T) {
	ResetDatabase()

	createPayload := map[string]interface{}{
		"type": "guest",
	}
	createJSON, _ := json.Marshal(createPayload)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(createJSON))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Logf("Unexpected response code: %d, body: %s", resp.Code, resp.Body.String())
	}

	assert.Equal(t, http.StatusCreated, resp.Code)

	var returnedUser responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedUser)
	assert.NoError(t, err)
	assert.Contains(t, returnedUser.Name, "Guest")

	verifyUserByGet(t, returnedUser)
}

func TestDeleteUserShouldDeleteUser(t *testing.T) {
	ResetDatabase()

	// Create a test user
	user := models.User{Name: "testuser"}
	config.DB.Create(&user)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", user.ID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)

	// Check that the user is deleted
	var deletedUser models.User
	err := config.DB.First(&deletedUser, user.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteUserShouldReturnUserNotFound(t *testing.T) {
	ResetDatabase()

	nonExistentUserID := uint(999)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", nonExistentUserID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotFound, response["error"])
}

func TestDeleteUserShouldReturnIDInvalid(t *testing.T) {
	ResetDatabase()

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", "invalidIDFormat"), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidUserID, response["error"])
}

func verifyUserByGet(t *testing.T, returnedUser responses.UserResponse) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", returnedUser.ID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var createdUser responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.Equal(t, returnedUser.ID, createdUser.ID)
	assert.Equal(t, returnedUser.Name, createdUser.Name)
	assert.Equal(t, returnedUser.GameRoomID, createdUser.GameRoomID)
}
