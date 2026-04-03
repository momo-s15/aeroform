# Template Reference

Aeroform ships with pre-built, security-hardened Terraform templates for AWS, Azure, and GCP. Templates are organized by mode (Simple or Pro) and cloud provider.

Every template follows these principles:

- **Encryption at rest** enabled by default
- **Public access blocked** unless the template specifically serves public content (e.g., static sites)
- **Least-privilege IAM** with scoped policies
- **Versioning and backups** enabled where applicable
- **Free-tier friendly** defaults for Simple Mode

---

## Simple Mode Templates

Simple Mode templates are single, self-contained Terraform configurations designed for one-command deployment. Each template deploys a complete, working project.

### AWS — Simple Mode (10 templates)

| Template | Description | Resources | Est. Cost |
|----------|-------------|-----------|-----------|
| `static-site` | Static website with CDN | S3 bucket (SSE, versioning, private), CloudFront distribution (HTTPS redirect, OAI), S3 bucket policy | $0.50/mo |
| `contact-form` | Static site with serverless contact form | S3 + CloudFront (same as static-site), Lambda (Node.js handler), API Gateway HTTP API (POST /contact, CORS), SES email identity, IAM role, CloudWatch logs | $0.50/mo |
| `lambda-api` | Serverless REST API with database | Lambda (CRUD handler), API Gateway HTTP API (GET/POST/DELETE), DynamoDB on-demand table (PITR, SSE), IAM role, CloudWatch logs | $0.00/mo |
| `tiny-db` | Private PostgreSQL database | Default VPC lookup, Security group (PostgreSQL port, VPC-only), RDS PostgreSQL 16.4 (t3.micro, gp3, encrypted, 7-day backup) | $14.99/mo |
| `discord-bot` | Persistent bot host | EC2 t3.micro (AL2023, IMDSv2, encrypted gp3), Elastic IP, Security group (SSH only), user_data (Node.js + systemd service) | $8.00/mo |
| `game-server` | Dedicated game server | EC2 t3.medium (AL2023, 30GB encrypted gp3), Elastic IP, Security group (configurable game port TCP+UDP + SSH), user_data (Java 21 + screen) | $30.00/mo |
| `fullstack-app` | Full-stack web application | App Runner service (configurable container image, health check), S3 assets bucket, Private RDS PostgreSQL (t3.micro), IAM role (S3 access) | $15.00/mo |
| `file-upload` | Serverless file upload service | S3 upload bucket (SSE, versioning), Lambda (presigned URL generator), API Gateway HTTP API (GET /upload-url, CORS), IAM role | $0.00/mo |
| `url-shortener` | Serverless URL shortener | Lambda (create/redirect handler), API Gateway HTTP API, DynamoDB table (TTL enabled), IAM role, CloudWatch logs | $0.00/mo |
| `cron-job` | Scheduled serverless task | Lambda function, EventBridge scheduled rule (configurable cron expression), IAM role, CloudWatch logs | $0.00/mo |

### Azure — Simple Mode (2 templates)

| Template | Description | Resources | Est. Cost |
|----------|-------------|-----------|-----------|
| `static-site` | Static website with CDN | Resource group, Storage Account (static website, TLS 1.2, versioning), CDN Standard profile, CDN endpoint (HTTPS) | $1.00/mo |
| `function-api` | Serverless API | Resource group, Storage Account, Service Plan (Linux Y1 consumption), Function App (HTTPS only, Node.js 20, FTPS disabled) | $0.00/mo |

### GCP — Simple Mode (2 templates)

| Template | Description | Resources | Est. Cost |
|----------|-------------|-----------|-----------|
| `static-site` | Static website with CDN | GCS bucket (versioning, public IAM), Backend bucket (Cloud CDN), URL map, HTTP proxy, Forwarding rule | $0.50/mo |
| `cloud-run-api` | Serverless container API | Cloud Run v2 service (scales to zero, configurable image, CPU/memory limits), IAM public invoker binding | $0.00/mo |

---

## Pro Mode Templates

Pro Mode templates are designed to be composed together. The LLM selects a set of templates based on your prompt, and Aeroform renders each as a Terraform module with a root `main.tf` that wires them together.

### AWS — Pro Mode (8 templates)

| Template | Description | Resources | Est. Cost |
|----------|-------------|-----------|-----------|
| `vpc` | Network foundation | VPC, public/private subnets (3 AZs), Internet Gateway, NAT Gateway, route tables, flow logs | $32.40/mo |
| `eks` | Kubernetes cluster | EKS cluster (private endpoint, Calico CNI, OIDC), managed node group (autoscaling, Shielded Instances), IAM roles, CloudWatch logging | $73.00/mo |
| `rds-private` | Private database | RDS PostgreSQL (encrypted gp3, PITR, private subnet group, security group), parameter group | $15.00/mo |
| `s3-private` | Encrypted object storage | S3 bucket (AES-256 SSE, versioning, public access blocked, lifecycle rules) | $0.50/mo |
| `alb` | Application load balancer | ALB (public subnets, access logs), target group (health check), security group (80/443), listener rules | $22.00/mo |
| `lambda-api` | Serverless API | Lambda + API Gateway + DynamoDB (same resources as Simple Mode) | $0.00/mo |
| `ecs-fargate` | Container service | ECS Fargate service, task definition (0.25 vCPU, 0.5GB), ALB integration, auto-scaling, CloudWatch logs | $36.00/mo |
| `cloudfront-api` | CDN with WAF | CloudFront distribution, WAFv2 web ACL (rate limiting, common rule set), API Gateway origin | $1.00/mo |

### Azure — Pro Mode (6 templates)

| Template | Description | Resources | Est. Cost |
|----------|-------------|-----------|-----------|
| `vnet` | Network foundation | Virtual Network, subnets, NSGs, service endpoints | $0.00/mo |
| `aks` | Kubernetes cluster | AKS (managed identity, private cluster, Calico, Workload Identity, OIDC, RBAC), Log Analytics, autoscaling node pool | $73.00/mo |
| `cosmos-db` | NoSQL database | Cosmos DB SQL account (private endpoint, periodic backup, configurable consistency), database, container | $25.00/mo |
| `app-service` | Web application | App Service Plan (B1 Linux), Web App (HTTPS only, managed identity, FTPS disabled, TLS 1.2, Node.js stack) | $13.14/mo |
| `storage` | Blob storage | Storage Account (TLS 1.2, no public blob access, versioning, soft delete 7d, lifecycle 30d -> cool), private container | $1.00/mo |
| `key-vault` | Secrets management | Key Vault (RBAC authorization, soft delete, purge protection, network ACLs deny-by-default, Azure Services bypass) | $0.03/mo |

### GCP — Pro Mode (4 templates)

| Template | Description | Resources | Est. Cost |
|----------|-------------|-----------|-----------|
| `vpc` | Network foundation | VPC, subnets, Cloud NAT, Cloud Router, firewall rules | $32.40/mo |
| `gke` | Kubernetes cluster | GKE (private nodes, Workload Identity, REGULAR channel, Shielded VMs, Secure Boot, auto-repair/upgrade, autoscaling) | $73.00/mo |
| `cloudsql` | Managed database | Cloud SQL PostgreSQL (private IP, SSL required, automated backups + PITR, maintenance window, SSD autoresize, HA option) | $7.67/mo |
| `gcs` | Object storage | GCS bucket (uniform access, public access prevention, versioning, lifecycle: 30d NEARLINE, 90d COLDLINE, old version cleanup) | $0.50/mo |

---

## Cost Estimates

All costs shown are estimates for the smallest viable configuration at low traffic. Actual costs depend on usage patterns, data transfer, and whether you're within the cloud provider's free tier.

- **$0.00** templates rely on free-tier allowances (e.g., Lambda 1M requests/mo, DynamoDB 25 read/write units)
- Costs do not include optional add-ons like custom domains (Route 53 hosted zone ~$0.50/mo) or data transfer overages
- Pro Mode templates composed together may share resources (e.g., VPC) which are counted once

## Adding Custom Templates

Aeroform reads templates from the `templates/` directory. To add a custom template:

1. Create a directory under `templates/simple/<cloud>/<name>/` or `templates/pro/<cloud>/<name>/`
2. Add `main.tf`, `variables.tf`, and `outputs.tf`
3. The template will appear in template selection automatically

All variables named `project_name` and `region`/`location` are auto-populated by Aeroform.
