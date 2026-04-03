# Unit Tests

Fast tests that run with:

```bash
go test ./...
```

They **exclude** build-tagged packages:

- `integration` — see [../integration/README.md](../integration/README.md)
- `e2e` — see [../e2e/README.md](../e2e/README.md)

## Coverage areas

- Config load/validate (`internal/config`)
- Engine (Simple/Pro planning, drift, upgrade)
- LLM clients and prompt building
- Security gates, auto-correction, policy helpers
- Cost estimation (`shopspring/decimal`)
- Terraform renderer/runner/workspace helpers (mocked `exec` where used)
- Providers and bootstrap rendering
- Cobra command wiring (`cmd`)

Use table-driven tests and keep cases close to the code they protect.
