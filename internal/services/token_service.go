package services

import (
	"context"
	"os"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	ctx       = context.Background()
)

type TokenService interface {
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	InvalidateToken(tokenString string) error
	IsTokenBlacklisted(tokenString string) (bool, error)
	RefreshTokens(refreshToken string) (string, string, error)
	GenerateJWT(userID uuid.UUID, userType domain.UserType) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
}

type TokenServiceImpl struct {
	redisClient *redis.Client
}

func NewTokenService(redisClient *redis.Client) TokenService {
	return &TokenServiceImpl{
		redisClient: redisClient,
	}
}

func (s *TokenServiceImpl) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, common.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, common.ErrInvalidToken
	}

	return claims, nil
}

func (s *TokenServiceImpl) InvalidateToken(tokenString string) error {
	err := s.redisClient.Set(ctx, "blacklist:"+tokenString, "true", 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *TokenServiceImpl) IsTokenBlacklisted(tokenString string) (bool, error) {
	val, err := s.redisClient.Get(ctx, "blacklist:"+tokenString).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "true", nil
}

func (s *TokenServiceImpl) RefreshTokens(refreshToken string) (string, string, error) {
	// Check if the refresh token is valid
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Extract and parse the user_id claim
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return "", "", common.ErrInvalidToken
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", err
	}

	// Check if the refresh token is still valid in Redis
	if _, err := s.redisClient.Get(ctx, "refresh_token:"+refreshToken).Result(); err != nil {
		if err == redis.Nil {
			return "", "", common.ErrInvalidToken
		}
		return "", "", err
	}

	// Invalidate the old refresh token
	if err := s.redisClient.Del(ctx, "refresh_token:"+refreshToken).Err(); err != nil {
		return "", "", err
	}

	// Generate new access token
	userType := domain.UserType(claims["user_type"].(string))
	newAccessToken, err := s.GenerateJWT(userID, userType)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	newRefreshToken, err := s.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *TokenServiceImpl) GenerateJWT(userID uuid.UUID, userType domain.UserType) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID.String(),
		"user_type": string(userType),
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 1 hour expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *TokenServiceImpl) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	duration := time.Hour * 24 * 7

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(duration).Unix(), // 1 week expiration
	})

	tokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	err = s.redisClient.Set(ctx, "refresh_token:"+tokenString, userID, duration).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
