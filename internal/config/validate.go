package config

import "fmt"

// Validate checks the configuration for consistency and required fields.
func Validate(cfg Config) error {
	if len(cfg.Dcctl.Environments) == 0 {
		return fmt.Errorf("dcctl.environments is required and must not be empty")
	}

	for name, env := range cfg.Dcctl.Environments {
		if len(env.Services) == 0 {
			return fmt.Errorf("environment %q must have at least one service", name)
		}
		seen := make(map[string]struct{}, len(env.Services))
		for _, service := range env.Services {
			if service == "" {
				return fmt.Errorf("environment %q has an empty service name", name)
			}
			if _, exists := seen[service]; exists {
				return fmt.Errorf("environment %q has duplicate service: %q", name, service)
			}
			seen[service] = struct{}{}
		}
	}

	if cfg.Dcctl.DefaultEnvironment != "" {
		if _, ok := cfg.Dcctl.Environments[cfg.Dcctl.DefaultEnvironment]; !ok {
			return fmt.Errorf("default_environment %q does not exist in environments", cfg.Dcctl.DefaultEnvironment)
		}
	}

	return nil
}
