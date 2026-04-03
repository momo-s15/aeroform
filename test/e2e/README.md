# End-to-End Tests

E2E tests deploy real infrastructure to a cloud account, verify it works, then
destroy everything. They are guarded by `//go:build e2e` and never run in
normal CI.

## Prerequisites

- Real AWS credentials (via env vars, profile, or OIDC role)
- Terraform CLI
- Go 1.24+

## Running locally

```bash
# Requires real AWS credentials in your environment
export AWS_REGION=us-east-1
make e2e
```

## What's tested

- **Simple Mode static-site**: full deploy → verify outputs → destroy
- **Simple Mode lambda-api**: full deploy → verify outputs → destroy
- **Pro Mode VPC**: full deploy → verify outputs → destroy

Every test uses `defer terraform.Destroy(workDir, true)` to ensure cleanup
even on failure.

## CI

The `.github/workflows/e2e.yml` workflow runs E2E tests only on version tag
pushes (`v*`). It uses OIDC to assume an IAM role from the `e2e` GitHub
environment, which should be configured with:

- `AWS_E2E_ROLE_ARN` — IAM role ARN with deploy permissions
- `AWS_ACCOUNT_ID` — AWS account ID for template variables

This ensures real cloud tests run before release artifacts are published.
