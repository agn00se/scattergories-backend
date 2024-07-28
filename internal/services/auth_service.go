package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"scattergories-backend/internal/common"

	"golang.org/x/crypto/argon2"
)

type ArgonParams struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

var params = ArgonParams{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

type AuthService interface {
	Login(email string, password string) (string, string, error)
	GenerateHash(password string) (string, string, error)
	ComputeHash(password string, salt []byte) string
}

type AuthServiceImpl struct {
	userService  UserService
	tokenService TokenService
}

func NewAuthService(userService UserService, tokenService TokenService) AuthService {
	return &AuthServiceImpl{userService: userService, tokenService: tokenService}
}

func (s *AuthServiceImpl) Login(email string, password string) (string, string, error) {
	// Retrieve the user from the database
	user, err := s.userService.GetUserByEmail(email)
	if err != nil {
		return "", "", common.ErrLoginFailed
	}

	// Decode the stored salt
	salt, err := base64.StdEncoding.DecodeString(*user.Salt)
	if err != nil {
		log.Println("Error decoding stored salt:", err)
		return "", "", common.ErrLoginFailed
	}

	// Compute the hash of the provided password using the decoded salt
	computedHash := s.ComputeHash(password, salt)

	// Compare the computed hash with the stored hash
	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(computedHash), []byte(*user.PasswordHash)) != 1 {
		return "", "", common.ErrLoginFailed
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateJWT(user.ID, user.Type)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthServiceImpl) GenerateHash(password string) (string, string, error) {
	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	// Hash password with salt
	hash := s.ComputeHash(password, salt)
	return hash, base64.StdEncoding.EncodeToString(salt), nil
}

func (s *AuthServiceImpl) ComputeHash(password string, salt []byte) string {
	// time - The number of iterations the algorithm should run.
	// memory - The amount of memory used by the algorithm in KiB.
	// threads - The number of parallel threads used for hashing.
	// keyLen - The desired length of the output hash in bytes.
	return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLen))
}
