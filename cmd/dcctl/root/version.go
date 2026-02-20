package root

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

type versionOutput struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func newVersionCmd() *cobra.Command {
	var short bool
	var output string
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show dcctl version",
		Long:  "Shows version, commit, and build date. Use -s for version only, -o for format.",
		Example: `  dcctl version
  dcctl version -s
  dcctl version -o text`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if short {
				fmt.Println(buildVersion)
				return nil
			}
			if output == "json" {
				data, err := json.Marshal(versionOutput{
					Version: buildVersion,
					Commit:  buildCommit,
					Date:    buildDate,
				})
				if err != nil {
					return fmt.Errorf("failed to encode version: %w", err)
				}
				fmt.Println(string(data))
				return nil
			}
			if output == "yaml" {
				fmt.Printf("Version: %s\nCommit: %s\nDate: %s\n", buildVersion, buildCommit, buildDate)
				return nil
			}
			fmt.Printf("version=%s commit=%s date=%s\n", buildVersion, buildCommit, buildDate)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&short, "short", "s", false, "Show only version number")
	cmd.Flags().StringVarP(&output, "output", "o", "json", "Output format (json, yaml, text)")
	return cmd
}
