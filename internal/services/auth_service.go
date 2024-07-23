package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"

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

func Register(userType string, name string, email string, password string) (*models.User, error) {
	var user *models.User
	var err error

	if userType == string(models.UserTypeGuest) {
		user, err = CreateGuestUser()
	} else {
		user, err = CreateRegisteredUser(name, email, password)
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

func Login(email string, password string) (string, string, error) {
	// Retrieve the user from the database
	user, err := getUserByEmail(email)
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
	computedHash := computeHash(password, salt)

	// Compare the computed hash with the stored hash
	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(computedHash), []byte(*user.PasswordHash)) != 1 {
		return "", "", common.ErrLoginFailed
	}

	// Generate access token for both guest and registered users
	accessToken, err := generateJWT(user.ID, user.Type)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token for registered users
	var refreshToken string
	if user.Type == models.UserTypeRegistered {
		refreshToken, err = generateRefreshToken(user.ID)
		if err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
}

func generateHash(password string) (string, string, error) {
	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	// Hash password with salt
	hash := computeHash(password, salt)
	return hash, base64.StdEncoding.EncodeToString(salt), nil
}

func computeHash(password string, salt []byte) string {
	// time - The number of iterations the algorithm should run.
	// memory - The amount of memory used by the algorithm in KiB.
	// threads - The number of parallel threads used for hashing.
	// keyLen - The desired length of the output hash in bytes.
	return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLen))
}
