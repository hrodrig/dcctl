package root

import (
	"context"

	"github.com/spf13/cobra"
)

func newStatusCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show status",
		Long:  "Shows container status for all services in the selected environment.",
		Example: `  dcctl status
  dcctl status -e prod`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withSignalCancel(context.Background())
			defer cancel()
			if err := ensureDockerAvailable(); err != nil {
				return err
			}
			loaded, err := loadOrBootstrap(opts)
			if err != nil {
				return err
			}
			if loaded == nil {
				return nil
			}
			if err := ensureEnvironmentDir(loaded); err != nil {
				return err
			}
			printSummary(loaded)
			compose, err := resolveComposeCommand()
			if err != nil {
				return err
			}
			projectName := loaded.Config.Common.Compose.ProjectName
			manifestPaths := resolveManifestPaths(loaded)
			if len(manifestPaths) == 0 {
				log.Warn("no manifests found for status")
				return nil
			}
			if err := ensureExternalNetworks(manifestPaths); err != nil {
				return err
			}
			return compose.status(ctx, manifestPaths, projectName)
		},
	}
}
