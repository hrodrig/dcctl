package root

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"dcctl/internal/config"
	"dcctl/internal/util"

	"gopkg.in/yaml.v3"
)

var log = util.Default()

func loadOrBootstrap(opts *options) (*config.LoadedConfig, error) {
	loaded, err := config.Load(config.LoadOptions{
		ConfigFile:  opts.configFile,
		Environment: opts.environment,
	})
	if err != nil {
		return nil, err
	}
	if loaded != nil {
		level := util.LogLevel(loaded.Config.Common.Dcctl.VerboseLevel)
		log.SetLevel(level)
	}
	if opts.debug {
		log.SetLevel(util.LevelDebug)
		log.Debug("debug logging enabled via --debug")
	}
	return loaded, nil
}

func printSummary(loaded *config.LoadedConfig) {
	env := loaded.Config.Dcctl.Environments[loaded.Environment]
	log.Info("config=%s env=%s", loaded.Path, loaded.Environment)
	log.Debug("services=%s", strings.Join(env.Services, ","))
}

func ensureEnvironmentDir(loaded *config.LoadedConfig) error {
	baseDir := filepath.Dir(loaded.Path)
	envDir := filepath.Join(baseDir, loaded.Environment)
	if !util.DirExists(envDir) {
		return fmt.Errorf(
			"environment directory not found: %s (config: %s)",
			envDir, loaded.Path,
		)
	}
	return nil
}

func resolveManifestPaths(loaded *config.LoadedConfig) []string {
	env := loaded.Config.Dcctl.Environments[loaded.Environment]
	var manifestPaths []string
	for _, service := range env.Services {
		p := manifestPathFor(loaded.Path, loaded.Environment, service)
		if !util.FileExists(p) {
			log.Warn("manifest not found for %s/%s", loaded.Environment, service)
			continue
		}
		log.Debug("found manifest: %s", p)
		manifestPaths = append(manifestPaths, p)
	}
	return manifestPaths
}

func manifestPathFor(configPath, environment, service string) string {
	baseDir := filepath.Dir(configPath)
	return filepath.Join(baseDir, environment, service+".yml")
}

func validatePortCollisions(configPath, environment string, services []string) bool {
	portMap := make(map[string][]string)
	hasCollision := false
	for _, service := range services {
		p := manifestPathFor(configPath, environment, service)
		if !util.FileExists(p) {
			continue
		}
		ports, err := extractPortsFromManifest(p)
		if err != nil {
			log.Warn("could not parse manifest for %s: %v", service, err)
			continue
		}
		for _, port := range ports {
			portMap[port] = append(portMap[port], service)
		}
	}
	for port, list := range portMap {
		if len(list) > 1 {
			hasCollision = true
			log.Warn("port collision: port %s used by: %s", port, strings.Join(list, ", "))
		}
	}
	return hasCollision
}

func extractPortsFromManifest(manifestPath string) ([]string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	var compose map[string]interface{}
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, err
	}
	services, ok := compose["services"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no services section found")
	}
	var ports []string
	for _, svc := range services {
		serviceMap, ok := svc.(map[string]interface{})
		if !ok {
			continue
		}
		portsList, ok := serviceMap["ports"].([]interface{})
		if !ok {
			continue
		}
		for _, entry := range portsList {
			if hp := extractHostPort(entry); hp != "" {
				ports = append(ports, hp)
			}
		}
	}
	return ports, nil
}

func extractHostPort(portEntry interface{}) string {
	if s, ok := portEntry.(string); ok {
		parts := strings.Split(s, ":")
		if len(parts) == 2 {
			return parts[0]
		}
		if len(parts) == 3 {
			return parts[1]
		}
		return ""
	}
	if m, ok := portEntry.(map[string]interface{}); ok {
		if p, ok := m["published"].(int); ok {
			return strconv.Itoa(p)
		}
		if p, ok := m["published"].(string); ok {
			return p
		}
	}
	return ""
}

func extractServiceNamesFromManifest(manifestPath string) ([]string, error) {
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
		return nil, fmt.Errorf("no services section found in %s", manifestPath)
	}
	names := make([]string, 0, len(services))
	for name := range services {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func asStringMap(input interface{}) (map[string]interface{}, bool) {
	switch t := input.(type) {
	case map[string]interface{}:
		return t, true
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, v := range t {
			if ks, ok := k.(string); ok {
				out[ks] = v
			}
		}
		return out, true
	default:
		return nil, false
	}
}

func resolveManifestPath(configPath, environment, manifestName string) string {
	baseDir := filepath.Join(filepath.Dir(configPath), environment)
	name := manifestName
	if !strings.HasSuffix(name, ".yml") {
		name += ".yml"
	}
	return filepath.Join(baseDir, name)
}

func resolveTargetManifestPaths(loaded *config.LoadedConfig, targets []string) ([]string, error) {
	baseDir := filepath.Join(filepath.Dir(loaded.Path), loaded.Environment)
	var paths []string
	seen := make(map[string]struct{})
	for _, target := range targets {
		target = strings.TrimSpace(target)
		if target == "" {
			continue
		}
		p := manifestPathFor(loaded.Path, loaded.Environment, target)
		if util.FileExists(p) {
			if _, ok := seen[p]; !ok {
				seen[p] = struct{}{}
				paths = append(paths, p)
			}
			continue
		}
		if strings.HasSuffix(target, ".yml") {
			p = filepath.Join(baseDir, target)
			if util.FileExists(p) {
				if _, ok := seen[p]; !ok {
					seen[p] = struct{}{}
					paths = append(paths, p)
				}
				continue
			}
		}
		manifestPath, err := findManifestByServiceName(baseDir, target)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[manifestPath]; !ok {
			seen[manifestPath] = struct{}{}
			paths = append(paths, manifestPath)
		}
	}
	return paths, nil
}

func findManifestByServiceName(baseDir, serviceName string) (string, error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to read environment dir: %w", err)
	}
	var matches []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}
		manifestPath := filepath.Join(baseDir, entry.Name())
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}
		var compose map[string]interface{}
		if err := yaml.Unmarshal(data, &compose); err != nil {
			continue
		}
		services, ok := asStringMap(compose["services"])
		if !ok {
			continue
		}
		if _, ok := services[serviceName]; ok {
			matches = append(matches, manifestPath)
		}
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("service %q not found in %s", serviceName, baseDir)
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("service %q found in multiple manifests: %s", serviceName, strings.Join(matches, ", "))
	}
	return matches[0], nil
}

func resolveTargetServices(loaded *config.LoadedConfig, args []string, manifestPaths []string) ([]string, error) {
	services := make(map[string]struct{})
	for _, arg := range args {
		name := strings.TrimSpace(arg)
		if name == "" {
			continue
		}
		p := resolveManifestPath(loaded.Path, loaded.Environment, name)
		if util.FileExists(p) {
			names, err := extractServiceNamesFromManifest(p)
			if err != nil {
				return nil, err
			}
			for _, svc := range names {
				services[svc] = struct{}{}
			}
			continue
		}
		services[name] = struct{}{}
	}
	if len(services) == 0 {
		for _, manifestPath := range manifestPaths {
			names, err := extractServiceNamesFromManifest(manifestPath)
			if err != nil {
				return nil, err
			}
			for _, svc := range names {
				services[svc] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(services))
	for svc := range services {
		out = append(out, svc)
	}
	sort.Strings(out)
	return out, nil
}

func validatePortCollisionsForManifests(manifestPaths []string) bool {
	portMap := make(map[string][]string)
	hasCollision := false
	for _, manifestPath := range manifestPaths {
		ports, err := extractPortsFromManifest(manifestPath)
		if err != nil {
			log.Warn("could not parse manifest %s: %v", manifestPath, err)
			continue
		}
		label := strings.TrimSuffix(filepath.Base(manifestPath), ".yml")
		for _, port := range ports {
			portMap[port] = append(portMap[port], label)
		}
	}
	for port, list := range portMap {
		if len(list) > 1 {
			hasCollision = true
			log.Warn("port collision: port %s used by: %s", port, strings.Join(list, ", "))
		}
	}
	return hasCollision
}

func ensureExternalNetworks(manifestPaths []string) error {
	networks := make(map[string]struct{})
	for _, manifestPath := range manifestPaths {
		names, err := extractExternalNetworks(manifestPath)
		if err != nil {
			log.Warn("could not parse networks for %s: %v", manifestPath, err)
			continue
		}
		for _, name := range names {
			networks[name] = struct{}{}
		}
	}
	for name := range networks {
		if networkExists(name) {
			log.Debug("network exists: %s", name)
			continue
		}
		log.Info("Creating network %s", name)
		if err := createNetwork(name); err != nil {
			return err
		}
	}
	return nil
}

func extractExternalNetworks(manifestPath string) ([]string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	var compose map[string]interface{}
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, err
	}
	networksSection, ok := compose["networks"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	var names []string
	for key, raw := range networksSection {
		network, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		external, ok := network["external"]
		if !ok {
			continue
		}
		externalEnabled := false
		var externalName string
		switch v := external.(type) {
		case bool:
			externalEnabled = v
		case map[string]interface{}:
			externalEnabled = true
			if n, ok := v["name"].(string); ok && n != "" {
				externalName = n
			}
		}
		if !externalEnabled {
			continue
		}
		if n, ok := network["name"].(string); ok && n != "" {
			externalName = n
		}
		if externalName == "" {
			externalName = key
		}
		names = append(names, externalName)
	}
	return names, nil
}

func networkExists(name string) bool {
	cmd := exec.Command("docker", "network", "inspect", name)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

func createNetwork(name string) error {
	cmd := exec.Command("docker", "network", "create", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("failed to create network %s: %s", name, msg)
	}
	log.Success("network created: %s", name)
	return nil
}
