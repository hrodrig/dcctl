# Sequence: Show ports

From `dcctl show-ports`: load config, resolve manifest paths for the current environment, extract published ports from each compose file, and print them per service. Lets the developer see at a glance where to connect (browser, API, DB client).

Uses the same config and manifest resolution as [03-manifest-resolution](./03-manifest-resolution.md); then parses each manifest for `services.*.ports` and normalizes host:container or host-only mappings.

```mermaid
sequenceDiagram
    participant User
    participant dcctl
    participant FS

    User->>dcctl: dcctl show-ports
    dcctl->>dcctl: load config, resolve environment
    dcctl->>dcctl: resolve manifest paths (env services)
    dcctl->>FS: read each manifest YAML
    dcctl->>dcctl: parse ports per service (host:container or host)
    dcctl->>dcctl: sort by service name
    dcctl->>User: print service → port(s) (e.g. web → 3000, redis → 6379)
```
