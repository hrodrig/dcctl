package root

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

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

// pullImages runs docker pull for each image (ensures you have the latest).
func pullImages(ctx context.Context, images []string) error {
	for _, image := range images {
		if err := pullImage(ctx, image); err != nil {
			log.Warn("failed to pull %s: %v", image, err)
			continue
		}
		log.Success("pulled: %s", image)
	}
	return nil
}

func pullImage(ctx context.Context, image string) error {
	cmd := exec.CommandContext(ctx, "docker", "image", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
