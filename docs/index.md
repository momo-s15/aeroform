# Aeroform Documentation

Aeroform is an open-source CLI tool that deploys cloud infrastructure using natural language. It supports AWS, Azure, and GCP through pre-validated, security-hardened Terraform templates.

**Install Aeroform:** use the **[installation guide](installation.md)** (curl install script for macOS/Linux, optional [Homebrew tap](homebrew-tap.md), `go install` for Go users, Windows `.exe` from Releases).

**Runtime requirements:** [Terraform](https://developer.hashicorp.com/terraform/install) and a supported LLM backend ([Ollama](https://ollama.com/download) by default). Building from source needs [Go 1.26.1+](https://go.dev/dl/) (matches `go.mod` and CI). Pro Mode security gating expects [Checkov](https://www.checkov.io/) (and optionally [tfsec](https://github.com/aquasecurity/tfsec)) on your `PATH` when you run scans locally.

**Two modes, one tool:**

- **Simple Mode** — Describe what you want, get it deployed. Guided prompts, cost estimates, auto-security. Perfect for students and solo developers.
- **Pro Mode** — Config-driven multi-template composition with Checkov/tfsec gating, workspace environments, and drift detection. Built for teams.

## Quick Start

```bash
# Install (macOS/Linux: curl -fsSL https://raw.githubusercontent.com/momo-s15/aeroform/main/install.sh | sh)
go install github.com/momo-s15/aeroform@latest

# Simple Mode — deploy a website in one command
aeroform setup
aeroform launch

# Pro Mode — config-driven deployment
aeroform init           # scaffold config.yaml
aeroform plan "kubernetes cluster with database"
aeroform generate "kubernetes cluster with database"
```

## Guides

| Guide | Audience | What you'll learn |
|-------|----------|-------------------|
| [Installation](installation.md) | Everyone | Install script, Homebrew tap, go install, Windows |
| [Getting Started — Simple](getting-started-simple.md) | Students, solo devs | Install to live website in 5 minutes |
| [Getting Started — Pro](getting-started-pro.md) | Teams, DevOps | Config to EKS cluster with security gating |

## References

| Page | What it covers |
|------|----------------|
| [Homebrew tap](homebrew-tap.md) | Tap setup; CI can auto-push the formula when `HOMEBREW_TAP_TOKEN` is set |
| [Template Reference](template-reference.md) | Every template (Simple + Pro, all clouds) with resources and costs |
| [LLM Backends](llm-backends.md) | Setup for Ollama, OpenAI, and AWS Bedrock |
| [Security Model](security-model.md) | Scanning, auto-correction, custom policies, drift detection |
| [Architecture Decisions](architecture-decisions.md) | Why two modes, constrained generation, Go, Terraform, Checkov |

## Contributing & CI

Contributors should follow [CONTRIBUTING.md](../CONTRIBUTING.md). The repository runs **lint** (golangci-lint **v2.9** in CI via `.github/workflows/ci.yml`), **unit tests**, **integration tests** (LocalStack in GitHub Actions), and **govulncheck** on pushes and pull requests. **E2E tests** (real AWS) run on version tags when `AWS_E2E_ROLE_ARN` and `AWS_ACCOUNT_ID` are configured for the `e2e` environment; otherwise that job skips cleanly.

## Key Commands

| Command | Mode | Description |
|---------|------|-------------|
| `aeroform setup` | Both | Install and verify LLM backend |
| `aeroform launch` | Simple | Guided deploy from natural language |
| `aeroform destroy --confirm` | Simple | Tear down tracked projects |
| `aeroform init` | Pro | Scaffold config.yaml |
| `aeroform bootstrap --repo` | Pro | Generate OIDC + state backend commands |
| `aeroform plan "prompt"` | Pro | Dry-run: select templates, scan, plan |
| `aeroform generate "prompt"` | Pro | Full pipeline: scan, plan, apply |
| `aeroform scan` | Pro | Standalone security scan |
| `aeroform drift <project>` | Both | Detect infrastructure drift |
| `aeroform env add <name>` | Pro | Create environment workspace |
| `aeroform policy add <file>` | Pro | Add custom Checkov policy |
| `aeroform status` | Simple | Show tracked projects |
| `aeroform cost` | Simple | Show estimated monthly spend |

Run `aeroform --help` for the full command tree (`template`, `domain`, `open`, `logs`, `upgrade`, …).

## License

Apache 2.0 — see [LICENSE](../LICENSE) in the repository root.

## Releases

Install a specific version with:

```bash
go install github.com/momo-s15/aeroform@v1.0.5
```

Prebuilt binaries are attached to [GitHub Releases](https://github.com/momo-s15/aeroform/releases) when the release workflow succeeds.
