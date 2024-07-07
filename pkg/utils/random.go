package utils

import (
	"math/rand"
	"time"
)

func GenerateRoomCode() string {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	code := make([]rune, 5)
	for i := range code {
		code[i] = letters[rng.Intn(len(letters))]
	}
	return string(code)
}

func GetRandomLetter() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	return string(letters[rng.Intn(len(letters))])
}
