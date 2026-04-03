# Aeroform

From a student's first cloud project to a team's production infrastructure.

Aeroform is an open-source CLI that turns plain English into infrastructure workflows.

Repository: https://github.com/momo-s15/aeroform

- Simple Mode is a guided, beginner-first path that favors safe defaults and cost awareness.
- Pro Mode is a config-driven workflow for platform and DevOps teams.

## Current Status

The project is in active build mode and ships in phases.

- Simple Mode launch and guardrails are implemented.
- Pro Mode command foundation and multi-cloud planning are implemented.
- Security scan and drift command scaffolds are implemented.

## Why Aeroform

- Plain-English workflows instead of hand-authoring every initial Terraform file.
- Mode-aware UX so beginner and professional users can share one tool.
- Local AI support through Ollama with fallback behavior when local AI is unavailable.
- Security gate integration for both Simple and Pro paths.

## Quick Start (Simple Mode)

Prerequisites:

- Go 1.24+
- Optional: Ollama for AI-backed template selection

Build and run:

```bash
go test ./...
go run . setup
go run . launch
```

Common Simple Mode commands:

- `aeroform launch`
- `aeroform status`
- `aeroform cost`
- `aeroform domain add mysite.com`
- `aeroform logs`
- `aeroform destroy --confirm`

## Quick Start (Pro Mode)

1. Prepare a `config.yaml` from `config.yaml.example`.
2. Run planning and generation commands.

```bash
go run . plan "resilient kubernetes platform with private database"
go run . generate "resilient kubernetes platform with private database"
go run . scan
go run . drift "resilient kubernetes platform with private database"
```

Common Pro Mode commands:

- `aeroform bootstrap`
- `aeroform plan "..."`
- `aeroform generate "..."`
- `aeroform scan`
- `aeroform drift "..."`

## Upgrade Path

When a Simple Mode project grows, generate a Pro config:

```bash
go run . upgrade
go run . upgrade --confirm
```

The second command writes `config.yaml` (or a custom output path) and enables Pro Mode workflows.

## Project Layout

Key directories:

- `cmd/` command entrypoints
- `internal/config/` config loading and validation
- `internal/engine/` planning and mode logic
- `internal/providers/` cloud-specific behavior
- `internal/security/` security gates and scanner integration
- `internal/prompt/` constrained prompt builders
- `internal/llm/` local AI client interfaces

## Multi-Cloud Support

Provider abstractions exist for:

- AWS
- Azure
- GCP

Each provider exposes region handling, Pro template support, and bootstrap guidance.

## Security Model

- Simple Mode blocks known unsafe behavior and applies auto-corrections for low-risk issues.
- Pro Mode integrates Checkov and tfsec when available and enforces scan gating on failure.
- Security reporting is plain language in Simple Mode and explicit scan status in Pro Mode.

## Contributing

See `CONTRIBUTING.md` for contribution workflow, expectations, and first-issue guidance.

## Responsible Disclosure

See `SECURITY.md` for how to report vulnerabilities responsibly.

## First GitHub Prototype Release

Suggested first release checklist:

- Run `go test ./...`
- Confirm `go run . launch` and `go run . plan "..."` both run locally
- Push to `main`
- Create tag `v0.1.0-prototype`
- Publish a GitHub Release with short demo notes and known limitations
