package test

import (
	"fmt"
	"os"
	"path/filepath"
	"scattergories-backend/config"
	"scattergories-backend/internal/client/routes"
	"scattergories-backend/internal/models"
	"scattergories-backend/pkg/validators"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func TestSetup() {
	SetProjectRoot()
	validators.RegisterCustomValidators()
	SetupTestDB()
	router = SetupRouter()
}

// SetProjectRoot sets the working directory to the project root
func SetProjectRoot() {
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

// SetupRouter initializes the Gin router with user routes
func SetupRouter() *gin.Engine {
	router := gin.Default()
	routes.RegisterUserRoutes(router)
	routes.RegisterGameRoomRoutes(router)
	return router
}

// SetupTestDB initializes the database for testing
func SetupTestDB() {
	config.ConnectDB()
}

func ResetDatabase() {
	dropTables()
	migrateTables()
}

func dropTables() {
	config.DB.Migrator().DropTable(&models.User{})
	config.DB.Migrator().DropTable(&models.GameRoom{})
	config.DB.Migrator().DropTable(&models.Game{})
	config.DB.Migrator().DropTable(&models.Player{})
	config.DB.Migrator().DropTable(&models.GamePrompt{})
	config.DB.Migrator().DropTable(&models.Answer{})
	config.DB.Migrator().DropTable(&models.GameRoomConfig{})
	config.DB.Migrator().DropTable(&models.Prompt{})
}

func migrateTables() {
	config.DB.AutoMigrate(&models.GameRoom{}, &models.User{}, &models.Game{}, &models.Player{}, &models.GameRoomConfig{}, &models.Prompt{}, &models.GamePrompt{}, &models.Answer{})
}
