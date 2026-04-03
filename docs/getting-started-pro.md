# Getting Started — Pro Mode

Pro Mode is for teams and professionals who need explicit control over infrastructure composition, security gating, and multi-environment workflows. You write a `config.yaml`, describe your infrastructure in natural language, and Aeroform composes multiple Terraform modules with full security scanning before deployment.

## Prerequisites

| Tool | Why | Install |
|------|-----|---------|
| **Go 1.21+** | Build and run Aeroform | [go.dev/dl](https://go.dev/dl/) |
| **Terraform** | Provisions cloud resources | [developer.hashicorp.com/terraform/install](https://developer.hashicorp.com/terraform/install) |
| **Ollama** | Local LLM for template selection | [ollama.com/download](https://ollama.com/download) |
| **Checkov** | Security scanning (required for Pro Mode gate) | `pip install checkov` |
| **tfsec** (optional) | Additional security scanning | [github.com/aquasecurity/tfsec](https://github.com/aquasecurity/tfsec) |
| **Cloud CLI** | Authentication with your cloud provider | See provider docs |

## Step 1 — Create config.yaml

Pro Mode activates when a `config.yaml` file exists in your project directory. Scaffold one from the example:

```bash
aeroform init
```

This creates `config.yaml` from the built-in example. Edit it for your cloud:

**AWS example:**

```yaml
cloud: aws
mode: pro
aws:
  region: us-east-1
  account_id: "123456789012"
llm:
  backend: ollama
  model: llama3.2
```

**Azure example:**

```yaml
cloud: azure
mode: pro
azure:
  location: eastus
llm:
  backend: ollama
  model: llama3.2
```

**GCP example:**

```yaml
cloud: gcp
mode: pro
gcp:
  region: us-central1
llm:
  backend: ollama
  model: llama3.2
```

You can override any config value with environment variables using the `AEROFORM_` prefix:

```bash
export AEROFORM_CLOUD=aws
export AEROFORM_AWS_REGION=eu-west-1
```

## Step 2 — Bootstrap Cloud Credentials (Optional)

Set up OIDC trust for GitHub Actions CI/CD and a remote Terraform state backend:

```bash
aeroform bootstrap --repo your-org/your-repo
```

Example output (AWS):

```
=== Create GitHub OIDC Identity Provider ===
  aws iam create-open-id-connect-provider \
    --url https://token.actions.githubusercontent.com \
    --client-id-list sts.amazonaws.com \
    --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1

=== Create IAM Role for GitHub Actions ===
  aws iam create-role --role-name aeroform-github-actions ...

=== Create S3 State Bucket + DynamoDB Lock Table ===
  aws s3api create-bucket --bucket aeroform-state-123456789012 ...
  aws dynamodb create-table --table-name aeroform-lock ...

=== Terraform Backend Config ===
  terraform {
    backend "s3" {
      bucket         = "aeroform-state-123456789012"
      key            = "aeroform/terraform.tfstate"
      region         = "us-east-1"
      dynamodb_table = "aeroform-lock"
      encrypt        = true
    }
  }
```

Copy and run the commands for your cloud. Bootstrap is a one-time setup.

## Step 3 — Dry-Run with Plan

Preview what Aeroform would deploy without actually creating resources:

```bash
aeroform plan "resilient kubernetes cluster with a database"
```

Example output:

```
Aeroform Pro Mode
Loading config...
  cloud: aws | region: us-east-1

Building plan...
  LLM selected templates: vpc, eks, rds-private

Estimated monthly cost:
  vpc            $32.40   NAT Gateway ($32.40/mo baseline) + data transfer
  eks            $73.00   EKS control plane ($73/mo) + node instances extra
  rds-private    $15.00   db.t3.micro with 20GB gp3
  ─────────────────────
  Total          $120.40/mo

Rendering templates into .aeroform/pro/resilient-kubernetes-cluster/
  module "vpc" -> vpc/
  module "eks" -> eks/
  module "rds_private" -> rds-private/

Security scan...
  checkov: pass
  tfsec: pass
  Security gate: pass

Terraform plan:
  Plan: 23 to add, 0 to change, 0 to destroy.
```

## Step 4 — Generate and Deploy

When satisfied with the plan, deploy:

```bash
aeroform generate "resilient kubernetes cluster with a database"
```

This runs the same pipeline as `plan` but continues to apply:

```
...
? Apply this plan? (Y/n) Y

Applying...
✓ Deployment complete!

Outputs:
  cluster_endpoint = https://ABC123.gr7.us-east-1.eks.amazonaws.com
  cluster_name     = resilient-kubernetes-cluster-eks
  db_endpoint      = resilient-kubernetes-cluster-db.abc123.us-east-1.rds.amazonaws.com

Post-deploy checklist:
  aws eks update-kubeconfig --name resilient-kubernetes-cluster-eks --region us-east-1
  ⚠ Remember to configure IAM access entries for your team.
  ⚠ Enable CloudWatch Container Insights for observability.
```

## Step 5 — Manage Environments

Create isolated environments from the same templates:

```bash
aeroform env add staging --dir .aeroform/pro/resilient-kubernetes-cluster
aeroform env add prod --dir .aeroform/pro/resilient-kubernetes-cluster
```

Switch between them:

```bash
aeroform env select staging
aeroform env list
```

Each environment has its own Terraform state — deploying to staging never touches prod.

## Step 6 — Detect Drift

Check if your live infrastructure has drifted from the Terraform state:

```bash
aeroform drift resilient-kubernetes-cluster
```

```
Drift report for resilient-kubernetes-cluster
  ✓ No drift detected — infrastructure matches desired state.
```

Or if drift is found:

```
Drift report for resilient-kubernetes-cluster
  ~ aws_security_group.eks_nodes (update in-place)
  + aws_iam_role_policy.extra (create)
  2 resources have drifted.
```

## Step 7 — Security Scanning

Run a standalone security scan on any Terraform directory:

```bash
aeroform scan
```

Add custom Checkov policies:

```bash
aeroform policy add checks/no_public_rds.py
aeroform policy list
```

Custom policies are automatically included in every scan and in the Pro Mode security gate.

## Available Pro Mode Templates

### AWS (8 templates)

| Template | What it deploys | Est. monthly cost |
|----------|----------------|-------------------|
| `vpc` | VPC + public/private subnets + NAT Gateway | $32.40 |
| `eks` | EKS cluster + managed node group + Calico | $73.00 |
| `rds-private` | Private RDS PostgreSQL + subnet group | $15.00 |
| `s3-private` | Encrypted S3 bucket + versioning | $0.50 |
| `alb` | Application Load Balancer + target group | $22.00 |
| `lambda-api` | Lambda + API Gateway + DynamoDB | $0.00 |
| `ecs-fargate` | ECS Fargate service + ALB integration | $36.00 |
| `cloudfront-api` | CloudFront + WAFv2 + API origin | $1.00 |

### Azure (6 templates)

| Template | What it deploys | Est. monthly cost |
|----------|----------------|-------------------|
| `vnet` | Virtual Network + subnets + NSGs | $0.00 |
| `aks` | AKS cluster + managed identity + RBAC | $73.00 |
| `cosmos-db` | Cosmos DB SQL + private endpoint | $25.00 |
| `app-service` | Linux App Service Plan + Web App | $13.14 |
| `storage` | Storage Account + lifecycle policies | $1.00 |
| `key-vault` | Key Vault + RBAC + purge protection | $0.03 |

### GCP (4 templates)

| Template | What it deploys | Est. monthly cost |
|----------|----------------|-------------------|
| `vpc` | VPC + Cloud NAT + private subnets | $32.40 |
| `gke` | GKE cluster + Workload Identity + Shielded VMs | $73.00 |
| `cloudsql` | Cloud SQL PostgreSQL + private IP + backups | $7.67 |
| `gcs` | GCS bucket + lifecycle rules + versioning | $0.50 |

## How Pro Mode Template Composition Works

Unlike Simple Mode (which deploys a single template), Pro Mode composes multiple templates into a single Terraform project. The LLM reads your prompt and selects a set of templates that work together. For example:

- "kubernetes cluster with a database" -> `vpc` + `eks` + `rds-private`
- "serverless API with CDN" -> `lambda-api` + `cloudfront-api`
- "web app with load balancer" -> `vpc` + `alb` + `ecs-fargate` + `rds-private`

Each template is rendered as a Terraform module, and a root `main.tf` is generated that wires them together.

## LLM Backend Options

Pro Mode supports three LLM backends. Configure in `config.yaml`:

```yaml
# Local (default)
llm:
  backend: ollama
  model: llama3.2

# OpenAI
llm:
  backend: openai
  api_key: sk-...

# AWS Bedrock
llm:
  backend: bedrock
  region: us-east-1
  role_arn: arn:aws:iam::123456789012:role/bedrock-role
```

See [LLM Backends](llm-backends.md) for detailed setup instructions.

## Next Steps

- Set up [OIDC bootstrap](../bootstrap) for CI/CD pipelines
- Add [custom security policies](security-model.md) for your organization
- Read the [Architecture Decisions](architecture-decisions.md) to understand the design
