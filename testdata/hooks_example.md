# Hooks Example Runbook

This runbook demonstrates lifecycle hook usage with pre-step, post-step,
and on-error hooks.

## Step 1: Prepare environment

```sh
echo "preparing environment"
mkdir -p /tmp/rr-work
```

**Rollback:**

```sh
rm -rf /tmp/rr-work
```

## Step 2: Run migration

```sh
echo "running migration"
```

**Rollback:**

```sh
echo "reversing migration"
```

## Step 3: Verify deployment

```sh
echo "verifying deployment"
```
