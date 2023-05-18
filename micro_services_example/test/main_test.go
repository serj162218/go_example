package test

import (
	"os"
	"testing"

	"github.com/serj162218/go_example/micro_services_example/initializer"
)

func TestMain(m *testing.M) {
	initializer.Initialize()
	code := m.Run()
	os.Exit(code)
}
