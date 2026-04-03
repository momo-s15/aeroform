package providers

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/cost"
	"github.com/momo-s15/aeroform/internal/engine"
)

type AWSProvider struct{}

func (AWSProvider) Name() string { return "aws" }

func (AWSProvider) Validate(cfg config.Config) error {
	if strings.TrimSpace(cfg.Cloud) == "" {
		return fmt.Errorf("cloud is required")
	}
	if strings.TrimSpace(cfg.AWS.Region) == "" {
		return fmt.Errorf("aws.region is required")
	}
	if strings.TrimSpace(cfg.AWS.AccountID) == "" {
		return fmt.Errorf("aws.account_id is required")
	}
	return nil
}

func (AWSProvider) GenerateVars(cfg config.Config, params map[string]string) map[string]string {
	vars := map[string]string{
		"cloud":          cfg.Cloud,
		"mode":           cfg.Mode,
		"aws_region":     cfg.AWS.Region,
		"aws_account_id": cfg.AWS.AccountID,
		"azure_location": cfg.Azure.Location,
		"gcp_region":     cfg.GCP.Region,
		"state_backend":  cfg.State.Backend,
		"llm_backend":    cfg.LLM.Backend,
		"llm_model":      cfg.LLM.Model,
		"llm_host":       cfg.LLM.Host,
		"llm_region":     cfg.LLM.Region,
		"llm_api_key":    cfg.LLM.APIKey,
		"llm_role_arn":   cfg.LLM.RoleARN,
	}
	for key, value := range params {
		vars[key] = value
	}
	return vars
}

func (AWSProvider) GetTemplateDir(mode engine.Mode, tmpl string) string {
	cloud := "aws"
	modeDir := "simple"
	if mode == engine.ProMode {
		modeDir = "pro"
	}
	return path.Join("templates", modeDir, cloud, tmpl)
}

func (AWSProvider) EstimateCost(templates []string) CostEstimate {
	return cost.EstimateForProTemplates("aws", templates)
}

func (AWSProvider) PostDeploy(cfg config.Config, w io.Writer) error {
	fmt.Fprintln(w, "AWS post-deploy checklist:")
	fmt.Fprintln(w, "  1. Run 'terraform output' in the project directory to see endpoints.")
	fmt.Fprintln(w, "  2. If EKS was deployed: aws eks update-kubeconfig --name <cluster> --region "+cfg.AWS.Region)
	fmt.Fprintln(w, "  3. Verify IAM roles and security groups in the AWS Console.")
	fmt.Fprintln(w, "  4. Set up CloudWatch alarms for billing and resource health.")
	return nil
}

func (AWSProvider) Region(cfg config.Config) string {
	return cfg.AWS.Region
}

func (AWSProvider) BootstrapSections(cfg config.Config, repo string) []BootstrapSection {
	region := cfg.AWS.Region
	if region == "" {
		region = "us-east-1"
	}
	accountID := cfg.AWS.AccountID
	if accountID == "" {
		accountID = "<YOUR_AWS_ACCOUNT_ID>"
	}
	if repo == "" {
		repo = "<owner/repo>"
	}

	roleName := "aeroform-github-actions"
	bucketName := fmt.Sprintf("aeroform-state-%s", accountID)
	tableName := "aeroform-state-lock"

	sections := []BootstrapSection{
		{
			Title: "Create OIDC identity provider for GitHub Actions",
			Commands: []string{
				"# One-time setup: register GitHub's OIDC thumbprint with AWS IAM",
				"aws iam create-open-id-connect-provider \\",
				"  --url https://token.actions.githubusercontent.com \\",
				"  --client-id-list sts.amazonaws.com \\",
				"  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1 \\",
				"  --region " + region,
			},
		},
		{
			Title: "Create IAM role with OIDC trust policy",
			Commands: []string{
				"# Create trust-policy.json for repo " + repo,
				`cat > /tmp/aeroform-trust-policy.json << 'POLICY'`,
				`{`,
				`  "Version": "2012-10-17",`,
				`  "Statement": [`,
				`    {`,
				`      "Effect": "Allow",`,
				`      "Principal": {`,
				`        "Federated": "arn:aws:iam::` + accountID + `:oidc-provider/token.actions.githubusercontent.com"`,
				`      },`,
				`      "Action": "sts:AssumeRoleWithWebIdentity",`,
				`      "Condition": {`,
				`        "StringEquals": {`,
				`          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"`,
				`        },`,
				`        "StringLike": {`,
				`          "token.actions.githubusercontent.com:sub": "repo:` + repo + `:*"`,
				`        }`,
				`      }`,
				`    }`,
				`  ]`,
				`}`,
				`POLICY`,
				"",
				"aws iam create-role \\",
				"  --role-name " + roleName + " \\",
				"  --assume-role-policy-document file:///tmp/aeroform-trust-policy.json",
				"",
				"# Attach AdministratorAccess (scope down for production)",
				"aws iam attach-role-policy \\",
				"  --role-name " + roleName + " \\",
				"  --policy-arn arn:aws:iam::aws:policy/AdministratorAccess",
			},
		},
		{
			Title: "Create S3 state bucket + DynamoDB lock table",
			Commands: []string{
				"aws s3api create-bucket \\",
				"  --bucket " + bucketName + " \\",
				"  --region " + region + " \\",
				"  --create-bucket-configuration LocationConstraint=" + region,
				"",
				"aws s3api put-bucket-versioning \\",
				"  --bucket " + bucketName + " \\",
				"  --versioning-configuration Status=Enabled",
				"",
				"aws s3api put-bucket-encryption \\",
				"  --bucket " + bucketName + " \\",
				`  --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'`,
				"",
				"aws s3api put-public-access-block \\",
				"  --bucket " + bucketName + " \\",
				"  --public-access-block-configuration BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true",
				"",
				"aws dynamodb create-table \\",
				"  --table-name " + tableName + " \\",
				"  --attribute-definitions AttributeName=LockID,AttributeType=S \\",
				"  --key-schema AttributeName=LockID,KeyType=HASH \\",
				"  --billing-mode PAY_PER_REQUEST \\",
				"  --region " + region,
			},
		},
		{
			Title: "Add to your Terraform backend config",
			Commands: []string{
				"# Add this block to your Terraform configuration:",
				`terraform {`,
				`  backend "s3" {`,
				`    bucket         = "` + bucketName + `"`,
				`    key            = "aeroform/terraform.tfstate"`,
				`    region         = "` + region + `"`,
				`    dynamodb_table = "` + tableName + `"`,
				`    encrypt        = true`,
				`  }`,
				`}`,
			},
		},
	}

	return sections
}

func (AWSProvider) SupportedProTemplates() []string {
	return []string{"vpc", "eks", "rds-private", "s3-private", "alb", "lambda-api", "ecs-fargate", "cloudfront-api"}
}

func (AWSProvider) SelectProTemplates(prompt string) []string {
	value := strings.ToLower(prompt)
	templates := []string{"vpc"}

	if strings.Contains(value, "eks") || strings.Contains(value, "kubernetes") {
		templates = append(templates, "eks")
	}
	if strings.Contains(value, "rds") || strings.Contains(value, "database") {
		templates = append(templates, "rds-private")
	}
	if strings.Contains(value, "s3") || strings.Contains(value, "bucket") || strings.Contains(value, "storage") {
		templates = append(templates, "s3-private")
	}
	if strings.Contains(value, "alb") || strings.Contains(value, "load balancer") {
		templates = append(templates, "alb")
	}
	if strings.Contains(value, "ecs") || strings.Contains(value, "fargate") || strings.Contains(value, "container") {
		templates = append(templates, "ecs-fargate")
	}
	if strings.Contains(value, "cdn") || strings.Contains(value, "cloudfront") {
		templates = append(templates, "cloudfront-api")
	}
	if len(templates) == 1 {
		templates = append(templates, "lambda-api")
	}

	return templates
}
