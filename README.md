# dcctl — Docker Compose Control

[![version](https://img.shields.io/badge/version-0.1.0-blue)](https://github.com/hrodrig/dcctl/releases)
[![release](https://img.shields.io/github/v/release/hrodrig/dcctl)](https://github.com/hrodrig/dcctl/releases)
[![Go 1.26](https://img.shields.io/badge/Go-1.26-00ADD8)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

CLI to manage Docker Compose stacks per **environment** using a single YAML config file.

**Documentation:** [Sequence diagrams](docs/README.md#sequence-diagrams) (Mermaid) for main flows and [terminal demo](docs/README.md#terminal-demo-vhs) (recorded with [VHS](https://github.com/charmbracelet/vhs)) — see [docs/](docs/README.md).

![Terminal demo](docs/demo.gif)

## Requirements

- Go 1.26+ (to build)
- Docker and Docker Compose (to run stacks)

## Install

```bash
go build -o dcctl ./cmd/dcctl
# or with version info:
# make build
# optional: move to PATH
# sudo mv dcctl /usr/local/bin/
```

## Quick start

1. Create a config file:

```bash
mkdir -p ~/.dcctl
dcctl config > ~/.dcctl/dcctl.yml
```

2. Edit `~/.dcctl/dcctl.yml`: set `dcctl.default_environment` and `dcctl.environments.<name>.services` (list of manifest names).

3. Create compose manifests next to the config:

- Path pattern: `<config_dir>/<environment>/<service>.yml`
- Example: `~/.dcctl/default/app.yml` for environment `default` and service `app`.
- Or run `dcctl manifests` to write a sample `sample-app.yml` (rename to `app.yml` and add `app` to services if needed).

4. Start services:

```bash
dcctl up
dcctl status
dcctl logs -f
dcctl down
```

## Config file

- **Default location:** `~/.dcctl/dcctl.yml`
- **Override:** `dcctl --config-file=/path/to/dcctl.yml`

Schema (see `dcctl config` output):

- `schema_version`: config format version
- `common.compose.project_name`: Docker Compose project name
- `common.compose.ignore_orphans`: whether to remove orphan containers
- `common.dcctl.verbose_level`: 0=quiet, 1=normal, 2=verbose, 3=debug
- `dcctl.default_environment`: default environment name
- `dcctl.environments`: map of environment name → list of service (manifest) names

## Commands

| Command | Description |
|--------|--------------|
| `dcctl config` | Print default dcctl.yml template |
| `dcctl up [service...]` | Start services (optional `--health-check`) |
| `dcctl down [service...]` | Stop and remove services |
| `dcctl restart` | Restart all services |
| `dcctl status` | Show container status (compose ps) |
| `dcctl logs [service]` | View logs (`-f` follow, `--tail N`) |
| `dcctl exec <service> <cmd> [args...]` | Run command in container |
| `dcctl manifests` | Write sample compose file to environment dir |
| `dcctl images <env[.target]> ls\|rm` | List or remove images used by manifests |
| `dcctl show-ports` | List published ports per service |
| `dcctl volumes clean [volume...]` | Remove volumes (interactive or by name, `--match` regex) |
| `dcctl version` | Show version (`-s` short, `-o json\|yaml\|text`) |
| `dcctl completion bash\|zsh\|fish\|powershell` | Shell completion |

Use `-e` / `--environment` to select an environment; `--config-file` to use a different config file; `--debug` for debug logging.

## Build and release

**Local build (version from `VERSION` file):**

```bash
make build
# override: make build VERSION=v0.2.0
```

**Docker image (local):**

```bash
make docker-build
# image tagged as dcctl, with version/commit/build date from VERSION and git
```

**Release (from `main`, requires [goreleaser](https://goreleaser.com)):**

```bash
brew install goreleaser   # or see goreleaser install docs
git tag v0.1.0
make release
```

- Builds binaries for linux/darwin/windows (amd64, arm64), archives, checksums, .deb/.rpm, Docker image (ghcr.io), and Homebrew tap.
- Create a Homebrew tap repo (e.g. `homebrew-dcctl`) and set `HOMEBREW_TAP_TOKEN` in CI for `brew install hrodrig/dcctl/dcctl`.

**Snapshot (no tag):** `make snapshot` — outputs to `dist/`.

## License

MIT. See [LICENSE](LICENSE).
