package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"dcctl/internal/util"

	"github.com/spf13/viper"
)

const defaultConfigName = "dcctl.yml"

// LoadOptions configures how the config is loaded.
type LoadOptions struct {
	ConfigFile  string
	Environment string
}

// LoadedConfig holds the resolved config path, selected environment, and parsed config.
type LoadedConfig struct {
	Path        string
	Environment string
	Config      Config
}

// Load reads and validates the config, resolving the selected environment.
func Load(options LoadOptions) (*LoadedConfig, error) {
	configPath, err := resolveConfigPath(options.ConfigFile)
	if err != nil {
		return nil, err
	}

	cfg, err := loadWithViper(configPath)
	if err != nil {
		return nil, err
	}

	if err := Validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	env := options.Environment
	if env == "" {
		env = cfg.Dcctl.DefaultEnvironment
		if env == "" {
			env = "default"
		}
	}
	if _, ok := cfg.Dcctl.Environments[env]; !ok {
		return nil, fmt.Errorf("environment %q does not exist in %s", env, configPath)
	}

	return &LoadedConfig{
		Path:        configPath,
		Environment: env,
		Config:      cfg,
	}, nil
}

func resolveConfigPath(configFile string) (string, error) {
	if configFile != "" {
		if !util.FileExists(configFile) {
			return "", fmt.Errorf("config file %q does not exist", configFile)
		}
		return configFile, nil
	}

	homePath, err := homeConfigPath()
	if err != nil {
		return "", err
	}
	if util.FileExists(homePath) {
		return homePath, nil
	}

	return "", errors.New(
		"no dcctl.yml config file found\n" +
			"  use a config file:\n" +
			"    dcctl --config-file=/path/to/dcctl.yml\n" +
			"  or create one at the default location:\n" +
			"    mkdir -p ~/.dcctl\n" +
			"    dcctl config > ~/.dcctl/dcctl.yml\n" +
			"  or at a custom path:\n" +
			"    mkdir -p /path/to/dir\n" +
			"    dcctl config > /path/to/dir/dcctl.yml",
	)
}

func homeConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}
	return filepath.Join(homeDir, ".dcctl", defaultConfigName), nil
}

func loadWithViper(configPath string) (Config, error) {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config %q: %w", configPath, err)
	}

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewBuffer(raw)); err != nil {
		return Config{}, fmt.Errorf("failed to parse config %q: %w", configPath, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to decode config %q: %w", configPath, err)
	}

	return cfg, nil
}
