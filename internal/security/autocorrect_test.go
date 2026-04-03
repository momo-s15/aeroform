package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTF(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readTF(t *testing.T, dir, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestAutoCorrectPubliclyAccessible(t *testing.T) {
	dir := t.TempDir()
	writeTF(t, dir, "main.tf", `resource "aws_db_instance" "this" {
  publicly_accessible = true
  storage_encrypted   = true
}`)

	corrections, err := AutoCorrect(dir)
	if err != nil {
		t.Fatalf("AutoCorrect() error = %v", err)
	}
	if len(corrections) != 1 {
		t.Fatalf("expected 1 correction, got %d", len(corrections))
	}
	if corrections[0].Rule != "no-public-db" {
		t.Fatalf("expected rule no-public-db, got %q", corrections[0].Rule)
	}

	content := readTF(t, dir, "main.tf")
	if strings.Contains(content, "publicly_accessible = true") {
		t.Fatal("publicly_accessible = true should have been fixed")
	}
	if !strings.Contains(content, "publicly_accessible = false") {
		t.Fatal("expected publicly_accessible = false")
	}
}

func TestAutoCorrectHTTPtoHTTPS(t *testing.T) {
	dir := t.TempDir()
	writeTF(t, dir, "main.tf", `resource "aws_lb_listener" "front" {
  protocol = "HTTP"
  port     = 80
}`)

	corrections, err := AutoCorrect(dir)
	if err != nil {
		t.Fatalf("AutoCorrect() error = %v", err)
	}
	if len(corrections) != 1 {
		t.Fatalf("expected 1 correction, got %d", len(corrections))
	}
	if corrections[0].Rule != "enforce-https" {
		t.Fatalf("expected rule enforce-https, got %q", corrections[0].Rule)
	}

	content := readTF(t, dir, "main.tf")
	if !strings.Contains(content, `protocol = "HTTPS"`) {
		t.Fatal("expected protocol to be changed to HTTPS")
	}
}

func TestAutoCorrectUnencryptedStorage(t *testing.T) {
	dir := t.TempDir()
	writeTF(t, dir, "main.tf", `resource "aws_db_instance" "this" {
  storage_encrypted = false
}`)

	corrections, err := AutoCorrect(dir)
	if err != nil {
		t.Fatalf("AutoCorrect() error = %v", err)
	}
	if len(corrections) != 1 {
		t.Fatalf("expected 1 correction, got %d", len(corrections))
	}
	if corrections[0].Rule != "enforce-encryption" {
		t.Fatalf("expected rule enforce-encryption, got %q", corrections[0].Rule)
	}

	content := readTF(t, dir, "main.tf")
	if !strings.Contains(content, "storage_encrypted = true") {
		t.Fatal("expected storage_encrypted = true")
	}
}

func TestAutoCorrectMultipleFixes(t *testing.T) {
	dir := t.TempDir()
	writeTF(t, dir, "main.tf", `resource "aws_db_instance" "this" {
  publicly_accessible = true
  storage_encrypted = false
}`)

	corrections, err := AutoCorrect(dir)
	if err != nil {
		t.Fatalf("AutoCorrect() error = %v", err)
	}
	if len(corrections) != 2 {
		t.Fatalf("expected 2 corrections, got %d", len(corrections))
	}

	content := readTF(t, dir, "main.tf")
	if strings.Contains(content, "publicly_accessible = true") {
		t.Fatal("publicly_accessible should have been fixed")
	}
	if strings.Contains(content, "storage_encrypted = false") {
		t.Fatal("storage_encrypted should have been fixed")
	}
}

func TestAutoCorrectNoChanges(t *testing.T) {
	dir := t.TempDir()
	writeTF(t, dir, "main.tf", `resource "aws_s3_bucket" "this" {
  bucket = "my-bucket"
}`)

	corrections, err := AutoCorrect(dir)
	if err != nil {
		t.Fatalf("AutoCorrect() error = %v", err)
	}
	if len(corrections) != 0 {
		t.Fatalf("expected 0 corrections, got %d", len(corrections))
	}
}

func TestAutoCorrectSkipsNonTFFiles(t *testing.T) {
	dir := t.TempDir()
	writeTF(t, dir, "readme.md", `publicly_accessible = true`)
	writeTF(t, dir, "meta.yaml", `storage_encrypted = false`)

	corrections, err := AutoCorrect(dir)
	if err != nil {
		t.Fatalf("AutoCorrect() error = %v", err)
	}
	if len(corrections) != 0 {
		t.Fatalf("non-.tf files should not be modified, got %d corrections", len(corrections))
	}
}

func TestAutoCorrectEmptyDir(t *testing.T) {
	dir := t.TempDir()
	corrections, err := AutoCorrect(dir)
	if err != nil {
		t.Fatalf("AutoCorrect() error = %v", err)
	}
	if len(corrections) != 0 {
		t.Fatalf("expected 0 corrections for empty dir, got %d", len(corrections))
	}
}
