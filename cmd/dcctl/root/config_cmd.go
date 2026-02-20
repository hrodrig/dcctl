package root

import (
	"fmt"
	"os"

	"dcctl/internal/config"

	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print the default dcctl.yml template",
		Long: `Print the default dcctl.yml template to stdout.

Use shell redirection to write it to a file.`,
		Example: `  dcctl config > ~/.dcctl/dcctl.yml
  dcctl config > ./dcctl.yml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := config.DefaultTemplate()
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(os.Stdout, string(data))
			return err
		},
	}
}
