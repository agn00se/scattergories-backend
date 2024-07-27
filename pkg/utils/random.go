package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	rng   *rand.Rand
	rngMu sync.Mutex
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GenerateRoomCode() string {
	return generateRandomString(5)
}

func GetRandomLetter() string {
	return generateRandomString(1)
}

func GenerateGuestName() string {
	rngMu.Lock()
	defer rngMu.Unlock()
	return fmt.Sprintf("Guest%d", rng.Intn(10000))
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rngMu.Lock()
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	rngMu.Unlock()
	return string(b)
}
