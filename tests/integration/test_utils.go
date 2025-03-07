package test

import (
	"challenge16/internal/config"
	"challenge16/internal/regions"
	"challenge16/internal/server"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
)

const (
	csvFile = "../../cities.csv"
	envPath = "../../.env"
)

type Response struct {
	Status       bool        `json:"status"`
	ResponseCode string      `json:"resp_code"`
	Data         interface{} `json:"data,omitempty"`
	Error        string      `json:"error,omitempty"`
}

type Permission struct {
	Included []string
	Excluded []string
}

// TestSetup contains all the dependencies needed for testing
type TestSetup struct {
	App            *fiber.App
	Cleanup        func()
	permStore      map[string]Permission
	permStoreMutex sync.RWMutex
}

// SetupIntegrationTest prepares the test environment
func SetupIntegrationTest(t *testing.T) *TestSetup {
	ts := &TestSetup{
		permStore: make(map[string]Permission),
	}

	err := regions.LoadDataIntoMap(csvFile)
	if err != nil {
		t.Fatalf("Error loading data into map: %v", err)
	}

	//initialize the environment configuration
	config.LoadEnv(envPath)

	app := server.NewServer()

	ts.App = app
	ts.Cleanup = func() {
		ts.permStoreMutex.Lock()
		ts.permStore = make(map[string]Permission)
		ts.permStoreMutex.Unlock()
	}

	return ts
}

// CleanupTest performs necessary cleanup after tests
func CleanupTest(t *testing.T, ts *TestSetup) {
	if ts.Cleanup != nil {
		ts.Cleanup()
	}
}
