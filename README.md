# Restic Controller

A controller for the restic backup program that manages restic repositories, configurable via a yaml file.

## Features

- Backupping
- Integrity Checks
- Retention
- Auto-reloading of the configuration at runtime

## Usage Example

### docker-compose.yml
```yaml
services:
  restic-controller:
    image: ghcr.io/syynek/restic-controller:latest
    volumes:
      - "./config.yml:/app/config.yml"
      - "./backup-data:/data:ro"
      - "./backups:/repositories"
```

> [!NOTE]  
> To use the restic controller with docker, you will need to bind mount any files you would like to back up into a folder within the docker container.

### config.yml

```yaml
log:
  level: info

repositories:
  - name: local
    url: /repositories/repository
    password: test
    auto_initialize: true
    backup:
      schedule: "0 * * * *"
      run_on_startup: true
      include_files:
        - /data
      exclude_files:
        - /data/example.txt
    integrity_check:
      schedule: "15 * * * *"
    retention:
      schedule: "30 * * * *"
      policy:
        keep_last: 1
```
