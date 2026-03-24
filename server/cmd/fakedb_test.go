package cmd

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/root-gg/skybook/common"
	"github.com/root-gg/skybook/metadata"
)

func TestFakeDB_HappyPath(t *testing.T) {
	// Create a temporary database file
	tempDir, err := os.MkdirTemp("", "skybook_fakedb_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // cleanup

	dbPath := filepath.Join(tempDir, "test-skybook.db")

	// Set CLI flag variables
	fakedbJumps = 50
	fakedbOutput = dbPath

	// Run the command logic
	runFakedb(nil, nil)

	// Verify the database was created and populated
	config := common.DatabaseConfig{Path: dbPath}
	backend, err := metadata.NewBackend(config, slog.Default())
	if err != nil {
		t.Fatalf("failed to open generated database: %v", err)
	}
	defer backend.Shutdown()

	count, err := backend.CountJumps(1) // user ID 1
	if err != nil {
		t.Fatalf("failed to count jumps: %v", err)
	}

	if count != 50 {
		t.Errorf("expected 50 jumps, got %d", count)
	}

	// Verify invariants: contiguous numbering
	jumps, total, err := backend.GetJumps(1, 0, 100, "number", "asc", metadata.JumpFilters{})
	if err != nil {
		t.Fatalf("failed to get jumps: %v", err)
	}

	if total != 50 || len(jumps) != 50 {
		t.Fatalf("expected to retrieve 50 jumps, got %d (total %d)", len(jumps), total)
	}

	for i, j := range jumps {
		expectedNumber := uint(i + 1)
		if j.Number != expectedNumber {
			t.Errorf("expected jump number %d, got %d", expectedNumber, j.Number)
		}
		if j.UserID != 1 {
			t.Errorf("expected user ID 1, got %d", j.UserID)
		}
	}
}
