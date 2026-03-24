package metadata

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/root-gg/skybook/common"
)

func TestNewBackend(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	config := common.DatabaseConfig{Path: dbPath}
	backend, err := NewBackend(config, slog.Default())
	if err != nil {
		t.Fatalf("NewBackend failed: %v", err)
	}
	defer backend.Shutdown()

	// Verify DB is accessible
	sqlDB, err := backend.DB().DB()
	if err != nil {
		t.Fatalf("get underlying DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("ping DB: %v", err)
	}
}

func TestNewBackend_WALMode(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	config := common.DatabaseConfig{Path: dbPath}
	backend, err := NewBackend(config, slog.Default())
	if err != nil {
		t.Fatalf("NewBackend failed: %v", err)
	}
	defer backend.Shutdown()

	// Verify WAL mode is enabled
	var journalMode string
	backend.DB().Raw("PRAGMA journal_mode").Scan(&journalMode)
	if journalMode != "wal" {
		t.Errorf("expected journal_mode 'wal', got %q", journalMode)
	}
}

func TestNewBackend_AutoCreateDir(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "sub", "dir", "test.db")

	config := common.DatabaseConfig{Path: dbPath}
	backend, err := NewBackend(config, slog.Default())
	if err != nil {
		t.Fatalf("NewBackend failed with nested dir: %v", err)
	}
	defer backend.Shutdown()
}

func TestBackend_Shutdown(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	config := common.DatabaseConfig{Path: dbPath}
	backend, err := NewBackend(config, slog.Default())
	if err != nil {
		t.Fatalf("NewBackend failed: %v", err)
	}

	if err := backend.Shutdown(); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}
