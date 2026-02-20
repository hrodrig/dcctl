package root

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

func newDownCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "down [service...]",
		Short: "Stop services",
		Long:  "Stops all services for the selected environment, or the given services/manifests.",
		Example: `  dcctl down
  dcctl down app
  dcctl down -e prod`,
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
			log.Info("Stopping services...")
			if err := ensureEnvironmentDir(loaded); err != nil {
				return err
			}
			printSummary(loaded)
			compose, err := resolveComposeCommand()
			if err != nil {
				return err
			}
			removeOrphans := !loaded.Config.Common.Compose.IgnoreOrphans
			projectName := loaded.Config.Common.Compose.ProjectName
			var manifestPaths []string
			var targetServices []string
			if len(args) > 0 {
				manifestPaths, err = resolveTargetManifestPaths(loaded, args)
				if err != nil {
					return err
				}
				targetServices, err = resolveTargetServices(loaded, args, manifestPaths)
				if err != nil {
					return err
				}
				if len(targetServices) == 0 {
					log.Warn("no services found to stop")
					return nil
				}
				log.Info("targeting services: %s", strings.Join(targetServices, ", "))
			} else {
				manifestPaths = resolveManifestPaths(loaded)
			}
			if len(manifestPaths) == 0 {
				log.Warn("no manifests found to stop")
				return nil
			}
			if len(args) > 0 {
				log.Verbose("stopping %d services", len(targetServices))
			} else {
				log.Verbose("stopping %d manifests", len(manifestPaths))
			}
			if err := ensureExternalNetworks(manifestPaths); err != nil {
				return err
			}
			if len(args) > 0 {
				if err := compose.stop(ctx, manifestPaths, targetServices, projectName); err != nil {
					return err
				}
				if err := compose.rm(ctx, manifestPaths, targetServices, projectName); err != nil {
					return err
				}
			} else {
				if err := compose.down(ctx, manifestPaths, removeOrphans, projectName); err != nil {
					return err
				}
			}
			log.Success("Done.")
			return nil
		},
	}
}
