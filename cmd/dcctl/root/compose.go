package root

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type composeCmd struct {
	bin      string
	baseArgs []string
}

func resolveComposeCommand() (composeCmd, error) {
	if err := ensureDockerAvailable(); err != nil {
		return composeCmd{}, err
	}
	if _, err := exec.LookPath("docker"); err == nil && isDockerComposeAvailable() {
		log.Debug("using docker compose")
		return composeCmd{bin: "docker", baseArgs: []string{"compose"}}, nil
	}
	if _, err := exec.LookPath("docker-compose"); err == nil {
		log.Debug("using docker-compose")
		return composeCmd{bin: "docker-compose"}, nil
	}
	return composeCmd{}, errors.New("docker compose not found (install Docker or docker-compose)")
}

func ensureDockerAvailable() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return errors.New("docker binary not found in PATH")
	}
	cmd := exec.Command("docker", "info")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return errors.New("docker is not available (is the daemon running?)")
	}
	return nil
}

func isDockerComposeAvailable() bool {
	cmd := exec.Command("docker", "compose", "version")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

func composeEnv(projectName string) []string {
	env := os.Environ()
	if projectName != "" {
		env = append(env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", projectName))
	}
	return env
}

func (c composeCmd) up(ctx context.Context, manifestPaths []string, removeOrphans bool, projectName string) error {
	args := append([]string{}, c.baseArgs...)
	for _, path := range manifestPaths {
		args = append(args, "-f", path)
	}
	args = append(args, "up", "-d")
	if removeOrphans {
		args = append(args, "--remove-orphans")
	}
	log.Debug("compose up: %s %s", c.bin, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = append(composeEnv(projectName),
		fmt.Sprintf("COMPOSE_REMOVE_ORPHANS=%t", removeOrphans),
		fmt.Sprintf("COMPOSE_IGNORE_ORPHANS=%t", !removeOrphans),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("compose up failed: %w", err)
	}
	return nil
}

func (c composeCmd) down(ctx context.Context, manifestPaths []string, removeOrphans bool, projectName string) error {
	args := append([]string{}, c.baseArgs...)
	for _, path := range manifestPaths {
		args = append(args, "-f", path)
	}
	args = append(args, "down")
	if removeOrphans {
		args = append(args, "--remove-orphans")
	}
	log.Debug("compose down: %s %s", c.bin, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = append(composeEnv(projectName),
		fmt.Sprintf("COMPOSE_REMOVE_ORPHANS=%t", removeOrphans),
		fmt.Sprintf("COMPOSE_IGNORE_ORPHANS=%t", !removeOrphans),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("compose down failed: %w", err)
	}
	return nil
}

func (c composeCmd) stop(ctx context.Context, manifestPaths []string, services []string, projectName string) error {
	args := append([]string{}, c.baseArgs...)
	for _, path := range manifestPaths {
		args = append(args, "-f", path)
	}
	args = append(args, "stop")
	args = append(args, services...)
	log.Debug("compose stop: %s %s", c.bin, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = append(composeEnv(projectName), "COMPOSE_IGNORE_ORPHANS=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("compose stop failed: %w", err)
	}
	return nil
}

func (c composeCmd) rm(ctx context.Context, manifestPaths []string, services []string, projectName string) error {
	args := append([]string{}, c.baseArgs...)
	for _, path := range manifestPaths {
		args = append(args, "-f", path)
	}
	args = append(args, "rm", "-f")
	args = append(args, services...)
	log.Debug("compose rm: %s %s", c.bin, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = append(composeEnv(projectName), "COMPOSE_IGNORE_ORPHANS=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("compose rm failed: %w", err)
	}
	return nil
}

func (c composeCmd) status(ctx context.Context, manifestPaths []string, projectName string) error {
	args := append([]string{}, c.baseArgs...)
	for _, path := range manifestPaths {
		args = append(args, "-f", path)
	}
	args = append(args, "ps")
	log.Debug("compose ps: %s %s", c.bin, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = append(composeEnv(projectName), "COMPOSE_IGNORE_ORPHANS=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("compose ps failed: %w", err)
	}
	return nil
}

func (c composeCmd) restart(ctx context.Context, manifestPaths []string, projectName string) error {
	args := append([]string{}, c.baseArgs...)
	for _, path := range manifestPaths {
		args = append(args, "-f", path)
	}
	args = append(args, "restart")
	log.Debug("compose restart: %s %s", c.bin, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = append(composeEnv(projectName), "COMPOSE_IGNORE_ORPHANS=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("compose restart failed: %w", err)
	}
	return nil
}

func (c composeCmd) logs(ctx context.Context, manifestPaths []string, service string, follow bool, tail string, projectName string) error {
	args := append([]string{}, c.baseArgs...)
	for _, path := range manifestPaths {
		args = append(args, "-f", path)
	}
	args = append(args, "logs")
	if follow {
		args = append(args, "-f")
	}
	if tail != "" {
		args = append(args, "--tail", tail)
	}
	if service != "" {
		args = append(args, service)
	}
	log.Debug("compose logs: %s %s", c.bin, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = append(composeEnv(projectName), "COMPOSE_IGNORE_ORPHANS=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return nil
		}
		return fmt.Errorf("compose logs failed: %w", err)
	}
	return nil
}

func (c composeCmd) exec(ctx context.Context, manifestPaths []string, service string, command []string, projectName string) error {
	args := append([]string{}, c.baseArgs...)
	for _, path := range manifestPaths {
		args = append(args, "-f", path)
	}
	args = append(args, "exec", service)
	args = append(args, command...)
	log.Debug("compose exec: %s %s", c.bin, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = append(composeEnv(projectName), "COMPOSE_IGNORE_ORPHANS=true")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("operation cancelled")
		}
		return fmt.Errorf("compose exec failed: %w", err)
	}
	return nil
}

func (c composeCmd) healthCheck(ctx context.Context, manifestPaths []string, projectName string) (bool, error) {
	args := append([]string{}, c.baseArgs...)
	for _, path := range manifestPaths {
		args = append(args, "-f", path)
	}
	args = append(args, "ps", "--format", "{{.Health}}")
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = composeEnv(projectName)
	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.Canceled {
			return false, fmt.Errorf("operation cancelled")
		}
		return false, err
	}
	for _, line := range splitLines(string(output)) {
		if line == "unhealthy" {
			return false, nil
		}
	}
	return true, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			if line != "" {
				lines = append(lines, line)
			}
			start = i + 1
		}
	}
	if start < len(s) {
		line := s[start:]
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
