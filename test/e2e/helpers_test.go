//go:build e2e

package e2e

import (
	"os"
	"testing"
)

func requireAWSCredentials(t *testing.T) {
	t.Helper()
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" && os.Getenv("AWS_PROFILE") == "" && os.Getenv("AWS_ROLE_ARN") == "" {
		t.Skip("skipping E2E test: no AWS credentials configured (set AWS_ACCESS_KEY_ID or AWS_PROFILE)")
	}
}

func awsRegion() string {
	if r := os.Getenv("AWS_REGION"); r != "" {
		return r
	}
	if r := os.Getenv("AWS_DEFAULT_REGION"); r != "" {
		return r
	}
	return "us-east-1"
}

func awsAccountID(t *testing.T) string {
	t.Helper()
	if id := os.Getenv("AWS_ACCOUNT_ID"); id != "" {
		return id
	}
	return "000000000000"
}
