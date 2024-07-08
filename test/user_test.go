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
	user1 := models.User{Name: "user1"}
	config.DB.Create(&user1)

	// Create a GameRoom with user1 as the host
	gameRoomID := uint(1)
	gameRoom := models.GameRoom{Model: gorm.Model{ID: gameRoomID}, RoomCode: "testroom", IsPrivate: false, HostID: &user1.ID}
	config.DB.Create(&gameRoom)

	// Create user2 associated with the GameRoom
	user2 := models.User{Name: "user2", GameRoomID: &gameRoomID}
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
	assert.Equal(t, "User not found", response["error"])
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
	assert.Equal(t, "Invalid user ID", response["error"])
}

func TestCreateUserShouldCreateUser(t *testing.T) {
	ResetDatabase()

	req, _ := http.NewRequest(http.MethodPost, "/users", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedUser responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedUser)
	assert.NoError(t, err)
	assert.Contains(t, returnedUser.Name, "Guest")

	verifyUserByGet(t, returnedUser)
}

func TestUpdateUserShouldUpdateUser(t *testing.T) {
	ResetDatabase()

	// Create a test user
	user := models.User{Name: "testuser"}
	config.DB.Create(&user)

	// Create a test game room
	gameRoom := models.GameRoom{RoomCode: "testroom", IsPrivate: false, HostID: &user.ID}
	config.DB.Create(&gameRoom)

	// Update user request
	updatePayload := map[string]interface{}{"name": "updateduser", "room_id": gameRoom.ID}
	updateJSON, _ := json.Marshal(updatePayload)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var returnedUser responses.UserResponse
	err := json.Unmarshal(resp.Body.Bytes(), &returnedUser)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", returnedUser.Name)
	assert.Equal(t, gameRoom.ID, *returnedUser.GameRoomID)

	verifyUserByGet(t, returnedUser)
}

func TestUpdateUserShouldReturnUserNotFound(t *testing.T) {
	ResetDatabase()

	// Update user request with a non-existent user
	updatePayload := map[string]string{"name": "updateduser"}
	updateJSON, _ := json.Marshal(updatePayload)

	nonExistentUserID := uint(999)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%d", nonExistentUserID), bytes.NewBuffer(updateJSON))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])
}

func TestUpdateUserShouldReturnIDInvalid(t *testing.T) {
	ResetDatabase()

	// Update user request with invalid ID format
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", "invalidIDFormat"), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])
}

func TestUpdateUserShouldFailValidationGivenNoName(t *testing.T) {
	ResetDatabase()

	// Update user request without the required name field
	gameRoomID := uint(2)
	updatePayload := map[string]uint{"room_id": gameRoomID}
	updateJSON, _ := json.Marshal(updatePayload)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%d", 1), bytes.NewBuffer(updateJSON))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'UserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag", response["error"])
}

func TestUpdateUserShouldFailValidationGivenBlankName(t *testing.T) {
	ResetDatabase()

	// Update user request without the required name field
	updatePayload := map[string]string{"name": "   "}
	updateJSON, _ := json.Marshal(updatePayload)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%d", 1), bytes.NewBuffer(updateJSON))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Key: 'UserRequest.Name' Error:Field validation for 'Name' failed on the 'not_blank' tag", response["error"])
}

func TestDeleteUserShouldDeleteUser(t *testing.T) {
	ResetDatabase()

	// Create a test user
	user := models.User{Name: "testuser"}
	config.DB.Create(&user)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", user.ID), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

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
	assert.Equal(t, "User not found", response["error"])
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
	assert.Equal(t, "Invalid user ID", response["error"])
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
