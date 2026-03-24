package metadata

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/root-gg/skybook/common"
)

// Backend is the metadata backend powered by GORM + SQLite.
type Backend struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewBackend creates a new metadata backend, connecting to SQLite and running migrations.
func NewBackend(config common.DatabaseConfig, log *slog.Logger) (*Backend, error) {
	// Ensure parent directory exists
	dbDir := filepath.Dir(config.Path)
	if dbDir != "." && dbDir != "" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("create database directory %s: %w", dbDir, err)
		}
	}

	// SQLite DSN with WAL mode and busy timeout
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON", config.Path)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	if log != nil {
		log.Debug("Opening database", "path", config.Path)
	}

	db, err := gorm.Open(sqlite.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	backend := &Backend{db: db, log: log}

	// Run migrations
	if err := backend.migrate(); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return backend, nil
}

// DB returns the underlying GORM database instance.
func (b *Backend) DB() *gorm.DB {
	return b.db
}

// Shutdown closes the database connection.
func (b *Backend) Shutdown() error {
	sqlDB, err := b.db.DB()
	if err != nil {
		return fmt.Errorf("get underlying db: %w", err)
	}
	if b.log != nil {
		b.log.Info("Closing database")
	}
	return sqlDB.Close()
}

// migrate runs all database migrations using gormigrate.
func (b *Backend) migrate() error {
	m := migrations()
	if len(m) == 0 {
		// No migrations defined yet — skip
		return nil
	}
	return gormigrate.New(b.db, gormigrate.DefaultOptions, m).Migrate()
}
