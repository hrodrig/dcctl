package root

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

func newUpCmd(opts *options) *cobra.Command {
	var checkHealth bool
	cmd := &cobra.Command{
		Use:   "up [service...]",
		Short: "Start services",
		Long:  "Starts all services for the selected environment, or the given services/manifests.",
		Example: `  dcctl up
  dcctl up app
  dcctl up -e prod`,
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
			log.Info("Starting services...")
			compose, err := resolveComposeCommand()
			if err != nil {
				return err
			}
			if err := ensureEnvironmentDir(loaded); err != nil {
				return err
			}
			removeOrphans := !loaded.Config.Common.Compose.IgnoreOrphans
			projectName := loaded.Config.Common.Compose.ProjectName
			var manifestPaths []string
			if len(args) > 0 {
				manifestPaths, err = resolveTargetManifestPaths(loaded, args)
				if err != nil {
					return err
				}
				log.Info("targeting services: %s", strings.Join(args, ", "))
			} else {
				manifestPaths = resolveManifestPaths(loaded)
			}
			if len(manifestPaths) == 0 {
				log.Warn("no manifests found to start")
				return nil
			}
			if len(args) > 0 {
				if validatePortCollisionsForManifests(manifestPaths) {
					log.Error("port collisions detected; aborting")
					return nil
				}
			} else {
				env := loaded.Config.Dcctl.Environments[loaded.Environment]
				if validatePortCollisions(loaded.Path, loaded.Environment, env.Services) {
					log.Error("port collisions detected; aborting")
					return nil
				}
			}
			log.Verbose("starting %d manifests", len(manifestPaths))
			if err := ensureExternalNetworks(manifestPaths); err != nil {
				return err
			}
			if err := compose.up(ctx, manifestPaths, removeOrphans, projectName); err != nil {
				return err
			}
			if checkHealth {
				log.Info("Checking service health...")
				healthy, err := compose.healthCheck(ctx, manifestPaths, projectName)
				if err != nil {
					log.Warn("could not check health: %v", err)
				} else if healthy {
					log.Success("All services are healthy")
				} else {
					log.Warn("Some services are unhealthy")
				}
			}
			log.Success("Done.")
			return nil
		},
	}
	cmd.Flags().BoolVar(&checkHealth, "health-check", false, "Check service health after starting")
	return cmd
}
