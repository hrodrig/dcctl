# Example 06: PostgreSQL + pgAdmin

PostgreSQL database with pgAdmin web UI. Shows **environment variables** for credentials and **depends_on** so pgAdmin starts only after Postgres is healthy.

**Source:** [docker/awesome-compose — postgresql-pgadmin](https://github.com/docker/awesome-compose/tree/master/postgresql-pgadmin).

## Environment variables and dependencies

| Service   | Env vars | Depends on   |
|-----------|-----------|--------------|
| postgres  | `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` | — |
| pgadmin   | `PGADMIN_EMAIL`, `PGADMIN_PASSWORD` | postgres (healthy) |

Defaults are in the compose file (`:-app`, `:-secret`, etc.). Override via `.env` or shell.

## Run with dcctl

```bash
dcctl --config-file=dcctl.yml up
# → Postgres :5432, pgAdmin http://localhost:5050 (login with PGADMIN_EMAIL / PGADMIN_PASSWORD)
```

In pgAdmin, add server: host `postgres`, user/password/database from your env (defaults: app / secret / appdb).

## Useful commands

```bash
dcctl --config-file=dcctl.yml status
dcctl --config-file=dcctl.yml logs -f postgres
dcctl --config-file=dcctl.yml down
```

## Files

- `dcctl.yml` — dcctl config
- `default/stack.yml` — Compose: postgres + pgadmin (env vars, healthcheck, depends_on)
- `.env.example` — optional env template
