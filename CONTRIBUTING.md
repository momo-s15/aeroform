# Contributing

Thanks for contributing to Aeroform.

The project is built in phases. Contributions are easiest to review when they are small, tested, and aligned with the active phase tracked in AEROFORM_PLAN.md.

## Development Workflow

1. Open an issue or pick an existing one.
2. Keep each pull request focused on one feature or fix.
3. Add or update tests for behavior changes.
4. Run `go test ./...` before opening a pull request.

## What To Work On

Good first contributions:

- Expand provider template mappings in `internal/providers/`
- Improve command UX wording in `cmd/`
- Improve security summary formatting in `internal/security/`
- Improve docs and examples in `README.md` and this file

Higher-scope contributions:

- Add deeper Terraform execution wiring in Pro Mode
- Add richer drift detection and live-state reconciliation
- Add docs and release workflow polish for launch readiness

## Coding Expectations

- Preserve existing project structure and naming patterns.
- Avoid unrelated refactors in feature PRs.
- Prefer explicit, readable logic over clever shortcuts.
- Keep user-facing output clear and actionable.
- Add tests for new package-level behavior.

## Pull Request Checklist

- [ ] Feature or fix is scoped and described clearly
- [ ] Tests added or updated
- [ ] `go test ./...` passes locally
- [ ] Docs updated when behavior changes

## Reporting Problems

If you find unexpected behavior, open an issue with:

- The command you ran
- Relevant input (prompt/config)
- Expected vs actual behavior
- Logs or terminal output snippets
