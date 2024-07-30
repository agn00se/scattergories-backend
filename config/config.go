package config

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"scattergories-backend/internal/api/websocket"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/internal/services"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AppConfig struct {
	DB                      *gorm.DB
	RedisClient             *redis.Client
	UserRepo                repositories.UserRepository
	AnswerRepo              repositories.AnswerRepository
	GamePromptRepo          repositories.GamePromptRepository
	GameRepo                repositories.GameRepository
	GameRoomRepo            repositories.GameRoomRepository
	GameRoomConfigRepo      repositories.GameRoomConfigRepository
	PlayerRepo              repositories.PlayerRepository
	PromptRepo              repositories.PromptRepository
	TokenService            services.TokenService
	UserService             services.UserService
	AuthService             services.AuthService
	UserRegistrationService services.UserRegistrationService
	PlayerService           services.PlayerService
	PromptService           services.PromptService
	GamePromptService       services.GamePromptService
	GameConfigService       services.GameConfigService
	AnswerService           services.AnswerService
	GameRoomService         services.GameRoomService
	GameService             services.GameService
	GameRoomDataService     services.GameRoomDataService
	GameRoomJoinService     services.GameRoomJoinService
	PermissionService       services.PermissionService
	MessageHandler          websocket.MessageHandler
}

type DBConfig struct {
	Host     string
	User     string
	Name     string
	Port     string
	SSLMode  string
	Password string
}

func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Name:     os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

func ConnectDB() (*gorm.DB, error) {
	dbConfig := GetDBConfig()
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s password=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Name, dbConfig.Port, dbConfig.SSLMode, dbConfig.Password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Automatically migrate the schema for all models. Order is important.
	if err := db.AutoMigrate(&domain.GameRoom{}, &domain.User{}, &domain.Game{}, &domain.Player{}, &domain.GameRoomConfig{}, &domain.Prompt{}, &domain.GamePrompt{}, &domain.Answer{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}
	return db, nil
}

func InitRedis() (*redis.Client, error) {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic("Invalid REDIS_DB value")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic("Could not connect to Redis: " + err.Error())
	}
	return redisClient, nil
}

func LoadPrompts(db *gorm.DB) {
	file, err := os.Open("config/prompts.txt")
	if err != nil {
		log.Fatalf("Failed to open prompts file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		promptText := scanner.Text()
		prompt := domain.Prompt{Text: promptText}
		if err := db.FirstOrCreate(&prompt, domain.Prompt{Text: promptText}).Error; err != nil {
			log.Printf("Failed to load prompt: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading prompts file: %v", err)
	}
}

func InitializeApp() (*AppConfig, error) {
	// Initialize DB
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	// Initialize Redis
	redisClient, err := InitRedis()
	if err != nil {
		return nil, err
	}

	// Load prompts into DB
	LoadPrompts(db)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	answerRepo := repositories.NewAnswerRepository(db)
	gamePromptRepo := repositories.NewGamePromptRepository(db)
	gameRepo := repositories.NewGameRepository(db)
	gameRoomRepo := repositories.NewGameRoomRepository(db)
	gameRoomConfigRepo := repositories.NewGameRoomConfigRepository(db)
	playerRepo := repositories.NewPlayerRepository(db)
	promptRepo := repositories.NewPromptRepository(db)

	// Initialize services
	tokenService := services.NewTokenService(redisClient)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userService, tokenService)
	userRegistrationService := services.NewUserRegistrationService(userService, authService)
	playerService := services.NewPlayerService(playerRepo, userService)
	promptService := services.NewPromptService(promptRepo)
	gamePromptService := services.NewGamePromptService(gamePromptRepo, promptService)
	gameConfigService := services.NewGameConfigService(db, gameRoomConfigRepo)
	answerService := services.NewAnswerService(db, answerRepo, gamePromptService, playerService)
	gameRoomService := services.NewGameRoomService(db, gameRoomRepo, userService, gameConfigService)
	gameService := services.NewGameService(db, gameRepo, playerService, gamePromptService, gameConfigService)
	gameRoomDataService := services.NewGameRoomDataService(db, answerService, gameService)
	gameRoomJoinService := services.NewGameRoomJoinService(db, gameRoomService, userService, gameService)
	permissionService := services.NewPermissionService(userService, gameRoomService)
	messageHandler := websocket.NewMessageHandler(gameService, gameRoomService, gameRoomDataService, permissionService, answerService, gameConfigService)

	return &AppConfig{
		DB:                      db,
		RedisClient:             redisClient,
		UserRepo:                userRepo,
		AnswerRepo:              answerRepo,
		GamePromptRepo:          gamePromptRepo,
		GameRepo:                gameRepo,
		GameRoomRepo:            gameRoomRepo,
		GameRoomConfigRepo:      gameRoomConfigRepo,
		PlayerRepo:              playerRepo,
		PromptRepo:              promptRepo,
		TokenService:            tokenService,
		UserService:             userService,
		AuthService:             authService,
		UserRegistrationService: userRegistrationService,
		PlayerService:           playerService,
		PromptService:           promptService,
		GamePromptService:       gamePromptService,
		GameConfigService:       gameConfigService,
		AnswerService:           answerService,
		GameRoomService:         gameRoomService,
		GameService:             gameService,
		GameRoomDataService:     gameRoomDataService,
		GameRoomJoinService:     gameRoomJoinService,
		PermissionService:       permissionService,
		MessageHandler:          messageHandler,
	}, nil
}
