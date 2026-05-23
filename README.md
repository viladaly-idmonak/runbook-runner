# runbook-runner

Execute structured markdown runbooks as automated shell scripts with step validation and rollback support.

## Installation

```bash
go install github.com/runbook-runner/runbook-runner@latest
```

Or build from source:

```bash
git clone https://github.com/runbook-runner/runbook-runner.git && cd runbook-runner && go build -o runbook-runner .
```

## Usage

Write a markdown runbook with fenced shell code blocks, then execute it:

```bash
runbook-runner run deploy.md
```

**Example runbook (`deploy.md`):**

```markdown
## Step 1: Pull latest image
```sh
docker pull myapp:latest
```

## Step 2: Restart service
```sh
systemctl restart myapp
```
```

Each step is validated before proceeding. If a step fails, `runbook-runner` automatically triggers any defined rollback blocks.

**Common flags:**

```bash
runbook-runner run deploy.md --dry-run        # Preview steps without executing
runbook-runner run deploy.md --step 2         # Start from a specific step
runbook-runner run deploy.md --rollback-on-fail  # Auto-rollback on any failure
```

## How It Works

1. Parses markdown headings as named steps
2. Executes fenced `sh` or `bash` code blocks in order
3. Validates exit codes between steps
4. Runs rollback blocks in reverse order on failure

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

MIT © runbook-runner contributors