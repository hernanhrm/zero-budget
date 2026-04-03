package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreSaveAndLoad(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "saves and loads config",
			cfg: Config{
				APIURL:      "http://localhost:8080",
				IdentityURL: "http://localhost:8081",
				APIKey:      "test-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(t.TempDir(), "zero-budget", "cli.json")

			store, err := NewStore(configPath)
			if err != nil {
				t.Fatalf("NewStore() error = %v", err)
			}

			if err := store.Save(tt.cfg); err != nil {
				t.Fatalf("Save() error = %v", err)
			}

			got, err := store.Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			if got != tt.cfg {
				t.Fatalf("Load() = %#v, want %#v", got, tt.cfg)
			}

			fileInfo, err := os.Stat(configPath)
			if err != nil {
				t.Fatalf("os.Stat() error = %v", err)
			}

			if gotPerms := fileInfo.Mode().Perm(); gotPerms != 0o600 {
				t.Fatalf("file permissions = %#o, want %#o", gotPerms, 0o600)
			}

			dirInfo, err := os.Stat(filepath.Dir(configPath))
			if err != nil {
				t.Fatalf("os.Stat(dir) error = %v", err)
			}

			if gotPerms := dirInfo.Mode().Perm(); gotPerms != 0o700 {
				t.Fatalf("dir permissions = %#o, want %#o", gotPerms, 0o700)
			}
		})
	}
}

func TestStoreLoadMissingConfig(t *testing.T) {
	store, err := NewStore(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	_, err = store.Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
}
