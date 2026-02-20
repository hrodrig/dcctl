package root

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"dcctl/internal/config"

	"github.com/spf13/cobra"
)

//go:embed assets/sample_app.yml
var sampleAppYAML []byte

func newManifestsCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "manifests",
		Short: "Write sample compose manifest",
		Long:  "Writes a sample app.yml compose file into the environment directory.",
		Example: `  dcctl manifests
  dcctl manifests -e default`,
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := config.Load(config.LoadOptions{
				ConfigFile:  opts.configFile,
				Environment: opts.environment,
			})
			if err != nil {
				return err
			}
			env := loaded.Environment
			if env == "" {
				return fmt.Errorf("environment could not be resolved")
			}
			baseDir := filepath.Dir(loaded.Path)
			targetDir := filepath.Join(baseDir, env)
			if err := os.MkdirAll(targetDir, 0o755); err != nil {
				return fmt.Errorf("failed to create environment directory: %w", err)
			}
			targetPath := filepath.Join(targetDir, "sample-app.yml")
			if err := os.WriteFile(targetPath, sampleAppYAML, 0o644); err != nil {
				return fmt.Errorf("failed to write %s: %w", targetPath, err)
			}
			fmt.Printf("Sample manifest written to %s\n", targetPath)
			return nil
		},
	}
}
