# Example 04: Full-stack (multiple manifests)

"Enredado" example: several Compose files combined by dcctl — **databases** (Postgres, Redis, Valkey, Mongo), **api** (Node), and **frontend** (nginx). Shows how dcctl merges multiple `-f` manifests for one environment.

## Stack

- **databases.yml** — Postgres, Redis, Valkey, Mongo (four data stores)
- **api.yml** — REST API (Postgres + Redis)
- **frontend.yml** — nginx serving a static HTML page

All services share the same Compose project network so the API can reach `postgres`, `redis`, etc.

## Run with dcctl

From the **dcctl repo root**:

```bash
dcctl --config-file=examples/04-full-stack/dcctl.yml up
```

Then:

- **Frontend:** http://localhost:8080
- **API:** http://localhost:4000/health — http://localhost:4000/api/items
- **Postgres:** localhost:5432, **Redis:** 6379, **Valkey:** 6380, **Mongo:** 27017

## Useful commands

```bash
dcctl --config-file=examples/04-full-stack/dcctl.yml status
dcctl --config-file=examples/04-full-stack/dcctl.yml show-ports
dcctl --config-file=examples/04-full-stack/dcctl.yml logs -f api
dcctl --config-file=examples/04-full-stack/dcctl.yml down
```

## Files

- `dcctl.yml` — three services in one environment: `databases`, `api`, `frontend`
- `default/databases.yml` — postgres, redis, valkey, mongo
- `default/api.yml` — API (build from `api/`)
- `default/frontend.yml` — nginx + static HTML from `frontend/html/`
- `api/` — Node app (Dockerfile, server.js)
- `frontend/html/index.html` — static page
