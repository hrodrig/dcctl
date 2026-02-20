# dcctl — Documentation

If `docs/demo.gif` is missing, generate it with: `vhs docs/demo.tape` (from repo root; requires [VHS](https://github.com/charmbracelet/vhs) and `dcctl` on PATH).

## Sequence diagrams

Sequence diagrams for main flows (Mermaid format). View in any Markdown viewer that supports Mermaid (e.g. GitHub, VS Code with Mermaid extension, or [Mermaid Live](https://mermaid.live)).

| Diagram | Description |
|---------|-------------|
| [01-config-and-up](./sequence/01-config-and-up.md) | Load config, resolve environment and manifests, run docker compose up |
| [02-down-and-cleanup](./sequence/02-down-and-cleanup.md) | Stop services, optional remove containers/volumes (down) |
| [03-manifest-resolution](./sequence/03-manifest-resolution.md) | How dcctl resolves config path, environment, and manifest file paths |
| [04-image](./sequence/04-image.md) | Image ls / pull / remove from current environment manifests |
| [05-show-ports](./sequence/05-show-ports.md) | Extract and list published ports per service (where to connect) |
| [06-manifests](./sequence/06-manifests.md) | Write sample compose manifest into environment directory |
| [07-volumes](./sequence/07-volumes.md) | Volumes clean: list (by project or all), select, confirm, remove |

## Terminal demo (VHS)

A [VHS](https://github.com/charmbracelet/vhs) tape records a short terminal demo of dcctl (help, version, config template).

### Prerequisites

- [VHS](https://github.com/charmbracelet/vhs): `brew install vhs` (or see project install docs)
- dcctl binary on `PATH` (e.g. `make build` then `export PATH="$PWD:$PATH"` from repo root)

### Render the demo

From the **repository root** (so `dcctl` resolves and paths match):

```bash
vhs docs/demo.tape
```

Output is written to `docs/demo.gif` (or the path set by `Output` in the tape). To produce MP4 instead, change the `Output` line in `demo.tape` to e.g. `Output docs/demo.mp4` and run again.

### Tape location

- Tape file: **`docs/demo.tape`**
- Rendered GIF (default): **`docs/demo.gif`**

### Prompt / shell issues

If you see a broken prompt in the GIF, run VHS from **bash** so it does not inherit a complex zsh setup:

```bash
bash
vhs docs/demo.tape
```
