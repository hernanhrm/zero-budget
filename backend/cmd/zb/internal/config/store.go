package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrConfigNotFound = errors.New("cli config not found")

type Config struct {
	APIURL      string `json:"api_url"`
	IdentityURL string `json:"identity_url"`
	APIKey      string `json:"api_key"`
}

type Store struct {
	path string
}

func NewStore(path string) (Store, error) {
	if strings.TrimSpace(path) == "" {
		defaultPath, err := DefaultPath()
		if err != nil {
			return Store{}, err
		}

		path = defaultPath
	}

	return Store{path: path}, nil
}

func DefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config directory: %w", err)
	}

	return filepath.Join(configDir, "zero-budget", "cli.json"), nil
}

func (s Store) Path() string {
	return s.path
}

func (s Store) Load() (Config, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("%w: run `zb login` first", ErrConfigNotFound)
		}

		return Config{}, fmt.Errorf("read CLI config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode CLI config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (s Store) Save(cfg Config) error {
	cfg = cfg.normalized()

	if err := cfg.Validate(); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode CLI config: %w", err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return fmt.Errorf("write CLI config: %w", err)
	}

	return nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.APIURL) == "" {
		return fmt.Errorf("api_url is required")
	}

	if strings.TrimSpace(c.IdentityURL) == "" {
		return fmt.Errorf("identity_url is required")
	}

	if strings.TrimSpace(c.APIKey) == "" {
		return fmt.Errorf("api_key is required")
	}

	return nil
}

func (c Config) normalized() Config {
	c.APIURL = strings.TrimSpace(c.APIURL)
	c.IdentityURL = strings.TrimSpace(c.IdentityURL)
	c.APIKey = strings.TrimSpace(c.APIKey)

	return c
}
