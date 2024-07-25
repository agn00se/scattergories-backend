package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testSetup()
	code := m.Run()
	resetDatabase()
	os.Exit(code)
}
