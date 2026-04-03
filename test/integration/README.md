# Integration Tests

Integration tests run against LocalStack (AWS emulator) and validate the full
pipeline from template selection through Terraform plan.

## Prerequisites

- Docker (for LocalStack)
- Terraform CLI
- Go 1.26.1+

## Running locally

```bash
# Start LocalStack + run tests + stop LocalStack
make integration

# Or manually:
docker compose up -d --wait
go test -tags integration -timeout 300s -v ./test/integration/...
docker compose down
```

## What's tested

- **Simple Mode pipeline**: template selection, security gate, render, tfvars, terraform init+plan
- **Pro Mode pipeline**: VPC template render, terraform init+plan, cost estimation
- **Security gate**: Checkov scan blocks bad templates, passes good ones
- **Cost estimation**: all simple templates return valid cost breakdowns

## CI

The `integration` job in `.github/workflows/ci.yml` runs these tests on every push and pull request using a **LocalStack** service container (pinned image), **Terraform** setup, and `go test -tags integration ./test/integration/...`.

For contributor commands and lint/test expectations, see [CONTRIBUTING.md](../../CONTRIBUTING.md).
