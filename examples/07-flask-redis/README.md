# Example 07: Flask + Redis

Python Flask app with Redis for a hit counter. Shows **environment variables** (`REDIS_URL`, `FLASK_PORT`) and **depends_on** so the web app starts after Redis is healthy.

**Source:** [docker/awesome-compose — flask-redis](https://github.com/docker/awesome-compose/tree/master/flask-redis).

## Environment variables and dependencies

| Service | Env vars | Depends on |
|---------|----------|------------|
| redis   | —        | —          |
| web     | `REDIS_URL`, `FLASK_PORT` | redis (healthy) |

## Run with dcctl

```bash
dcctl --config-file=dcctl.yml up
# → http://localhost:8000 (refresh to see counter)
```

## Useful commands

```bash
dcctl --config-file=dcctl.yml status
dcctl --config-file=dcctl.yml logs -f web
dcctl --config-file=dcctl.yml down
```

## Files

- `dcctl.yml` — dcctl config
- `default/stack.yml` — Compose: redis + web (env, healthcheck, depends_on)
- `Dockerfile`, `app.py`, `requirements.txt` — Flask app
- `.env.example` — optional env template
