package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	configDirName  = ".sunday"
	configFileName = "config.json"
	configDirPerm  = 0700
	configFilePerm = 0600
)

// Config holds the authentication state for the CLI.
type Config struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserEmail    string    `json:"user_email,omitempty"`
	IdentityName string    `json:"identity_name,omitempty"`
	PINSalt      string    `json:"pin_salt,omitempty"`
	PublicKey    string    `json:"public_key,omitempty"`
	PrivateKey   string    `json:"private_key,omitempty"`
}

// Path returns the path to the config file (~/.sunday/config.json).
func Path() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fall back to current directory if home dir unavailable
		return filepath.Join(".", configDirName, configFileName)
	}
	return filepath.Join(homeDir, configDirName, configFileName)
}

// Load reads the config from disk. Returns an empty config if the file doesn't exist.
func Load() (*Config, error) {
	path := Path()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the config to disk, creating the directory if needed.
func Save(cfg *Config) error {
	path := Path()
	dir := filepath.Dir(path)

	// Create config directory with restricted permissions
	if err := os.MkdirAll(dir, configDirPerm); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	if err := os.WriteFile(path, data, configFilePerm); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// Clear deletes the config file. Returns nil if the file doesn't exist.
func Clear() error {
	path := Path()

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("removing config file: %w", err)
	}

	return nil
}
