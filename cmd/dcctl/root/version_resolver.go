package root

import (
	"os"
	"path/filepath"
	"strings"
)

func resolvedVersion() string {
	if buildVersion != "" && buildVersion != "dev" {
		return buildVersion
	}
	if v := readVersionFile("version"); v != "" {
		return v
	}
	if v := readVersionFile(filepath.Join("..", "version")); v != "" {
		return v
	}
	return buildVersion
}

func readVersionFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
