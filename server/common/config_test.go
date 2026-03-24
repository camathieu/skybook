package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c := NewConfig()
	if c.Server.ListenAddress != "0.0.0.0" {
		t.Errorf("expected default ListenAddress '0.0.0.0', got %q", c.Server.ListenAddress)
	}
	if c.Server.ListenPort != 8080 {
		t.Errorf("expected default ListenPort 8080, got %d", c.Server.ListenPort)
	}
	if c.Database.Path != "./skybook.db" {
		t.Errorf("expected default Database.Path './skybook.db', got %q", c.Database.Path)
	}
	if c.Defaults.UnitSystem != "imperial" {
		t.Errorf("expected default UnitSystem 'imperial', got %q", c.Defaults.UnitSystem)
	}
	if c.Defaults.DefaultJumpType != "FF" {
		t.Errorf("expected default DefaultJumpType 'FF', got %q", c.Defaults.DefaultJumpType)
	}
}

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "skybook.cfg")

	content := `
[server]
ListenAddress = "127.0.0.1"
ListenPort = 9090
Debug = true

[database]
Path = "/data/skybook.db"

[defaults]
UnitSystem = "metric"
DefaultJumpType = "WS"
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	c, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	if c.Server.ListenAddress != "127.0.0.1" {
		t.Errorf("expected ListenAddress '127.0.0.1', got %q", c.Server.ListenAddress)
	}
	if c.Server.ListenPort != 9090 {
		t.Errorf("expected ListenPort 9090, got %d", c.Server.ListenPort)
	}
	if c.Server.Debug != true {
		t.Error("expected Debug true")
	}
	if c.Database.Path != "/data/skybook.db" {
		t.Errorf("expected Path '/data/skybook.db', got %q", c.Database.Path)
	}
	if c.Defaults.UnitSystem != "metric" {
		t.Errorf("expected UnitSystem 'metric', got %q", c.Defaults.UnitSystem)
	}
}

func TestApplyEnvironment(t *testing.T) {
	c := NewConfig()

	t.Setenv("SKYBOOK_SERVER_LISTENPORT", "9999")
	t.Setenv("SKYBOOK_DATABASE_PATH", "/tmp/test.db")
	t.Setenv("SKYBOOK_SERVER_DEBUG", "true")

	c.ApplyEnvironment()

	if c.Server.ListenPort != 9999 {
		t.Errorf("expected ListenPort 9999 from env, got %d", c.Server.ListenPort)
	}
	if c.Database.Path != "/tmp/test.db" {
		t.Errorf("expected Database.Path '/tmp/test.db' from env, got %q", c.Database.Path)
	}
	if c.Server.Debug != true {
		t.Error("expected Debug true from env")
	}
}

func TestValidate(t *testing.T) {
	c := NewConfig()
	if err := c.Validate(); err != nil {
		t.Errorf("default config should validate: %v", err)
	}

	// Invalid port
	c.Server.ListenPort = 0
	if err := c.Validate(); err == nil {
		t.Error("expected error for port 0")
	}

	// Invalid unit system
	c = NewConfig()
	c.Defaults.UnitSystem = "furlongs"
	if err := c.Validate(); err == nil {
		t.Error("expected error for invalid unit system")
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path.cfg")
	if err == nil {
		t.Error("expected error for missing config file")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected os.IsNotExist error, got: %v", err)
	}
}
