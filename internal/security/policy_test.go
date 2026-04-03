package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddPolicyPython(t *testing.T) {
	origDir := PolicyDir
	t.Cleanup(func() { PolicyDir = origDir })

	tmp := t.TempDir()
	PolicyDir = filepath.Join(tmp, "policies")

	src := filepath.Join(tmp, "my_check.py")
	if err := os.WriteFile(src, []byte("from checkov.common.models.enums import CheckResult\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	dst, err := AddPolicy(src)
	if err != nil {
		t.Fatalf("AddPolicy() error = %v", err)
	}
	if filepath.Base(dst) != "my_check.py" {
		t.Errorf("dst = %q, want my_check.py", dst)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read copied file: %v", err)
	}
	if len(content) == 0 {
		t.Error("copied file is empty")
	}
}

func TestAddPolicyYAML(t *testing.T) {
	origDir := PolicyDir
	t.Cleanup(func() { PolicyDir = origDir })

	tmp := t.TempDir()
	PolicyDir = filepath.Join(tmp, "policies")

	src := filepath.Join(tmp, "my_rule.yaml")
	if err := os.WriteFile(src, []byte("metadata:\n  name: test\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	dst, err := AddPolicy(src)
	if err != nil {
		t.Fatalf("AddPolicy() error = %v", err)
	}
	if filepath.Base(dst) != "my_rule.yaml" {
		t.Errorf("dst = %q, want my_rule.yaml", dst)
	}
}

func TestAddPolicyRejectsUnsupportedExtension(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "policy.txt")
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := AddPolicy(src)
	if err == nil {
		t.Fatal("expected error for .txt file")
	}
}

func TestAddPolicyRejectsDirectory(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "subdir.py")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := AddPolicy(dir)
	if err == nil {
		t.Fatal("expected error for directory")
	}
}

func TestAddPolicyRejectsMissingFile(t *testing.T) {
	_, err := AddPolicy("/nonexistent/path/check.py")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestListPoliciesEmpty(t *testing.T) {
	origDir := PolicyDir
	t.Cleanup(func() { PolicyDir = origDir })
	PolicyDir = filepath.Join(t.TempDir(), "policies")

	policies, err := ListPolicies()
	if err != nil {
		t.Fatalf("ListPolicies() error = %v", err)
	}
	if len(policies) != 0 {
		t.Errorf("expected 0 policies, got %d", len(policies))
	}
}

func TestListPoliciesFindsFiles(t *testing.T) {
	origDir := PolicyDir
	t.Cleanup(func() { PolicyDir = origDir })

	tmp := t.TempDir()
	PolicyDir = filepath.Join(tmp, "policies")
	os.MkdirAll(PolicyDir, 0o755)

	for _, name := range []string{"check1.py", "check2.yaml", "check3.yml", "readme.md"} {
		os.WriteFile(filepath.Join(PolicyDir, name), []byte("x"), 0o644)
	}

	policies, err := ListPolicies()
	if err != nil {
		t.Fatalf("ListPolicies() error = %v", err)
	}
	if len(policies) != 3 {
		t.Errorf("expected 3 policies, got %d: %v", len(policies), policies)
	}
}

func TestRemovePolicy(t *testing.T) {
	origDir := PolicyDir
	t.Cleanup(func() { PolicyDir = origDir })

	tmp := t.TempDir()
	PolicyDir = filepath.Join(tmp, "policies")
	os.MkdirAll(PolicyDir, 0o755)
	os.WriteFile(filepath.Join(PolicyDir, "old.py"), []byte("x"), 0o644)

	if err := RemovePolicy("old.py"); err != nil {
		t.Fatalf("RemovePolicy() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(PolicyDir, "old.py")); !os.IsNotExist(err) {
		t.Error("policy file should have been deleted")
	}
}

func TestRemovePolicyNotFound(t *testing.T) {
	origDir := PolicyDir
	t.Cleanup(func() { PolicyDir = origDir })
	PolicyDir = filepath.Join(t.TempDir(), "policies")
	os.MkdirAll(PolicyDir, 0o755)

	err := RemovePolicy("nonexistent.py")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}

func TestHasCustomPolicies(t *testing.T) {
	origDir := PolicyDir
	t.Cleanup(func() { PolicyDir = origDir })

	tmp := t.TempDir()
	PolicyDir = filepath.Join(tmp, "policies")

	if HasCustomPolicies() {
		t.Error("should be false when directory doesn't exist")
	}

	os.MkdirAll(PolicyDir, 0o755)
	os.WriteFile(filepath.Join(PolicyDir, "check.py"), []byte("x"), 0o644)

	if !HasCustomPolicies() {
		t.Error("should be true when policies exist")
	}
}
