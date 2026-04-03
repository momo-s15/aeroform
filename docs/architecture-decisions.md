# Architecture Decisions

This document explains the key design decisions behind Aeroform and why they were made.

---

## Why Two Modes?

Cloud infrastructure tools tend to target one audience: either beginners (with heavy abstraction and hidden complexity) or professionals (with steep learning curves and configuration overload). Aeroform serves both from a single CLI binary.

**Simple Mode** is for students, solo developers, and first-time cloud users. It uses natural language, guided prompts, cost estimates in dollars, and plain-English security reports. The user never sees raw Terraform output unless they want to.

**Pro Mode** is for teams, DevOps engineers, and production workloads. It uses explicit `config.yaml`, multi-template composition, Terraform module rendering, Checkov/tfsec gating, and workspace-based environments.

Both modes share the same engine, templates, security pipeline, and provider abstraction. The difference is in the user interface and the level of automation.

### Why not separate tools?

- Shared template library means improvements benefit both audiences
- Users can graduate from Simple to Pro without switching tools (`aeroform upgrade`)
- One binary to install, one set of docs, one CI pipeline

---

## Why Constrained Generation?

Aeroform uses an LLM, but the LLM never generates Terraform code. Instead, it selects from a library of pre-validated templates. This is the **constrained generation** pattern.

### The problem with unconstrained generation

Tools that let an LLM write arbitrary HCL face serious issues:

1. **Security** — LLMs frequently generate insecure configurations (public S3 buckets, open security groups, no encryption)
2. **Correctness** — Generated Terraform often has syntax errors, missing required arguments, or invalid resource combinations
3. **Reproducibility** — The same prompt can produce different infrastructure on different runs
4. **Auditability** — You can't review a template library if there is no template library

### How constrained generation works

1. The user describes what they want in natural language
2. The LLM receives the description and a list of available template names with descriptions
3. The LLM returns one or more template names (Simple Mode: 1 template, Pro Mode: multiple)
4. Aeroform validates the response against the known template list
5. If validation fails, keyword matching provides a deterministic fallback
6. The selected template(s) are rendered with the user's variables

The LLM is a **classification** layer, not a **generation** layer. This keeps the security and correctness properties of hand-written templates while providing the natural language interface of AI tools.

### Trade-offs

- **Pro**: Every deployment uses reviewed, tested, security-hardened Terraform
- **Pro**: Deterministic — same template always produces the same infrastructure
- **Pro**: Small local models work well for classification tasks
- **Con**: Limited to templates in the library (can't deploy arbitrary infrastructure)
- **Con**: Adding new deployment patterns requires writing a new template

This trade-off is intentional. Aeroform optimizes for safety and reliability over flexibility.

---

## Why Go?

Aeroform is written in Go for several reasons:

1. **Single binary** — `go build` produces a statically-linked binary with no runtime dependencies. Users don't need Python, Node.js, or Docker to run Aeroform.
2. **Cross-platform** — Go cross-compiles to Linux, macOS, and Windows from a single codebase. GoReleaser builds all six binaries in CI.
3. **Cloud SDK ecosystem** — AWS, Azure, and GCP all have first-class Go SDKs.
4. **Terraform ecosystem** — Terraform itself is written in Go. The HCL library and Terraform provider SDKs are Go-native.
5. **Performance** — Template rendering, file I/O, and HTTP calls to the LLM backend are fast without async complexity.
6. **Simplicity** — Go's straightforward concurrency model, strong typing, and standard library reduce the surface area for bugs.

### Why not Python/TypeScript/Rust?

- **Python**: Would require users to manage virtual environments. Packaging as a single binary is possible (PyInstaller) but fragile.
- **TypeScript**: Requires Node.js runtime. Good for web tools, less natural for CLI + cloud SDK work.
- **Rust**: Would work well for the binary story, but the cloud SDK ecosystem is less mature, and contribution barrier is higher for an open-source project.

---

## Why Terraform?

Aeroform uses Terraform as its infrastructure provisioning engine rather than CloudFormation, Pulumi, CDK, or direct API calls.

1. **Multi-cloud** — Terraform works with AWS, Azure, and GCP from the same tool. CloudFormation is AWS-only. CDK generates CloudFormation.
2. **Declarative** — Templates describe desired state, not imperative steps. This makes them reviewable and auditable.
3. **State management** — Terraform tracks what it created and can destroy it cleanly. This is critical for Simple Mode's `aeroform destroy`.
4. **Ecosystem** — Thousands of providers and modules. Well-understood by the industry.
5. **Installable** — Single binary, same as Go. No runtime dependencies.

Aeroform wraps Terraform rather than replacing it. Users who outgrow Aeroform can take the rendered Terraform files and manage them directly.

---

## Why Checkov?

Checkov is the primary security scanner for several reasons:

1. **Terraform-native** — Built specifically for scanning Terraform, CloudFormation, Kubernetes, and other IaC formats.
2. **Comprehensive** — 1000+ built-in checks covering CIS benchmarks, SOC2, HIPAA, and cloud-specific best practices.
3. **Custom policies** — Supports Python and YAML custom checks, which Aeroform exposes via `aeroform policy add`.
4. **Exit codes** — Clean exit code behavior (0 = pass, 1 = findings) makes it easy to use as a gate.
5. **Open source** — Apache 2.0 licensed, same as Aeroform.

tfsec is included as a secondary scanner because it catches some cloud-specific issues that Checkov misses, and vice versa. Both must pass for the Pro Mode gate to open.

---

## Provider Abstraction

Cloud-specific behavior is isolated behind the `CloudProvider` interface:

```go
type CloudProvider interface {
    Name() string
    GetTemplateDir(mode, template string) string
    GenerateVars(cfg config.Config, extra map[string]string) map[string]string
    EstimateCost(templates []string) CostEstimate
    PostDeploy(out io.Writer, templates []string)
    BootstrapSections(cfg config.Config, repo string) []BootstrapSection
}
```

Each cloud (AWS, Azure, GCP) implements this interface. The engine, commands, and security pipeline work with the interface, never with cloud-specific types directly. This means:

- Adding a new cloud provider is a single file implementing the interface
- The CLI commands don't change when a new cloud is added
- Templates are organized by `templates/{mode}/{cloud}/{template}/`

---

## Project Structure

```
aeroform/
├── cmd/                    # Cobra command definitions
├── internal/
│   ├── config/             # Viper-based config loading + validation
│   ├── cost/               # Cost estimation (Simple + Pro, per-cloud)
│   ├── engine/             # Simple Mode planner, Pro Mode planner, drift
│   ├── llm/                # LLM backends (Ollama, OpenAI, Bedrock)
│   ├── logger/             # Zap structured logging
│   ├── prompt/             # Promptui-based interactive prompts
│   ├── providers/          # CloudProvider implementations (AWS, Azure, GCP)
│   ├── security/           # Checkov/tfsec scanning, auto-correction, policy mgmt
│   ├── setup/              # Ollama detection and model pulling
│   ├── simplestate/        # JSON-based project tracking for Simple Mode
│   ├── terraform/          # Template renderer, Terraform runner, workspace mgmt
│   └── ui/                 # Colored terminal output helpers
├── bootstrap/              # Cloud bootstrap script generation
├── templates/
│   ├── simple/{cloud}/     # Simple Mode Terraform templates
│   └── pro/{cloud}/        # Pro Mode Terraform templates
├── docs/                   # Documentation site
├── test/
│   ├── integration/        # LocalStack-based integration tests
│   └── e2e/                # Real cloud E2E tests (release-only)
└── .github/workflows/      # CI, release, E2E pipelines
```

---

## Security Pipeline Design

The security pipeline is intentionally ordered:

1. **Auto-correction** (Simple Mode only) — Apply safe fixes before scanning
2. **Checkov scan** — Comprehensive IaC security analysis
3. **tfsec scan** (Pro Mode only) — Additional cloud-specific checks
4. **Gate decision** — Block on HIGH severity (Simple) or any failure (Pro)
5. **Terraform plan** — Show what will change
6. **User confirmation** — Explicit approval before apply

This ordering ensures that:
- Auto-corrected templates get scanned (no bypassing the scanner)
- Users see the plan after security passes (no wasted time on blocked deployments)
- Nothing is applied without explicit human confirmation

---

## Cost Estimation Design

Cost estimates use a lookup table approach rather than real-time pricing APIs:

- **Predictable** — Same template always shows the same estimate
- **Fast** — No network calls, no API keys needed
- **Offline** — Works without internet access
- **Honest** — Estimates are for the smallest viable configuration, clearly labeled

The trade-off is that estimates may drift from actual pricing over time. This is acceptable because:
- Aeroform targets small-scale, free-tier-friendly deployments
- Exact pricing requires usage patterns that aren't known at plan time
- The estimates serve as a ballpark guide, not a billing commitment

All cost arithmetic uses `shopspring/decimal` to avoid floating-point rounding errors.

---

## CI, linting, and supply chain

- **Go toolchain** is pinned via `go.mod` / `toolchain` so GitHub Actions, `govulncheck`, and contributors use a consistent compiler (currently Go 1.26.1+).
- **Lint** uses **golangci-lint v2** (`.golangci.yml` `version: "2"`) so linters stay compatible with the supported Go release.
- **Integration tests** exercise **Terraform plan** against **LocalStack** in CI (pinned image) so AWS-shaped templates are validated without a live account.
- **E2E tests** are optional at release time: they run on `v*` tag pushes when AWS OIDC secrets are configured; otherwise the workflow skips after detecting missing configuration.
- **Release builds** are produced by the tag-triggered workflow (multi-platform binaries uploaded to GitHub Releases).

This keeps fast feedback on every PR while still allowing opt-in real-cloud validation for maintainers who wire AWS credentials.
