package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type options struct {
	environment string
	configFile  string
	debug       bool
}

// Execute runs the root command.
func Execute() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	opts := &options{}

	cmd := &cobra.Command{
		Use:   "dcctl",
		Short: "Docker Compose Control",
		Long: `dcctl manages Docker Compose stacks per environment using a single config file.

Use -e/--environment to select an environment (defaults to dcctl.default_environment).`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       resolvedVersion(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.environment, "environment", "e", "", "Environment name (default: dcctl.default_environment)")
	cmd.PersistentFlags().StringVar(&opts.configFile, "config-file", "", "Path to dcctl.yml")
	cmd.PersistentFlags().BoolVar(&opts.debug, "debug", false, "Enable debug logging")

	cmd.AddCommand(
		newConfigCmd(),
		newUpCmd(opts),
		newDownCmd(opts),
		newRestartCmd(opts),
		newStatusCmd(opts),
		newLogsCmd(opts),
		newExecCmd(opts),
		newManifestsCmd(opts),
		newImagesCmd(opts),
		newShowPortsCmd(opts),
		newVolumesCmd(opts),
		newVersionCmd(),
		newCompletionCmd(),
	)

	return cmd
}
