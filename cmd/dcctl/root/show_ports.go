package root

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type portMapping struct {
	hostPort string
	protocol string
}

type serviceEndpoint struct {
	service  string
	endpoint string
	scheme   string
}

func newShowPortsCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "show-ports",
		Short: "List published ports for services",
		Long:  "Scans manifest files for published ports and prints service endpoints.",
		Example: `  dcctl show-ports
  dcctl show-ports -e prod`,
		RunE: func(cmd *cobra.Command, args []string) error {
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
			manifestPaths := resolveManifestPaths(loaded)
			if len(manifestPaths) == 0 {
				log.Warn("no manifests found for environment %s", loaded.Environment)
				return nil
			}
			endpoints, err := collectServiceEndpoints(manifestPaths)
			if err != nil {
				return err
			}
			if len(endpoints) == 0 {
				log.Info("No published ports found.")
				return nil
			}
			sort.Slice(endpoints, func(i, j int) bool {
				if endpoints[i].service != endpoints[j].service {
					return endpoints[i].service < endpoints[j].service
				}
				return endpoints[i].endpoint < endpoints[j].endpoint
			})
			for _, e := range endpoints {
				fmt.Printf("%s - %s\n", e.service, e.endpoint)
			}
			return nil
		},
	}
}

func collectServiceEndpoints(manifestPaths []string) ([]serviceEndpoint, error) {
	var endpoints []serviceEndpoint
	for _, manifestPath := range manifestPaths {
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
			continue
		}
		for name, raw := range services {
			serviceMap, ok := asStringMap(raw)
			if !ok {
				continue
			}
			portsList, ok := serviceMap["ports"].([]interface{})
			if !ok {
				continue
			}
			for _, entry := range portsList {
				mapping, ok := parsePortMapping(entry)
				if !ok || mapping.hostPort == "" {
					continue
				}
				scheme, endpoint := formatEndpoint(mapping)
				endpoints = append(endpoints, serviceEndpoint{service: name, endpoint: endpoint, scheme: scheme})
			}
		}
	}
	return endpoints, nil
}

func parsePortMapping(entry interface{}) (portMapping, bool) {
	switch v := entry.(type) {
	case string:
		return parsePortString(v)
	case map[string]interface{}:
		return parsePortMap(v)
	default:
		return portMapping{}, false
	}
}

func parsePortString(spec string) (portMapping, bool) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return portMapping{}, false
	}
	protocol := ""
	if strings.Contains(spec, "/") {
		parts := strings.SplitN(spec, "/", 2)
		spec, protocol = parts[0], parts[1]
	}
	parts := strings.Split(spec, ":")
	if len(parts) == 2 {
		return portMapping{hostPort: parts[0], protocol: protocol}, true
	}
	if len(parts) == 3 {
		return portMapping{hostPort: parts[1], protocol: protocol}, true
	}
	return portMapping{}, false
}

func parsePortMap(m map[string]interface{}) (portMapping, bool) {
	var hostPort string
	if p, ok := m["published"].(int); ok {
		hostPort = strconv.Itoa(p)
	} else if p, ok := m["published"].(string); ok {
		hostPort = p
	}
	if hostPort == "" {
		return portMapping{}, false
	}
	protocol, _ := m["protocol"].(string)
	return portMapping{hostPort: hostPort, protocol: protocol}, true
}

func formatEndpoint(m portMapping) (scheme, endpoint string) {
	scheme = inferScheme(m)
	if scheme == "" {
		scheme = "tcp"
	}
	return scheme, fmt.Sprintf("%s://localhost:%s", scheme, m.hostPort)
}

func inferScheme(m portMapping) string {
	switch m.hostPort {
	case "80", "8080", "3000", "8000":
		return "http"
	case "443":
		return "https"
	case "5432":
		return "postgres"
	case "6379":
		return "redis"
	case "27017":
		return "mongodb"
	case "5672":
		return "amqp"
	}
	if strings.EqualFold(m.protocol, "udp") {
		return "udp"
	}
	return ""
}
