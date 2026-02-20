package config

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: Config{
				Dcctl: Dcctl{
					DefaultEnvironment: "default",
					Environments: map[string]Environment{
						"default": {Services: []string{"app", "db"}},
					},
				},
			},
			expectError: false,
		},
		{
			name: "empty environments",
			config: Config{
				Dcctl: Dcctl{
					Environments: map[string]Environment{},
				},
			},
			expectError: true,
			errorMsg:    "environments",
		},
		{
			name: "environment with no services",
			config: Config{
				Dcctl: Dcctl{
					Environments: map[string]Environment{
						"default": {Services: []string{}},
					},
				},
			},
			expectError: true,
			errorMsg:    "must have",
		},
		{
			name: "environment with empty service name",
			config: Config{
				Dcctl: Dcctl{
					Environments: map[string]Environment{
						"default": {Services: []string{"app", ""}},
					},
				},
			},
			expectError: true,
			errorMsg:    "empty service",
		},
		{
			name: "environment with duplicate services",
			config: Config{
				Dcctl: Dcctl{
					Environments: map[string]Environment{
						"default": {Services: []string{"app", "app"}},
					},
				},
			},
			expectError: true,
			errorMsg:    "duplicate",
		},
		{
			name: "invalid default environment",
			config: Config{
				Dcctl: Dcctl{
					DefaultEnvironment: "nonexistent",
					Environments: map[string]Environment{
						"default": {Services: []string{"app"}},
					},
				},
			},
			expectError: true,
			errorMsg:    "does not exist",
		},
		{
			name: "valid config without default_environment set",
			config: Config{
				Dcctl: Dcctl{
					Environments: map[string]Environment{
						"default": {Services: []string{"app"}},
					},
				},
			},
			expectError: false,
		},
		{
			name: "missing required default environment",
			config: Config{
				Dcctl: Dcctl{
					Environments: map[string]Environment{
						"uat": {Services: []string{"app"}},
					},
				},
			},
			expectError: true,
			errorMsg:    "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %q", err.Error())
				}
			}
		})
	}
}
