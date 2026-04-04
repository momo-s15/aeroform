# End-to-End Tests

E2E tests deploy real infrastructure in **AWS**, verify behavior, then destroy resources. They are guarded by `//go:build e2e` and are **not** part of `go test ./...`.

## Prerequisites

- **AWS** credentials (environment variables, shared config, or instance role)
- **Terraform** on `PATH`
- **Go 1.26.1+**

## Running locally

```bash
export AWS_REGION=us-east-1
# Ensure AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY (or equivalent) are set
make e2e
```

Or:

```bash
go test -tags e2e -timeout 900s -v ./test/e2e/...
```

**Warning:** E2E tests create billable resources. Use a dedicated sandbox account when possible.

## What is tested

| Test | Area |
|------|------|
| `TestSimpleModeDeployStaticSite` | Simple Mode static-site on AWS |
| `TestSimpleModeDeployLambdaAPI` | Simple Mode `lambda-api` on AWS |
| `TestProModeDeployVPC` | Pro Mode VPC module composition on AWS |

Tests use `defer` cleanup so `terraform destroy` still runs after failures.

## CI (GitHub Actions)

Workflow: `.github/workflows/e2e.yml` — triggers on **tags** matching `v*` (e.g. `v1.0.5`).

1. **Detect E2E configuration** — If `AWS_E2E_ROLE_ARN` and `AWS_ACCOUNT_ID` are **empty** (repository or `e2e` environment secrets), subsequent steps are **skipped** and the workflow **succeeds** with a notice. This is intentional for forks and repos that have not wired OIDC yet.
2. When both secrets are set, the job uses **OIDC** (`aws-actions/configure-aws-credentials`) to assume the role and runs `go test -tags e2e ./test/e2e/...`.

Configure the **`e2e`** GitHub Environment (or repository secrets) with:

- `AWS_E2E_ROLE_ARN` — IAM role ARN trusted for GitHub OIDC
- `AWS_ACCOUNT_ID` — used where tests need an account id variable

See also [CONTRIBUTING.md](../../CONTRIBUTING.md) for the full CI matrix.
