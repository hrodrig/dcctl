# Example 08: Nginx + Go backend

Nginx reverse proxy in front of a Go (Chi) backend. Shows **depends_on** (proxy waits for backend) and **environment variables** (`BACKEND_PORT`, `HTTP_PORT`).

**Source:** [docker/awesome-compose — nginx-golang](https://github.com/docker/awesome-compose/tree/master/nginx-golang).

## Environment variables and dependencies

| Service | Env vars | Depends on |
|---------|----------|------------|
| backend | `BACKEND_PORT` (default 80) | — |
| proxy   | — (host port via `HTTP_PORT`) | backend |

## Run with dcctl

```bash
dcctl --config-file=dcctl.yml up
# → http://localhost:80 (or HTTP_PORT) — Nginx proxies to Go backend
```

## Useful commands

```bash
dcctl --config-file=dcctl.yml status
dcctl --config-file=dcctl.yml logs -f backend
dcctl --config-file=dcctl.yml down
```

## Files

- `dcctl.yml` — dcctl config
- `default/stack.yml` — Compose: backend + proxy (depends_on)
- `backend/` — Go app (Chi), Dockerfile
- `proxy/nginx.conf` — Nginx config
- `.env.example` — optional
