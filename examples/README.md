# dcctl examples

Example stacks you can run with [dcctl](https://github.com/hrodrig/dcctl) to see how it orchestrates Docker Compose per environment. From a simple counter to a multi-manifest full-stack.

**Config location:** dcctl only reads config from `~/.dcctl/dcctl.yml` by default; it does not look in the current directory. So for these examples you **must** pass `--config-file` to the config in the example (or copy it to `~/.dcctl/dcctl.yml`). All commands below assume you run **from the dcctl repository root** with `dcctl --config-file=examples/.../dcctl.yml`, or from the example directory with `dcctl --config-file=dcctl.yml`.

**Volumes and local data:** All examples use **named Docker volumes** (e.g. `pg_data`, `portainer_data`) for persistent data, so nothing is written into the repo. The only exception is **05-multi-env** (Minecraft), which bind-mounts `./minecraft_data` for the world; that directory is listed in `05-multi-env/minecraft/.gitignore` so it is not committed.

| Example | Description | Stack | Ports |
|---------|-------------|--------|-------|
| [01-counter-redis](./01-counter-redis/) | Simple: web counter stored in Redis | Node app + Redis | 3000 (web), 6379 (redis) |
| [02-wordpress-mariadb](./02-wordpress-mariadb/) | Classic CMS | WordPress + MariaDB | 8080 (wordpress) |
| [03-api-postgres-redis](./03-api-postgres-redis/) | Polyglot persistence | API + Postgres + Redis | 4000 (api), 5432, 6379 |
| [04-full-stack](./04-full-stack/) | Multiple manifests (enredado) | Postgres, Redis, Valkey, Mongo + API + nginx frontend | 8080 (frontend), 4000 (api), 5432, 6379, 6380, 27017 |
| [05-multi-env](./05-multi-env/) | Two environments, default_environment ≠ default | default: counter+Redis; minecraft: Minecraft server (awesome-compose) | default: 3000, 6379; minecraft: 25565 |
| [06-postgresql-pgadmin](./06-postgresql-pgadmin/) | Postgres + pgAdmin (env vars, depends_on) | PostgreSQL + pgAdmin (awesome-compose) | 5432 (postgres), 5050 (pgAdmin) |
| [07-flask-redis](./07-flask-redis/) | Flask + Redis (env, depends_on) | Python Flask + Redis (awesome-compose) | 8000 (web), 6379 (redis) |
| [08-nginx-golang](./08-nginx-golang/) | Nginx + Go backend (depends_on) | Nginx proxy + Chi backend (awesome-compose) | 80 (nginx) |
| [09-portainer](./09-portainer/) | Portainer CE (env) | Docker UI (awesome-compose) | 9000 (portainer) |
| [10-prometheus-grafana](./10-prometheus-grafana/) | Prometheus + Grafana (env, depends_on) | Monitoring stack (awesome-compose) | 9090 (prometheus), 3000 (grafana) |

## Quick run

```bash
# From repo root, with dcctl on PATH (e.g. make build && export PATH="$PWD:$PATH")
dcctl --config-file=examples/01-counter-redis/dcctl.yml up
# → http://localhost:3000

dcctl --config-file=examples/02-wordpress-mariadb/dcctl.yml up
# → http://localhost:8080

dcctl --config-file=examples/03-api-postgres-redis/dcctl.yml up
# → http://localhost:4000/health

dcctl --config-file=examples/04-full-stack/dcctl.yml up
# → http://localhost:8080 (frontend), http://localhost:4000 (API)

dcctl --config-file=examples/05-multi-env/dcctl.yml up
# → Minecraft on 25565 (default_environment is minecraft)
dcctl --config-file=examples/05-multi-env/dcctl.yml -e default up
# → http://localhost:3000 (counter+Redis)

dcctl --config-file=examples/06-postgresql-pgadmin/dcctl.yml up
# → Postgres :5432, pgAdmin http://localhost:5050
dcctl --config-file=examples/07-flask-redis/dcctl.yml up
# → http://localhost:8000
dcctl --config-file=examples/08-nginx-golang/dcctl.yml up
# → http://localhost:80 (Nginx → Go)
dcctl --config-file=examples/09-portainer/dcctl.yml up
# → http://localhost:9000
dcctl --config-file=examples/10-prometheus-grafana/dcctl.yml up
# → Prometheus :9090, Grafana http://localhost:3000
```

**See published ports for any example:** after `up`, run `dcctl --config-file=examples/<name>/dcctl.yml show-ports` to list ports per service (e.g. from repo root: `dcctl --config-file=examples/01-counter-redis/dcctl.yml show-ports`).

## What each example shows

- **01** — One manifest, one environment; good first step.
- **02** — Off-the-shelf images (WordPress, MariaDB), volumes.
- **03** — Single stack with two data stores (Postgres + Redis).
- **04** — Several compose files in one environment; dcctl passes all of them to `docker compose -f ... -f ... -f ...`.
- **05** — Multiple environments: `default` (counter+Redis) and `minecraft` (Minecraft from [awesome-compose](https://github.com/docker/awesome-compose/tree/master/minecraft)); `default_environment: minecraft` so `dcctl up` runs Minecraft unless you use `-e default`.
- **06–10** — All use **environment variables** and **service dependencies** (`depends_on`, healthchecks). Sources: [docker/awesome-compose](https://github.com/docker/awesome-compose). **06** Postgres + pgAdmin (credentials via env). **07** Flask + Redis (`REDIS_URL`). **08** Nginx + Go (proxy depends on backend). **09** Portainer (single service, port/env). **10** Prometheus + Grafana (admin env, Grafana depends on Prometheus).

See each example’s `README.md` for details and `dcctl down` when done.
