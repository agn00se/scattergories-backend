package services

import (
	"context"
	"os"
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	ctx       = context.Background()
)

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
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

func InvalidateToken(tokenString string) error {
	err := config.RedisClient.Set(ctx, "blacklist:"+tokenString, "true", 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func IsTokenBlacklisted(tokenString string) (bool, error) {
	val, err := config.RedisClient.Get(ctx, "blacklist:"+tokenString).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "true", nil
}

func RefreshTokens(refreshToken string) (string, string, error) {
	// Check if the refresh token is valid
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	userID := uint(claims["user_id"].(float64))

	// Check if the refresh token is still valid in Redis
	if _, err := config.RedisClient.Get(ctx, "refresh_token:"+refreshToken).Result(); err != nil {
		if err == redis.Nil {
			return "", "", common.ErrInvalidToken
		}
		return "", "", err
	}

	// Invalidate the old refresh token
	if err := config.RedisClient.Del(ctx, "refresh_token:"+refreshToken).Err(); err != nil {
		return "", "", err
	}

	// Generate new access token
	userType := models.UserType(claims["user_type"].(string))
	newAccessToken, err := generateJWT(userID, userType)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	newRefreshToken, err := generateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func generateJWT(userID uint, userType models.UserType) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_type": string(userType),
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 1 hour expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func generateRefreshToken(userID uint) (string, error) {
	duration := time.Hour * 24 * 7

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(), // 1 week expiration
	})

	tokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	err = config.RedisClient.Set(ctx, "refresh_token:"+tokenString, userID, duration).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
