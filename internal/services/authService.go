package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
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

func Login(email string, password string) error {
	user, err := getUserByEmail(email)
	if err != nil {
		return common.ErrLoginFailed
	}

	salt, err := base64.StdEncoding.DecodeString(*user.Salt)
	if err != nil {
		log.Println("Error decoding stored salt:", err)
		return fmt.Errorf("login failed")
	}

	// Compute the hash of the provided password using the stored salt
	computedHash := computeHash(password, salt)

	// Compare the computed hash with the stored hash
	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(computedHash), []byte(*user.PasswordHash)) != 1 {
		return common.ErrLoginFailed
	}

	return nil
}

func generateHash(password string) (string, string, error) {
	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		fmt.Println("Failed to generate salt for password hashing")
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
