package root

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newImagesCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "images <env[.target]> <ls|rm>",
		Short: "List or remove images used by manifests",
		Long:  "Lists or removes images referenced by manifests. Target: <environment> or <environment>.<service-or-manifest>.",
		Example: `  dcctl images default ls
  dcctl images default.app rm`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withSignalCancel(context.Background())
			defer cancel()
			if err := ensureDockerAvailable(); err != nil {
				return err
			}
			envName, target := parseImagesTarget(args[0])
			if envName == "" {
				return errors.New("missing environment (use <env> or <env>.<service>)")
			}
			action := strings.ToLower(strings.TrimSpace(args[1]))
			if action != "ls" && action != "rm" {
				return fmt.Errorf("unsupported action %q (use ls or rm)", action)
			}
			loaded, err := loadOrBootstrap(&options{
				environment: envName,
				configFile:  opts.configFile,
				debug:       opts.debug,
			})
			if err != nil {
				return err
			}
			if loaded == nil {
				return nil
			}
			if err := ensureEnvironmentDir(loaded); err != nil {
				return err
			}
			var manifestPaths []string
			if target != "" {
				manifestPaths, err = resolveTargetManifestPaths(loaded, []string{target})
				if err != nil {
					return err
				}
			} else {
				manifestPaths = resolveManifestPaths(loaded)
			}
			if len(manifestPaths) == 0 {
				log.Warn("no manifests found")
				return nil
			}
			images, err := collectImages(manifestPaths)
			if err != nil {
				return err
			}
			if len(images) == 0 {
				log.Warn("no images found")
				return nil
			}
			switch action {
			case "ls":
				printImages(images)
				return nil
			case "rm":
				printImages(images)
				if !confirmImageRemoval(images) {
					log.Warn("cancelled")
					return nil
				}
				return removeImages(ctx, images)
			}
			return nil
		},
	}
}

func parseImagesTarget(input string) (env, target string) {
	parts := strings.SplitN(strings.TrimSpace(input), ".", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func collectImages(manifestPaths []string) ([]string, error) {
	seen := make(map[string]struct{})
	for _, manifestPath := range manifestPaths {
		images, err := extractImagesFromManifest(manifestPath)
		if err != nil {
			return nil, err
		}
		for _, image := range images {
			seen[image] = struct{}{}
		}
	}
	return sortedKeys(seen), nil
}

func extractImagesFromManifest(manifestPath string) ([]string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	var compose map[string]interface{}
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, err
	}
	services, ok := asStringMap(compose["services"])
	if !ok {
		return nil, fmt.Errorf("no services section in %s", manifestPath)
	}
	images := make(map[string]struct{})
	for _, serviceRaw := range services {
		serviceMap, ok := asStringMap(serviceRaw)
		if !ok {
			continue
		}
		if image, _ := serviceMap["image"].(string); strings.TrimSpace(image) != "" {
			images[strings.TrimSpace(image)] = struct{}{}
		}
	}
	return sortedKeys(images), nil
}

func sortedKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func printImages(images []string) {
	fmt.Println("Images:")
	for _, name := range images {
		sha := imageID(name)
		if sha == "" {
			sha = "not found"
		}
		fmt.Printf("  - %s (%s)\n", name, sha)
	}
}

func imageID(image string) string {
	cmd := exec.Command("docker", "image", "inspect", "--format", "{{.Id}}", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func confirmImageRemoval(images []string) bool {
	fmt.Printf("Remove %d images? Type 'yes' to confirm: ", len(images))
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(line), "yes")
}

func removeImages(ctx context.Context, images []string) error {
	var failed []string
	for _, image := range images {
		if err := removeImage(ctx, image); err != nil {
			log.Warn("failed to remove image %s: %v", image, err)
			failed = append(failed, image)
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("failed to remove: %s", strings.Join(failed, ", "))
	}
	log.Success("Done.")
	return nil
}

func removeImage(ctx context.Context, image string) error {
	cmd := exec.CommandContext(ctx, "docker", "image", "rm", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("%s", msg)
	}
	log.Success("image removed: %s", image)
	return nil
}
