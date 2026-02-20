# Example 03: API + Postgres + Redis

Polyglot stack: a small REST API that uses **PostgreSQL** for persistence and **Redis** for a counter cache. Demonstrates multiple data stores in one stack.

## Stack

- **PostgreSQL** — main store (table `items`)
- **Redis** — cache/counter (e.g. `items:count`)
- **API** — Node.js Express: `GET/POST /items`, `GET /stats`, `GET /health`

## Run with dcctl

From the **dcctl repo root**:

```bash
dcctl --config-file=examples/03-api-postgres-redis/dcctl.yml up
```

Then:

- http://localhost:4000/health
- http://localhost:4000/items — list items (empty at first)
- `curl -X POST http://localhost:4000/items -H "Content-Type: application/json" -d '{"name":"test"}'`
- http://localhost:4000/stats — Redis vs Postgres counts

## Useful commands

```bash
dcctl --config-file=examples/03-api-postgres-redis/dcctl.yml status
dcctl --config-file=examples/03-api-postgres-redis/dcctl.yml logs -f api
dcctl --config-file=examples/03-api-postgres-redis/dcctl.yml down
```

## Files

- `dcctl.yml` — dcctl config
- `default/stack.yml` — Compose: postgres, redis, api
- `Dockerfile`, `package.json`, `server.js` — API app
