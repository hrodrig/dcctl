# Example 09: Portainer

Portainer CE — web UI to manage Docker. Single service; shows **environment variables** for port and timezone.

**Source:** [docker/awesome-compose — portainer](https://github.com/docker/awesome-compose/tree/master/portainer).

## Environment variables and dependencies

| Service   | Env vars | Depends on |
|-----------|----------|------------|
| portainer | `PORTAINER_PORT` (host port), `TZ` (optional) | — (uses Docker socket) |

No other services; Portainer connects to the host Docker daemon via the mounted socket.

## Run with dcctl

```bash
dcctl --config-file=dcctl.yml up
# → http://localhost:9000 — create admin user on first run
```

## Useful commands

```bash
dcctl --config-file=dcctl.yml status
dcctl --config-file=dcctl.yml logs -f portainer
dcctl --config-file=dcctl.yml down
```

## Files

- `dcctl.yml` — dcctl config
- `default/stack.yml` — Compose: portainer (env, volumes)
- `.env.example` — optional
