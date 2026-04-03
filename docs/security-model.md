# Security Model

Security is a first-class concern in Aeroform. Every deployment — Simple or Pro — passes through a security gate before any infrastructure is created. The security behavior is mode-aware: beginner-friendly in Simple Mode, strict and technical in Pro Mode.

---

## Core Principles

1. **Secure by default** — Templates are pre-hardened with encryption, private networking, and least-privilege IAM.
2. **Scan before deploy** — No Terraform applies without a security check passing first.
3. **Auto-correct when safe** — Simple Mode automatically fixes common misconfigurations.
4. **Block when unsafe** — High-severity findings stop the deployment with a clear explanation.
5. **Transparent** — All findings are shown to the user with severity and description.

---

## Simple Mode Security

Simple Mode targets users who may not know security best practices. The security pipeline has three stages:

### Stage 1: Auto-Correction

After the template is rendered but before the security scan, Aeroform applies safe automatic fixes:

| Rule | What it fixes | Example |
|------|---------------|---------|
| `no-public-db` | Ensures databases are not publicly accessible | Sets `publicly_accessible = false` on RDS |
| `enforce-https` | Ensures HTTPS/TLS is enforced | Sets `viewer_protocol_policy = "redirect-to-https"` on CloudFront |
| `enforce-encryption` | Ensures storage encryption is enabled | Sets `server_side_encryption_configuration` on S3 buckets |

Auto-correction only applies changes that are unambiguously safe. If a fix would change the user's intended behavior, it is left for the security gate to flag.

### Stage 2: Checkov Scan

If Checkov is installed, a full scan runs on the rendered Terraform. Findings are translated into human-readable messages:

```
Security check
  ✗ [HIGH] S3 bucket should have access logging enabled
  ⚠ [MEDIUM] CloudFront distribution should use TLS 1.2
  ℹ [LOW] S3 bucket should have lifecycle configuration
```

### Stage 3: Blocking Decision

- **High-severity findings** block the deployment. The user sees a clear explanation and the deployment stops.
- **Medium and low findings** are shown as warnings but do not block.
- **No findings** shows a green checkmark:

```
✓ Security check: no issues found
```

### Simple Mode Security Report Format

```
Security check
  ✗ [HIGH] Public access to database is not allowed
  ⚠ [MEDIUM] Consider enabling access logging
  ℹ [LOW] Add lifecycle rules for cost optimization

  1 high, 1 medium, 1 low
  Blocked: security gate requires all HIGH issues resolved
```

---

## Pro Mode Security

Pro Mode targets teams who need full scanner output and strict gating. The security pipeline runs two independent scanners.

### Scanners

| Scanner | Role | Gate behavior |
|---------|------|--------------|
| **Checkov** | Primary scanner, comprehensive Terraform checks | Blocks on any failure (non-zero exit) |
| **tfsec** | Secondary scanner, focused on cloud-specific risks | Blocks on any failure (non-zero exit) |

Both scanners must pass for the security gate to open. If either scanner is not installed, it is skipped (with a warning).

### Custom Policies

You can add organization-specific Checkov policies:

```bash
aeroform policy add checks/no_public_rds.py
aeroform policy add rules/naming_convention.yaml
```

Custom policies are stored in `.aeroform/policies/` and automatically passed to Checkov via `--external-checks-dir` on every scan.

Manage policies:

```bash
aeroform policy list      # show registered policies
aeroform policy remove no_public_rds.py  # remove a policy
```

### Supported Policy Formats

- **Python (.py)** — Checkov Python-based custom checks using the `BaseCheck` class
- **YAML (.yaml/.yml)** — Checkov YAML-based custom checks for simple attribute checks

Example Python policy:

```python
from checkov.terraform.checks.resource.base_resource_check import BaseResourceCheck
from checkov.common.models.enums import CheckResult, CheckCategories

class NoPublicRDS(BaseResourceCheck):
    def __init__(self):
        name = "Ensure RDS is not publicly accessible"
        id = "CUSTOM_RDS_001"
        supported_resources = ["aws_db_instance"]
        categories = [CheckCategories.NETWORKING]
        super().__init__(name=name, id=id, categories=categories,
                         supported_resources=supported_resources)

    def scan_resource_conf(self, conf):
        if conf.get("publicly_accessible") == [True]:
            return CheckResult.FAILED
        return CheckResult.PASSED

check = NoPublicRDS()
```

Example YAML policy:

```yaml
metadata:
  id: "CUSTOM_S3_001"
  name: "Ensure S3 bucket names follow naming convention"
  category: "CONVENTION"
definition:
  cond_type: "attribute"
  resource_types:
    - "aws_s3_bucket"
  attribute: "bucket"
  operator: "regex_match"
  value: "^(dev|staging|prod)-.*"
```

### Pro Mode Scan Output

```
Aeroform Pro Mode scan
  checkov: pass
  tfsec: pass
  Security gate: pass
```

Or when blocked:

```
Aeroform Pro Mode scan
  checkov: fail (exit 1)
  tfsec: pass
  Security gate: blocked
```

### Standalone Scanning

Run a security scan without deploying:

```bash
aeroform scan
```

This runs Checkov and tfsec against the current directory's Terraform files.

---

## Template Security Hardening

All built-in templates follow these security defaults:

### Storage
- Server-side encryption (AES-256 or KMS)
- Public access blocked (except static site buckets that serve public content)
- Versioning enabled
- Soft delete / lifecycle rules where supported

### Databases
- Private subnets only (`publicly_accessible = false`)
- Encrypted storage
- Automated backups with point-in-time recovery
- TLS/SSL required for connections
- Security groups scoped to VPC CIDR

### Compute
- IMDSv2 required (AWS EC2)
- Encrypted EBS volumes
- Minimal security group rules (SSH only where needed)
- Managed identities over static credentials

### Kubernetes
- Private cluster endpoints
- Workload Identity / IRSA for pod-level auth
- Network policies (Calico)
- Shielded VMs with Secure Boot
- RBAC enabled, legacy ABAC disabled

### Networking
- HTTPS enforced on all public endpoints
- TLS 1.2 minimum
- WAF on public-facing CDNs (Pro Mode)
- VPC flow logs enabled

---

## Drift Detection

Aeroform can detect when live infrastructure has drifted from its Terraform state:

```bash
aeroform drift my-project
```

This runs `terraform plan -detailed-exitcode` against the project directory and reports:
- **No drift** — infrastructure matches desired state
- **Drift detected** — lists each changed resource with its change type (create, update, destroy, replace)

Drift detection helps catch manual console changes, external modifications, or stale state.
