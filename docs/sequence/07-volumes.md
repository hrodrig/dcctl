# Sequence: Volumes clean

From `dcctl volumes clean [volume...]`: optionally load config (for project name), then list Docker volumes (filtered by project label or all), let the user select which to remove (by name, by number, or via `--match` regex), ask for confirmation, and run `docker volume rm` on the selected volumes.

- **With volume names as args:** no need to list; remove those volumes (after checking they exist).
- **Without args:** list volumes (by default for the compose project from config), then either interactive selection or `--match` regex; confirm and remove.

```mermaid
sequenceDiagram
    participant User
    participant dcctl
    participant Docker

    User->>dcctl: dcctl volumes clean [volume...] [--match regex] [--all]
    dcctl->>dcctl: ensure Docker available
    alt volume names given
        dcctl->>dcctl: unique volume names from args
        loop each volume
            dcctl->>Docker: docker volume inspect (exists?)
            alt not found
                dcctl->>User: warn, skip
            else exists
                dcctl->>Docker: docker volume rm <name>
                Docker-->>dcctl: success / error
            end
        end
        dcctl->>User: Done
    else no args
        dcctl->>dcctl: load config (for project name)
        alt --all
            dcctl->>Docker: docker volume ls (no filter)
        else default
            dcctl->>Docker: docker volume ls --filter label=com.docker.compose.project=<name>
            alt no volumes for project
                dcctl->>Docker: docker volume ls (all), warn
            end
        end
        Docker-->>dcctl: volume names
        alt --match set
            dcctl->>dcctl: filter names by regex
            dcctl->>User: confirm removal (yes)
            alt confirmed
                dcctl->>Docker: docker volume rm <each>
            else cancelled
                dcctl->>User: cancelled
            end
        else interactive
            dcctl->>User: print "Available volumes:" + list
            User->>dcctl: selection (numbers, names, ranges)
            dcctl->>dcctl: parse selection
            dcctl->>User: confirm removal (yes)
            alt confirmed
                dcctl->>Docker: docker volume rm <selected>
            else cancelled
                dcctl->>User: cancelled
            end
        end
    end
```
