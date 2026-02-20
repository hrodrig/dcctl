package root

import (
	"context"

	"github.com/spf13/cobra"
)

func newLogsCmd(opts *options) *cobra.Command {
	var follow bool
	var tail string
	cmd := &cobra.Command{
		Use:   "logs [service]",
		Short: "View service logs",
		Long:  "Shows logs for services. If a service name is given, only that service; otherwise all. Use -f to follow.",
		Example: `  dcctl logs
  dcctl logs app
  dcctl logs -f app
  dcctl logs --tail 100`,
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
				log.Warn("no manifests found")
				return nil
			}
			if err := ensureExternalNetworks(manifestPaths); err != nil {
				return err
			}
			var service string
			if len(args) > 0 {
				service = args[0]
				log.Verbose("showing logs for service: %s", service)
			}
			return compose.logs(ctx, manifestPaths, service, follow, tail, projectName)
		},
	}
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	cmd.Flags().StringVar(&tail, "tail", "", "Number of lines to show from the end (e.g. 100)")
	return cmd
}
