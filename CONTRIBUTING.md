# Contributing

Thanks for contributing to Aeroform.

Keep pull requests **focused** (one feature or fix), **tested**, and consistent with existing style. For historical implementation context, see `AEROFORM_PLAN.md` and `aeroform_blueprint.md`.

End-user install options are in [docs/installation.md](docs/installation.md) (curl `install.sh`, Releases, [Homebrew tap](docs/homebrew-tap.md), `go install`).

## What you need

- **Go 1.26.1+** (match `go.mod` / `toolchain`; CI uses `go-version-file: go.mod`)
- **Terraform** on `PATH` for integration-style checks locally
- **Docker** (optional) for `make integration` with LocalStack
- **golangci-lint v2** (CI uses **v2.9**) for `make lint` — [install](https://golangci-lint.run/welcome/install/) or rely on CI

## Local checks

```bash
go test ./...                    # unit tests (default; no integration/e2e tags)
golangci-lint run ./...          # same linters as CI (.golangci.yml v2)
make integration                 # LocalStack + integration tests (needs Docker)
```

Integration tests only:

```bash
go test -tags integration -timeout 300s ./test/integration/...
```

E2E tests (real AWS; **charges possible**):

```bash
export AWS_REGION=us-east-1
# ... credentials or instance profile ...
go test -tags e2e -timeout 900s ./test/e2e/...
```

## CI (GitHub Actions)

| Workflow | When | What |
|----------|------|------|
| `ci.yml` | push / PR | lint (golangci-lint **v2.9**), unit tests, LocalStack integration |
| `security.yml` | push / PR | `govulncheck`, dependency review (PRs) |
| `e2e.yml` | tags `v*` | E2E if `AWS_E2E_ROLE_ARN` + `AWS_ACCOUNT_ID` set for `e2e` env; else skip |
| `release.yml` | tags `v*.*.*` | tests + build release assets |

## Good first contributions

- Template or provider tweaks in `internal/providers/` and `templates/`
- Command help text and UX in `cmd/`
- Security scan summaries in `internal/security/`
- Docs in `README.md`, `docs/`, and this file

## Coding expectations

- Match surrounding naming, structure, and import style.
- No unrelated refactors in feature PRs.
- Prefer clear, explicit logic over clever shortcuts.
- User-facing errors and logs should be actionable.
- Add tests for new package-level behavior.

## Pull request checklist

- [ ] Scope and intent are clear in the description
- [ ] Tests added or updated where behavior changed
- [ ] `go test ./...` passes
- [ ] `golangci-lint run ./...` passes (or CI green)
- [ ] User-facing docs updated if commands or flags changed

## Reporting problems

Open an issue with:

- The exact command you ran
- Relevant prompt, `config.yaml`, or flags
- Expected vs actual behavior
- Logs or terminal snippets (redact secrets)
