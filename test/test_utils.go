package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"scattergories-backend/config"
	"scattergories-backend/internal/api/routes"
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/validators"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	router    *gin.Engine
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	testUser  domain.User
)

func testSetup() {
	setProjectRoot()
	validators.RegisterCustomValidators()
	setupTestDB()
	router = setupRouter()
	testUser = createTestUser()
}

// setProjectRoot sets the working directory to the project root
func setProjectRoot() {
	// Get the absolute path to the project root
	projectRoot, err := filepath.Abs("../")
	if err != nil {
		panic(err)
	}
	fmt.Println("Setting working directory to:", projectRoot)

	// Set the working directory to the project root
	err = os.Chdir(projectRoot)
	if err != nil {
		panic(err)
	}
}

// setupRouter initializes the Gin router with user routes
func setupRouter() *gin.Engine {
	router := gin.Default()
	routes.RegisterUserRoutes(router)
	routes.RegisterGameRoomRoutes(router)
	return router
}

func makeUnauthenticatedRequest(method, url string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func makeAuthenticatedRequest(method, url string, body io.Reader, userID uint, userType string) *httptest.ResponseRecorder {
	token, err := generateTestToken(userID, userType)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate test token: %v", err))
	}

	req, _ := http.NewRequest(method, url, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func generateTestToken(userID uint, userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func createTestUser() domain.User {
	user := domain.User{
		Name: "Test User",
		Type: domain.UserTypeGuest,
	}
	result := config.DB.Create(&user)
	if result.Error != nil {
		panic(fmt.Sprintf("Failed to create test user: %v", result.Error))
	}
	return user
}

// setupTestDB initializes the database for testing
func setupTestDB() {
	config.ConnectDB()
	config.InitRedis()
}

func resetDatabase() {
	dropTables()
	migrateTables()
}

func dropTables() {
	config.DB.Migrator().DropTable(&domain.User{})
	config.DB.Migrator().DropTable(&domain.GameRoom{})
	config.DB.Migrator().DropTable(&domain.Game{})
	config.DB.Migrator().DropTable(&domain.Player{})
	config.DB.Migrator().DropTable(&domain.GamePrompt{})
	config.DB.Migrator().DropTable(&domain.Answer{})
	config.DB.Migrator().DropTable(&domain.GameRoomConfig{})
	config.DB.Migrator().DropTable(&domain.Prompt{})
}

func migrateTables() {
	config.DB.AutoMigrate(&domain.GameRoom{}, &domain.User{}, &domain.Game{}, &domain.Player{}, &domain.GameRoomConfig{}, &domain.Prompt{}, &domain.GamePrompt{}, &domain.Answer{})
}
