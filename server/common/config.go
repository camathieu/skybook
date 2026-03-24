package common

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
)

const envPrefix = "SKYBOOK_"

// Config is the top-level configuration.
type Config struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Defaults DefaultsConfig `toml:"defaults"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	ListenAddress string `toml:"ListenAddress"`
	ListenPort    int    `toml:"ListenPort"`
	Debug         bool   `toml:"Debug"`
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	Path string `toml:"Path"`
}

// DefaultsConfig holds default values for jump creation.
type DefaultsConfig struct {
	UnitSystem      string `toml:"UnitSystem"`
	DefaultJumpType string `toml:"DefaultJumpType"`
}

// NewConfig returns a Config with sensible defaults.
func NewConfig() *Config {
	c := &Config{}
	c.InitializeDefaults()
	return c
}

// InitializeDefaults fills zero values with sensible defaults.
func (c *Config) InitializeDefaults() {
	if c.Server.ListenAddress == "" {
		c.Server.ListenAddress = "0.0.0.0"
	}
	if c.Server.ListenPort == 0 {
		c.Server.ListenPort = 8080
	}
	if c.Database.Path == "" {
		c.Database.Path = "./skybook.db"
	}
	if c.Defaults.UnitSystem == "" {
		c.Defaults.UnitSystem = "imperial"
	}
	if c.Defaults.DefaultJumpType == "" {
		c.Defaults.DefaultJumpType = "FF"
	}
}

// Validate checks config values for obvious errors.
func (c *Config) Validate() error {
	if c.Server.ListenPort < 1 || c.Server.ListenPort > 65535 {
		return fmt.Errorf("invalid server.ListenPort: %d (must be 1-65535)", c.Server.ListenPort)
	}
	if c.Database.Path == "" {
		return fmt.Errorf("database.Path must not be empty")
	}
	validUnits := map[string]bool{"imperial": true, "metric": true}
	if !validUnits[c.Defaults.UnitSystem] {
		return fmt.Errorf("invalid defaults.UnitSystem: %q (must be 'imperial' or 'metric')", c.Defaults.UnitSystem)
	}
	return nil
}

// LoadConfig reads a TOML config file and applies defaults.
func LoadConfig(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := toml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", absPath, err)
	}

	config.InitializeDefaults()
	return config, nil
}

// ApplyEnvironment overrides config values with SKYBOOK_ prefixed environment variables.
// Pattern: SKYBOOK_SECTION_KEY (e.g. SKYBOOK_SERVER_LISTENPORT=9090)
func (c *Config) ApplyEnvironment() {
	applyEnvToStruct(envPrefix, reflect.ValueOf(c).Elem())
}

// applyEnvToStruct recursively applies environment variable overrides to a struct.
func applyEnvToStruct(prefix string, v reflect.Value) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		envKey := prefix + strings.ToUpper(field.Name)

		if field.Type.Kind() == reflect.Struct {
			applyEnvToStruct(envKey+"_", fieldVal)
			continue
		}

		envVal, ok := os.LookupEnv(envKey)
		if !ok {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			fieldVal.SetString(envVal)
		case reflect.Int:
			var n int
			if _, err := fmt.Sscanf(envVal, "%d", &n); err == nil {
				fieldVal.SetInt(int64(n))
			}
		case reflect.Bool:
			fieldVal.SetBool(envVal == "true" || envVal == "1" || envVal == "yes")
		}
	}
}
