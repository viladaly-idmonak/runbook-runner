# Example Runbook

This runbook demonstrates the supported markdown format for runbook-runner.

## Check disk space

Ensure there is enough free disk space before proceeding.

```sh
df -h /
```

<!-- rollback -->
```sh
echo "No rollback needed for disk check"
```

## Create backup directory

Create a timestamped backup directory.

```sh
mkdir -p /tmp/backup/$(date +%Y%m%d%H%M%S)
```

<!-- rollback -->
```sh
rm -rf /tmp/backup
```

## Sync application files

Copy application files to the backup directory.

```sh
rsync -av /opt/myapp/ /tmp/backup/
```
