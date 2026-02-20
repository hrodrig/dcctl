# Example 10: Prometheus + Grafana

Prometheus metrics and Grafana dashboards. Shows **environment variables** for ports and admin credentials, and **depends_on** so Grafana starts after Prometheus. Grafana is provisioned with a Prometheus datasource.

**Source:** [docker/awesome-compose — prometheus-grafana](https://github.com/docker/awesome-compose/tree/master/prometheus-grafana).

## Environment variables and dependencies

| Service    | Env vars | Depends on |
|------------|----------|------------|
| prometheus | `PROMETHEUS_PORT` | — |
| grafana    | `GF_ADMIN_USER`, `GF_ADMIN_PASSWORD`, `GRAFANA_PORT` | prometheus |

## Run with dcctl

```bash
dcctl --config-file=dcctl.yml up
# → Prometheus http://localhost:9090, Grafana http://localhost:3000 (login with GF_ADMIN_*)
```

Grafana is pre-provisioned with a Prometheus datasource (Configuration → Data sources).

## Useful commands

```bash
dcctl --config-file=dcctl.yml status
dcctl --config-file=dcctl.yml logs -f grafana
dcctl --config-file=dcctl.yml down
```

## Files

- `dcctl.yml` — dcctl config
- `default/stack.yml` — Compose: prometheus + grafana (env, depends_on)
- `default/prometheus/prometheus.yml` — Prometheus config
- `default/grafana/provisioning/datasources/datasource.yml` — Grafana datasource
- `.env.example` — optional
