package http

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	go startGin()
	code := m.Run()
	os.Exit(code)
}
