//go:build integration

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// writeLocalStackAWSProviderOverride writes a *_override.tf provider block so Terraform
// talks to LocalStack in CI instead of real AWS (templates only set region by default).
func writeLocalStackAWSProviderOverride(t *testing.T, workDir string) {
	t.Helper()
	ep := os.Getenv("AWS_ENDPOINT_URL")
	if ep == "" {
		ep = "http://localhost:4566"
	}
	// Broad service map for simple + pro AWS modules under LocalStack.
	const tmpl = `provider "aws" {
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway   = %[1]q
    apigatewayv2 = %[1]q
    cloudfront   = %[1]q
    cloudwatch   = %[1]q
    dynamodb     = %[1]q
    ec2          = %[1]q
    elbv2        = %[1]q
    iam          = %[1]q
    lambda       = %[1]q
    logs         = %[1]q
    rds          = %[1]q
    s3           = %[1]q
    sns          = %[1]q
    sqs          = %[1]q
    sts          = %[1]q
  }
}
`
	content := fmt.Sprintf(tmpl, ep)
	path := filepath.Join(workDir, "localstack_aws_override.tf")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write localstack aws override: %v", err)
	}
}
