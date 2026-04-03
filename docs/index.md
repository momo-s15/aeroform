# Aeroform Documentation

Aeroform is an open-source CLI tool that deploys cloud infrastructure using natural language. It supports AWS, Azure, and GCP through pre-validated, security-hardened Terraform templates.

**Two modes, one tool:**

- **Simple Mode** — Describe what you want, get it deployed. Guided prompts, cost estimates, auto-security. Perfect for students and solo developers.
- **Pro Mode** — Config-driven multi-template composition with Checkov/tfsec gating, workspace environments, and drift detection. Built for teams.

## Quick Start

```bash
# Install
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
| [Getting Started — Simple](getting-started-simple.md) | Students, solo devs | Install to live website in 5 minutes |
| [Getting Started — Pro](getting-started-pro.md) | Teams, DevOps | Config to EKS cluster with security gating |

## References

| Page | What it covers |
|------|----------------|
| [Template Reference](template-reference.md) | Every template (Simple + Pro, all clouds) with resources and costs |
| [LLM Backends](llm-backends.md) | Setup for Ollama, OpenAI, and AWS Bedrock |
| [Security Model](security-model.md) | Scanning, auto-correction, custom policies, drift detection |
| [Architecture Decisions](architecture-decisions.md) | Why two modes, constrained generation, Go, Terraform, Checkov |

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

## License

Apache 2.0 — see [LICENSE](../LICENSE) in the repository root.
