package config

import (
	"embed"
	"fmt"
)

//go:embed assets/dcctl.yml
var templateFS embed.FS

const templatePath = "assets/dcctl.yml"

// DefaultTemplate returns the default dcctl.yml template bytes.
func DefaultTemplate() ([]byte, error) {
	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config template: %w", err)
	}
	return data, nil
}
