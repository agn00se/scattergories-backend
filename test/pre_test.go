package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	TestSetup()
	code := m.Run()
	ResetDatabase()
	os.Exit(code)
}
