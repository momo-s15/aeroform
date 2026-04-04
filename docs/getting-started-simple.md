# Getting Started — Simple Mode

Simple Mode is the guided path for students, solo developers, and anyone deploying their first cloud project. You describe what you want in plain English, and Aeroform picks the right template, estimates costs, runs a security check, and deploys — all in one command.

## Prerequisites

| Tool | Why | Install |
|------|-----|---------|
| **Terraform** | Provisions cloud resources | [developer.hashicorp.com/terraform/install](https://developer.hashicorp.com/terraform/install) |
| **Ollama** | Local LLM for template selection | [ollama.com/download](https://ollama.com/download) |
| **AWS / Azure / GCP CLI** | Cloud authentication | See your cloud's docs |

You do **not** need Go installed if you use the **install script** (macOS/Linux). To build from source or use `go install`, see [Go 1.26.1+](https://go.dev/dl/) on `go.mod`.

On **Windows**, use the `.exe` from [Releases](https://github.com/momo-s15/aeroform/releases) (see [installation.md](installation.md)); arrow-key menus fall back to **typed answers** in some terminals. For **GCP**, you can set **`AEROFORM_GCP_PROJECT_ID`** if the project prompt is awkward.

Make sure you are authenticated with your chosen cloud provider before running Aeroform. For AWS, that means `aws configure` or environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`).

For **Azure**, `aeroform launch` asks for an **Azure region**. Some subscriptions (including Azure for Education) **block** regions like `eastus` via policy; if storage creation fails with `RequestDisallowedByAzure`, pick an allowed region (often **`canadacentral`** or **`canadaeast`** in Canada). You can list names with `az account list-locations -o table`. Set a default with `AEROFORM_AZURE_LOCATION` to skip retyping.

## Install Aeroform

**Recommended (macOS / Linux, no Go):**

```bash
curl -fsSL https://raw.githubusercontent.com/momo-s15/aeroform/main/install.sh | sh
```

See [installation.md](installation.md) for version pinning, custom paths, Windows, and [Homebrew](homebrew-tap.md).

**From source (contributors):**

```bash
git clone https://github.com/momo-s15/aeroform.git
cd aeroform
go build -o aeroform .
```

**Go install:**

```bash
go install github.com/momo-s15/aeroform@latest
```

## Step 1 — Run Setup

Setup installs and verifies the local AI backend, then pulls the default model.

```bash
aeroform setup
```

Example output:

```
Checking Ollama...
  ✓ Ollama is installed (version 0.5.4)
  ✓ Ollama is running
Pulling model llama3.2...
  ✓ Model llama3.2 ready
Setup complete.
```

You can specify a different model with `--model`:

```bash
aeroform setup --model mistral
```

## Step 2 — Launch Your First Project

```bash
aeroform launch
```

Aeroform will walk you through an interactive flow:

```
? Choose a cloud provider: (use arrow keys)
  > aws
    azure
    gcp

? What do you want to launch?
> a portfolio website

? Project name:
> my-portfolio
```

## Step 3 — Review the Plan

Aeroform shows you exactly what will be deployed and what it will cost:

```
Template: static-site
Provider: aws

Resources:
  + aws_s3_bucket.site
  + aws_cloudfront_distribution.cdn
  + aws_cloudfront_origin_access_identity.oai
  + aws_s3_bucket_policy.site_policy

Estimated monthly cost: $0.50
  S3 and CloudFront are free tier friendly; Route 53 adds a small hosted zone cost

Security check: ✓ no issues found

Plan: 4 to add, 0 to change, 0 to destroy.
? Apply this plan? (Y/n)
```

## Step 4 — Deploy

Type `Y` to confirm. Aeroform runs `terraform apply` and shows the outputs:

```
✓ Deployment complete!

Outputs:
  cloudfront_url    = https://d1234abcdef.cloudfront.net
  s3_bucket_name    = my-portfolio-site
  upload_command    = aws s3 sync ./public s3://my-portfolio-site
```

Upload your site files using the command shown in the output.

## Step 5 — Manage Your Project

Check status of tracked projects:

```bash
aeroform status
```

See estimated costs:

```bash
aeroform cost
```

Tear it down when you're done:

```bash
aeroform destroy --confirm
```

## Available Simple Mode Templates

### AWS (10 templates)

| Template | What it deploys | Est. monthly cost |
|----------|----------------|-------------------|
| `static-site` | S3 + CloudFront CDN | $0.50 |
| `contact-form` | Static site + Lambda + API Gateway + SES | $0.50 |
| `lambda-api` | Lambda + API Gateway + DynamoDB | $0.00 |
| `tiny-db` | Private RDS PostgreSQL (t3.micro) | $14.99 |
| `discord-bot` | EC2 t3.micro + Elastic IP | $8.00 |
| `game-server` | EC2 t3.medium + 30GB EBS | $30.00 |
| `fullstack-app` | App Runner + RDS + S3 | $15.00 |
| `file-upload` | S3 + Lambda + presigned URLs | $0.00 |
| `url-shortener` | Lambda + API Gateway + DynamoDB | $0.00 |
| `cron-job` | Lambda + EventBridge schedule | $0.00 |

### Azure (2 templates)

| Template | What it deploys | Est. monthly cost |
|----------|----------------|-------------------|
| `static-site` | Storage Account static website (HTTPS) | $0.50 |
| `function-api` | Azure Functions (consumption plan) | $0.00 |

### GCP (2 templates)

| Template | What it deploys | Est. monthly cost |
|----------|----------------|-------------------|
| `static-site` | GCS bucket + Cloud CDN + HTTP LB | $0.50 |
| `cloud-run-api` | Cloud Run v2 (scales to zero) | $0.00 |

## How Template Selection Works

You don't pick templates directly — you describe what you want in plain English. The local LLM analyzes your prompt and selects the best matching template from the list above. For example:

- "a portfolio website" -> `static-site`
- "a REST API for my mobile app" -> `lambda-api`
- "a Minecraft server" -> `game-server`
- "a contact form for my business" -> `contact-form`

## Security in Simple Mode

Every deployment goes through a security gate before Terraform runs:

1. **Auto-correction** — Safe fixes are applied automatically (e.g., enabling encryption, blocking public database access, enforcing HTTPS).
2. **Checkov scan** — If Checkov is installed, a full scan runs on the rendered template.
3. **Blocking** — High-severity findings block the deployment with a clear explanation.

You never need to think about security best practices — Aeroform enforces them by default.

## Custom Domains

After deploying a `static-site`, you can record a custom domain for the **most recently tracked** Simple Mode project:

```bash
aeroform domain add example.com
```

## Debugging

Add `--debug` to any command to see detailed logs on stderr:

```bash
aeroform launch --debug
```

## Next Steps

- Try deploying a different template type (API, database, full-stack app)
- Explore [Pro Mode](getting-started-pro.md) for team workflows and multi-template compositions
- Read the [Template Reference](template-reference.md) for detailed resource listings
