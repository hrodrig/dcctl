# Example 05: Multiple environments (default + minecraft)

Same idea as example 01 (counter + Redis) but with **two environments**: `default` and `minecraft`. The config sets **default_environment: minecraft**, so `dcctl up` starts the Minecraft server unless you pass `-e default`.

Inspired by [docker/awesome-compose — Minecraft server](https://github.com/docker/awesome-compose/tree/master/minecraft).

## Environments

| Environment | Services | Description |
|-------------|----------|-------------|
| **default** | `stack` | Counter app + Redis (like example 01) — http://localhost:3000 |
| **minecraft** | `server` | Minecraft server (itzg/minecraft-server) — port 25565 |

## Run with dcctl

From the **dcctl repo root** (dcctl on PATH):

```bash
# Start Minecraft (default_environment is minecraft)
dcctl --config-file=examples/05-multi-env/dcctl.yml up
# → Minecraft on port 25565; world data in examples/05-multi-env/minecraft/minecraft_data/

# Start counter + Redis instead
dcctl --config-file=examples/05-multi-env/dcctl.yml -e default up
# → http://localhost:3000
```

From this example directory:

```bash
dcctl --config-file=dcctl.yml up          # Minecraft
dcctl --config-file=dcctl.yml -e default up   # Counter + Redis
```

## Useful commands

```bash
dcctl --config-file=examples/05-multi-env/dcctl.yml status
dcctl --config-file=examples/05-multi-env/dcctl.yml -e default status
dcctl --config-file=examples/05-multi-env/dcctl.yml image ls
dcctl --config-file=examples/05-multi-env/dcctl.yml down
dcctl --config-file=examples/05-multi-env/dcctl.yml -e default down
```

## Files

- `dcctl.yml` — two environments (`default`, `minecraft`), **default_environment: minecraft**
- `default/stack.yml` — Compose: counter app + Redis
- `minecraft/server.yml` — Compose: Minecraft server (from awesome-compose)
- `Dockerfile`, `package.json`, `server.js` — counter app for the `default` environment
