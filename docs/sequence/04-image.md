# Sequence: Image (ls, pull, remove)

From `dcctl image ls [service]`, `dcctl image pull [service]`, or `dcctl image remove [service]`: load config for the current environment, resolve manifest paths (all or for the given service), collect image names from those manifests, then list them, pull them, or remove them (with confirmation).

Same manifest resolution as other commands (see [03-manifest-resolution](./03-manifest-resolution.md)); images are taken from the `image:` field of each service in the resolved compose files.

```mermaid
sequenceDiagram
    participant User
    participant dcctl
    participant FS
    participant Docker

    User->>dcctl: dcctl image <ls|pull|remove> [service]
    dcctl->>dcctl: load config (--config-file or ~/.dcctl/dcctl.yml)
    dcctl->>dcctl: resolve environment (-e or default_environment)
    dcctl->>dcctl: resolve manifest paths (all or for [service])
    dcctl->>FS: read each manifest YAML
    dcctl->>dcctl: collect unique image names from services.image
    alt no images
        dcctl->>User: warn, exit 0
    end
    alt ls
        dcctl->>Docker: docker image inspect (optional, for local ID)
        dcctl->>User: print image list (name + local ID or "not found")
    else pull
        loop each image
            dcctl->>Docker: docker image pull <image>
            Docker-->>dcctl: success / warn
        end
        dcctl->>User: Done
    else remove
        dcctl->>User: print image list, ask confirmation (yes)
        alt user confirms
            loop each image
                dcctl->>Docker: docker image rm <image>
            end
            dcctl->>User: Done
        else cancelled
            dcctl->>User: cancelled
        end
    end
```
