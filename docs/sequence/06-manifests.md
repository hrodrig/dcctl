# Sequence: Manifests

From `dcctl manifests`: load config, resolve the environment (from `-e` or `default_environment`), ensure the environment directory exists, and write an embedded sample compose file (`sample-app.yml`) into that directory. Used for onboarding when the user has no manifests yet.

Config and environment resolution follow [03-manifest-resolution](./03-manifest-resolution.md). The written file is a minimal compose template so the user can rename it (e.g. to `app.yml`) and add the service name to `dcctl.environments.<env>.services`.

```mermaid
sequenceDiagram
    participant User
    participant dcctl
    participant Config
    participant FS

    User->>dcctl: dcctl manifests [-e env]
    dcctl->>Config: load config (--config-file or ~/.dcctl/dcctl.yml)
    dcctl->>dcctl: resolve environment (-e or default_environment)
    dcctl->>FS: MkdirAll(config_dir / environment)
    alt mkdir fails
        dcctl->>User: error, exit
    end
    dcctl->>dcctl: target = config_dir / environment / "sample-app.yml"
    dcctl->>FS: write embedded sample-app.yml to target
    alt write fails
        dcctl->>User: error, exit
    end
    dcctl->>User: "Sample manifest written to <path>"
```
