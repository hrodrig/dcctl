package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dcctl.yml")
	configContent := `---
schema_version: "1.0"
common:
  dcctl:
    verbose_level: 2
  compose:
    ignore_orphans: true
    project_name: "test"
dcctl:
  default_environment: "default"
  environments:
    default:
      services:
        - app
        - db
    prod:
      services:
        - app
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	t.Run("load with explicit config file", func(t *testing.T) {
		loaded, err := Load(LoadOptions{ConfigFile: configPath})
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if loaded.Path != configPath {
			t.Errorf("Path = %q, want %q", loaded.Path, configPath)
		}
		if loaded.Environment != "default" {
			t.Errorf("Environment = %q, want default", loaded.Environment)
		}
		services := loaded.Config.Dcctl.Environments["default"].Services
		if len(services) != 2 {
			t.Errorf("expected 2 services, got %d", len(services))
		}
	})

	t.Run("load with specific environment", func(t *testing.T) {
		loaded, err := Load(LoadOptions{
			ConfigFile:  configPath,
			Environment: "prod",
		})
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if loaded.Environment != "prod" {
			t.Errorf("Environment = %q, want prod", loaded.Environment)
		}
	})

	t.Run("missing config file", func(t *testing.T) {
		_, err := Load(LoadOptions{ConfigFile: filepath.Join(tmpDir, "nonexistent.yml")})
		if err == nil {
			t.Fatal("expected error for missing file")
		}
		if err != nil && !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("expected 'does not exist' in error, got: %v", err)
		}
	})

	t.Run("invalid environment", func(t *testing.T) {
		_, err := Load(LoadOptions{
			ConfigFile:  configPath,
			Environment: "nonexistent",
		})
		if err == nil {
			t.Fatal("expected error for nonexistent environment")
		}
	})
}

