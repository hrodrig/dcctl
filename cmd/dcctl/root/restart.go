package root

import (
	"context"

	"github.com/spf13/cobra"
)

func newRestartCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart services",
		Long:  "Restarts all services for the selected environment.",
		Example: `  dcctl restart
  dcctl restart -e prod`,
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
			log.Info("Restarting services...")
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
				log.Warn("no manifests found to restart")
				return nil
			}
			log.Verbose("restarting %d manifests", len(manifestPaths))
			if err := ensureExternalNetworks(manifestPaths); err != nil {
				return err
			}
			if err := compose.restart(ctx, manifestPaths, projectName); err != nil {
				return err
			}
			log.Success("Done.")
			return nil
		},
	}
}
