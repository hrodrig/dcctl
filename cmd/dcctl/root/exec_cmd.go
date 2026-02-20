package root

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
)

func newExecCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "exec <service> <command> [args...]",
		Short: "Execute command in a service container",
		Long:  "Runs a command inside a running service container.",
		Example: `  dcctl exec app sh
  dcctl exec app ls -la /app`,
		Args: cobra.MinimumNArgs(2),
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
			compose, err := resolveComposeCommand()
			if err != nil {
				return err
			}
			projectName := loaded.Config.Common.Compose.ProjectName
			manifestPaths := resolveManifestPaths(loaded)
			if len(manifestPaths) == 0 {
				return errors.New("no manifests found")
			}
			if err := ensureExternalNetworks(manifestPaths); err != nil {
				return err
			}
			service := args[0]
			command := args[1:]
			log.Debug("executing in %s: %v", service, command)
			return compose.exec(ctx, manifestPaths, service, command, projectName)
		},
	}
}
