package config

import (
	"bufio"
	"context"
	"log"
	"os"
	"scattergories-backend/internal/models"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var RedisClient *redis.Client

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
	err = DB.AutoMigrate(&models.GameRoom{}, &models.User{}, &models.Game{}, &models.Player{}, &models.GameRoomConfig{}, &models.Prompt{}, &models.GamePrompt{}, &models.Answer{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
}

func InitRedis() {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic("Invalid REDIS_DB value")
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic("Could not connect to Redis: " + err.Error())
	}
}

func LoadPrompts() {
	file, err := os.Open("config/prompts.txt")
	if err != nil {
		log.Fatalf("Failed to open prompts file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		promptText := scanner.Text()
		prompt := models.Prompt{Text: promptText}
		if err := DB.FirstOrCreate(&prompt, models.Prompt{Text: promptText}).Error; err != nil {
			log.Printf("Failed to load prompt: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading prompts file: %v", err)
	}
}
