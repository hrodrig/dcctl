# Example 01: Counter + Redis

Minimal stack: a small Node.js web app that stores a counter in Redis. Classic "first steps" example for Docker Compose and dcctl.

## Stack

- **Redis** — in-memory store for the counter
- **Web** — Express app that increments and displays the count on each page load

## Run with dcctl

From the **dcctl repo root** (so `dcctl` is on PATH):

```bash
dcctl --config-file=examples/01-counter-redis/dcctl.yml up
```

Then open http://localhost:3000 and refresh to see the counter increment.

From this example directory you can also run (if dcctl is on PATH):

```bash
dcctl --config-file=dcctl.yml -e default up
```

## Useful commands

```bash
dcctl --config-file=examples/01-counter-redis/dcctl.yml status
dcctl --config-file=examples/01-counter-redis/dcctl.yml logs -f
dcctl --config-file=examples/01-counter-redis/dcctl.yml down
```

## Files

- `dcctl.yml` — dcctl config (one environment, one manifest `stack`)
- `default/stack.yml` — Compose: redis + web
- `Dockerfile`, `package.json`, `server.js` — counter app
