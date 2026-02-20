package root

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func newVolumesCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volumes",
		Short: "Manage Docker volumes",
	}
	cmd.AddCommand(newVolumesCleanCmd(opts))
	return cmd
}

func newVolumesCleanCmd(opts *options) *cobra.Command {
	var listAll bool
	var match string
	cmd := &cobra.Command{
		Use:   "clean [volume...]",
		Short: "Remove Docker volumes",
		Long:  "Removes Docker volumes. If volume names are given, only those; otherwise interactive selection.",
		Example: `  dcctl volumes clean
  dcctl volumes clean myvol1 myvol2
  dcctl volumes clean --match "dcctl_.*"`,
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
			projectName := loaded.Config.Common.Compose.ProjectName
			if len(args) > 0 {
				return removeVolumes(ctx, uniqueStrings(args))
			}
			volumes, err := listVolumesForSelection(projectName, listAll)
			if err != nil {
				return err
			}
			if len(volumes) == 0 {
				log.Warn("no volumes found")
				return nil
			}
			if match != "" {
				filtered, err := filterVolumes(volumes, match)
				if err != nil {
					return err
				}
				if len(filtered) == 0 {
					log.Warn("no volumes matched %q", match)
					return nil
				}
				if !confirmVolumeRemoval(filtered) {
					log.Warn("cancelled")
					return nil
				}
				return removeVolumes(ctx, filtered)
			}
			selected, err := promptForVolumes(volumes)
			if err != nil {
				return err
			}
			if len(selected) == 0 {
				log.Warn("no volumes selected")
				return nil
			}
			if !confirmVolumeRemoval(selected) {
				log.Warn("cancelled")
				return nil
			}
			return removeVolumes(ctx, selected)
		},
	}
	cmd.Flags().BoolVar(&listAll, "all", false, "List all volumes (ignore project filter)")
	cmd.Flags().StringVar(&match, "match", "", "Regex to match volume names")
	return cmd
}

func listVolumesForSelection(projectName string, listAll bool) ([]string, error) {
	if listAll || projectName == "" {
		return listVolumes("")
	}
	labelFilter := fmt.Sprintf("com.docker.compose.project=%s", projectName)
	volumes, err := listVolumes(labelFilter)
	if err != nil {
		return nil, err
	}
	if len(volumes) > 0 {
		return volumes, nil
	}
	log.Warn("no volumes for project %q, listing all", projectName)
	return listVolumes("")
}

func listVolumes(labelFilter string) ([]string, error) {
	args := []string{"volume", "ls", "--format", "{{.Name}}"}
	if labelFilter != "" {
		args = append(args, "--filter", "label="+labelFilter)
	}
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("failed to list volumes: %s", msg)
	}
	lines := splitLines(string(output))
	sort.Strings(lines)
	return lines, nil
}

func filterVolumes(volumes []string, pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %w", err)
	}
	var out []string
	for _, name := range volumes {
		if re.MatchString(name) {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out, nil
}

func promptForVolumes(volumes []string) ([]string, error) {
	fmt.Println("Available volumes:")
	for i, name := range volumes {
		fmt.Printf("  %2d) %s\n", i+1, name)
	}
	fmt.Print("Select by number or name (comma/space separated), or Enter to cancel: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}
	return parseVolumeSelection(line, volumes), nil
}

func parseVolumeSelection(input string, volumes []string) []string {
	repl := strings.NewReplacer(",", " ", ";", " ")
	clean := repl.Replace(input)
	tokens := strings.Fields(clean)
	indexed := make(map[int]struct{})
	selected := make(map[string]struct{})
	for _, token := range tokens {
		if strings.Contains(token, "-") {
			applyRange(token, len(volumes), indexed)
			continue
		}
		if idx, err := strconv.Atoi(token); err == nil && idx >= 1 && idx <= len(volumes) {
			indexed[idx-1] = struct{}{}
			continue
		}
		for _, name := range volumes {
			if name == token {
				selected[name] = struct{}{}
				break
			}
		}
	}
	for idx := range indexed {
		selected[volumes[idx]] = struct{}{}
	}
	return sortedKeys(selected)
}

func applyRange(token string, maxIdx int, indexed map[int]struct{}) {
	parts := strings.Split(token, "-")
	if len(parts) != 2 {
		return
	}
	start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || start <= 0 || end <= 0 {
		return
	}
	if start > end {
		start, end = end, start
	}
	for i := start - 1; i <= end-1 && i < maxIdx; i++ {
		indexed[i] = struct{}{}
	}
}

func confirmVolumeRemoval(volumes []string) bool {
	fmt.Printf("Remove %d volumes? Type 'yes' to confirm: ", len(volumes))
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(line), "yes")
}

func removeVolumes(ctx context.Context, volumes []string) error {
	var failed []string
	for _, name := range volumes {
		if !volumeExists(name) {
			log.Warn("volume not found: %s", name)
			continue
		}
		if err := removeVolume(ctx, name); err != nil {
			log.Warn("failed to remove volume %s: %v", name, err)
			failed = append(failed, name)
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("failed to remove: %s", strings.Join(failed, ", "))
	}
	log.Success("Done.")
	return nil
}

func volumeExists(name string) bool {
	cmd := exec.Command("docker", "volume", "inspect", name)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

func removeVolume(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "docker", "volume", "rm", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("%s", msg)
	}
	log.Success("volume removed: %s", name)
	return nil
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	for _, s := range input {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		seen[s] = struct{}{}
	}
	return sortedKeys(seen)
}
