//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/momo-s15/aeroform/internal/security"
)

func TestCheckovGateBlocksBadTemplate(t *testing.T) {
	workDir := t.TempDir()

	badTF := `
resource "aws_s3_bucket" "bad" {
  bucket = "totally-public-bucket"
}

resource "aws_s3_bucket_public_access_block" "bad" {
  bucket                  = aws_s3_bucket.bad.id
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}
`
	if err := os.WriteFile(filepath.Join(workDir, "main.tf"), []byte(badTF), 0o644); err != nil {
		t.Fatalf("write bad template: %v", err)
	}

	report := security.RunProSecurityScan(workDir)
	t.Logf("scan found %d findings, blocking=%v", len(report.SummaryLines()), report.Blocking)

	for _, line := range report.SummaryLines() {
		t.Log(line)
	}
}

func TestCheckovGatePassesGoodTemplate(t *testing.T) {
	workDir := t.TempDir()

	goodTF := `
resource "aws_s3_bucket" "good" {
  bucket = "secure-private-bucket"
}

resource "aws_s3_bucket_public_access_block" "good" {
  bucket                  = aws_s3_bucket.good.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "good" {
  bucket = aws_s3_bucket.good.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "good" {
  bucket = aws_s3_bucket.good.id
  versioning_configuration {
    status = "Enabled"
  }
}
`
	if err := os.WriteFile(filepath.Join(workDir, "main.tf"), []byte(goodTF), 0o644); err != nil {
		t.Fatalf("write good template: %v", err)
	}

	report := security.RunProSecurityScan(workDir)
	t.Logf("scan found %d lines, blocking=%v", len(report.SummaryLines()), report.Blocking)

	for _, line := range report.SummaryLines() {
		t.Log(line)
	}
}
