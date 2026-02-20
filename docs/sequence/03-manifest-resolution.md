# Sequence: Manifest resolution

How dcctl resolves the config file path, the selected environment, and the list of manifest file paths (per-environment services or explicit targets). **Config path:** dcctl uses only `--config-file` if set, otherwise **only** `~/.dcctl/dcctl.yml` (it does not look in the current directory).

```mermaid
sequenceDiagram
    participant User
    participant dcctl
    participant FS

    User->>dcctl: dcctl [flags] <cmd> [args]
    dcctl->>dcctl: parse flags (--config-file, -e, --debug)
    alt --config-file set
        dcctl->>FS: file exists?
        FS-->>dcctl: yes → use path
    else default
        dcctl->>FS: ~/.dcctl/dcctl.yml exists?
        alt not found
            dcctl->>User: error: no config found (hint: dcctl config > ...)
        end
    end
    dcctl->>dcctl: load YAML, validate (environments, services)
    dcctl->>dcctl: environment = -e or dcctl.default_environment or "default"
    dcctl->>dcctl: env dir = dir(config_path) / environment
    alt args = service or manifest names
        loop for each target
            dcctl->>dcctl: manifest path = env_dir / target.yml or find by service name
            dcctl->>FS: file exists?
        end
    else no args
        dcctl->>dcctl: for each env.Services: env_dir / service.yml
        dcctl->>FS: collect existing paths, warn missing
    end
    dcctl->>dcctl: return list of manifest paths
```
