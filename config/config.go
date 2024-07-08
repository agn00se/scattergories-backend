package config

import (
	"log"
	"os"
	"scattergories-backend/internal/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE") +
		" password=" + os.Getenv("DB_PASSWORD")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Automatically migrate the schema for all models. Order is important.
	err = DB.AutoMigrate(&models.GameRoom{}, &models.User{}, &models.Game{}, &models.Player{}, &models.GamePrompt{}, &models.Answer{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
}
