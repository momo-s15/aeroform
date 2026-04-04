package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Renderer tests
// ---------------------------------------------------------------------------

func TestRenderTemplateCreatesOutputFiles(t *testing.T) {
	templateDir := t.TempDir()
	outputDir := t.TempDir()

	mainTF := `resource "aws_s3_bucket" "site" {
  bucket = var.project_name
  tags = {
    Name = "${var.project_name}"
  }
}
`
	if err := os.WriteFile(filepath.Join(templateDir, "main.tf"), []byte(mainTF), 0o644); err != nil {
		t.Fatal(err)
	}

	readme := "This is a readme."
	if err := os.WriteFile(filepath.Join(templateDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"project_name": "my-site",
	}

	if err := RenderTemplate(templateDir, vars, outputDir); err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	rendered, err := os.ReadFile(filepath.Join(outputDir, "main.tf"))
	if err != nil {
		t.Fatalf("rendered main.tf not found: %v", err)
	}
	content := string(rendered)
	if strings.Contains(content, "var.project_name") {
		t.Fatalf("var.project_name was not substituted:\n%s", content)
	}
	if !strings.Contains(content, `"my-site"`) {
		t.Fatalf("expected substituted value \"my-site\" in output:\n%s", content)
	}

	readmeOut, err := os.ReadFile(filepath.Join(outputDir, "README.md"))
	if err != nil {
		t.Fatalf("README.md not copied: %v", err)
	}
	if string(readmeOut) != readme {
		t.Fatalf("README.md content changed: got %q", string(readmeOut))
	}
}

func TestRenderTemplateFailsOnMissingDir(t *testing.T) {
	err := RenderTemplate("/nonexistent/path", nil, t.TempDir())
	if err == nil {
		t.Fatal("expected error for missing template dir")
	}
}

func TestRenderTemplateUsesEmbedWhenNotOnDisk(t *testing.T) {
	t.Setenv("AEROFORM_DISK_TEMPLATES", "")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	out := t.TempDir()
	err = RenderTemplate("templates/simple/aws/static-site", map[string]string{
		"project_name": "embed-smoke",
		"region":       "us-east-1",
	}, out)
	if err != nil {
		t.Fatalf("RenderTemplate from embed: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(out, "main.tf"))
	if err != nil || len(b) < 50 {
		t.Fatalf("expected embedded main.tf: %v len=%d", err, len(b))
	}
}

func TestRenderTemplatePrefersEmbeddedOverStaleDisk(t *testing.T) {
	t.Setenv("AEROFORM_DISK_TEMPLATES", "")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	root := t.TempDir()
	staleDir := filepath.Join(root, "templates", "simple", "aws", "static-site")
	if err := os.MkdirAll(staleDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staleDir, "main.tf"), []byte("# STALE_DISK_OVERRIDE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

	out := t.TempDir()
	if err := RenderTemplate("templates/simple/aws/static-site", nil, out); err != nil {
		t.Fatalf("RenderTemplate: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(out, "main.tf"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if strings.Contains(s, "STALE_DISK_OVERRIDE") {
		t.Fatal("embedded template should win over a templates/ tree in CWD")
	}
	if !strings.Contains(s, `resource "aws_s3_bucket"`) {
		t.Fatalf("expected shipped AWS static-site template, got:\n%s", s)
	}
}

func TestRenderTemplateDiskTemplatesEnvUsesCWD(t *testing.T) {
	t.Setenv("AEROFORM_DISK_TEMPLATES", "1")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	root := t.TempDir()
	staleDir := filepath.Join(root, "templates", "simple", "aws", "static-site")
	if err := os.MkdirAll(staleDir, 0o755); err != nil {
		t.Fatal(err)
	}
	want := "# DISK_OVERRIDE_MARKER\n"
	if err := os.WriteFile(filepath.Join(staleDir, "main.tf"), []byte(want), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

	out := t.TempDir()
	if err := RenderTemplate("templates/simple/aws/static-site", nil, out); err != nil {
		t.Fatalf("RenderTemplate: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(out, "main.tf"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "DISK_OVERRIDE_MARKER") {
		t.Fatalf("AEROFORM_DISK_TEMPLATES=1 should use CWD disk copy; got:\n%s", string(b))
	}
}

func TestGenerateTfvars(t *testing.T) {
	dir := t.TempDir()
	vars := map[string]string{
		"project_name": "demo",
		"region":       "us-east-1",
	}
	if err := GenerateTfvars(vars, dir); err != nil {
		t.Fatalf("GenerateTfvars() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "terraform.tfvars"))
	if err != nil {
		t.Fatalf("terraform.tfvars not created: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, `project_name = "demo"`) {
		t.Fatalf("missing project_name in tfvars:\n%s", content)
	}
	if !strings.Contains(content, `region = "us-east-1"`) {
		t.Fatalf("missing region in tfvars:\n%s", content)
	}
}

func TestGenerateTfvarsEmptyVars(t *testing.T) {
	dir := t.TempDir()
	if err := GenerateTfvars(map[string]string{}, dir); err != nil {
		t.Fatalf("GenerateTfvars() error = %v", err)
	}
	_, err := os.Stat(filepath.Join(dir, "terraform.tfvars"))
	if err == nil {
		t.Fatal("expected no terraform.tfvars for empty vars")
	}
}

// ---------------------------------------------------------------------------
// Pro project rendering tests
// ---------------------------------------------------------------------------

func TestRenderProProjectCreatesModulesAndRoot(t *testing.T) {
	tmpDir := t.TempDir()
	vpcDir := filepath.Join(tmpDir, "tpl-vpc")
	eksDir := filepath.Join(tmpDir, "tpl-eks")
	os.MkdirAll(vpcDir, 0o755)
	os.MkdirAll(eksDir, 0o755)

	os.WriteFile(filepath.Join(vpcDir, "main.tf"), []byte(`resource "aws_vpc" "this" {}`), 0o644)
	os.WriteFile(filepath.Join(vpcDir, "variables.tf"), []byte(`variable "cidr_block" {}`), 0o644)
	os.WriteFile(filepath.Join(eksDir, "main.tf"), []byte(`resource "aws_eks_cluster" "this" {}`), 0o644)

	outDir := filepath.Join(tmpDir, "output")
	sources := []ProTemplateSource{
		{Name: "vpc", Dir: vpcDir},
		{Name: "eks", Dir: eksDir},
	}
	vars := map[string]string{"region": "us-east-1"}

	if err := RenderProProject(sources, vars, outDir); err != nil {
		t.Fatalf("RenderProProject() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "vpc", "main.tf")); err != nil {
		t.Fatal("vpc/main.tf not created")
	}
	if _, err := os.Stat(filepath.Join(outDir, "vpc", "variables.tf")); err != nil {
		t.Fatal("vpc/variables.tf not created")
	}
	if _, err := os.Stat(filepath.Join(outDir, "eks", "main.tf")); err != nil {
		t.Fatal("eks/main.tf not created")
	}

	rootMain, err := os.ReadFile(filepath.Join(outDir, "main.tf"))
	if err != nil {
		t.Fatalf("root main.tf not created: %v", err)
	}
	content := string(rootMain)
	if !strings.Contains(content, `module "vpc"`) {
		t.Fatal("root main.tf missing vpc module block")
	}
	if !strings.Contains(content, `module "eks"`) {
		t.Fatal("root main.tf missing eks module block")
	}
	if !strings.Contains(content, `source = "./vpc"`) {
		t.Fatal("root main.tf missing vpc source path")
	}

	tfvars, err := os.ReadFile(filepath.Join(outDir, "terraform.tfvars"))
	if err != nil {
		t.Fatalf("terraform.tfvars not created: %v", err)
	}
	if !strings.Contains(string(tfvars), `region = "us-east-1"`) {
		t.Fatal("tfvars missing region")
	}
}

func TestRenderProProjectEmptyTemplates(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "output")
	if err := RenderProProject(nil, nil, outDir); err != nil {
		t.Fatalf("RenderProProject(nil) error = %v", err)
	}

	rootMain, err := os.ReadFile(filepath.Join(outDir, "main.tf"))
	if err != nil {
		t.Fatalf("root main.tf not created: %v", err)
	}
	if !strings.Contains(string(rootMain), "required_version") {
		t.Fatal("root main.tf missing terraform block")
	}
}

// ---------------------------------------------------------------------------
// Runner tests
// ---------------------------------------------------------------------------

func TestEnsureBinaryMissing(t *testing.T) {
	original := execLookPath
	t.Cleanup(func() { execLookPath = original })

	execLookPath = func(file string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	err := EnsureBinary()
	if err == nil {
		t.Fatal("expected error when terraform is not in PATH")
	}
	if !strings.Contains(err.Error(), "not installed") {
		t.Fatalf("error message should mention installation, got: %v", err)
	}
}

func TestEnsureBinaryFound(t *testing.T) {
	original := execLookPath
	t.Cleanup(func() { execLookPath = original })

	execLookPath = func(file string) (string, error) {
		return "/usr/local/bin/terraform", nil
	}

	if err := EnsureBinary(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInitFailsWhenTerraformMissing(t *testing.T) {
	original := execLookPath
	t.Cleanup(func() { execLookPath = original })

	execLookPath = func(file string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	err := Init(t.TempDir())
	if err == nil {
		t.Fatal("expected error when terraform is missing")
	}
}

func TestPlanFailsWhenTerraformMissing(t *testing.T) {
	original := execLookPath
	t.Cleanup(func() { execLookPath = original })

	execLookPath = func(file string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	_, err := Plan(t.TempDir())
	if err == nil {
		t.Fatal("expected error when terraform is missing")
	}
}

func TestParsePlanCounts(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		add     int
		change  int
		destroy int
	}{
		{
			name:    "typical plan output",
			output:  "Plan: 5 to add, 2 to change, 1 to destroy.",
			add:     5,
			change:  2,
			destroy: 1,
		},
		{
			name:    "no changes",
			output:  "No changes. Your infrastructure matches the configuration.",
			add:     0,
			change:  0,
			destroy: 0,
		},
		{
			name:    "add only",
			output:  "Plan: 12 to add, 0 to change, 0 to destroy.",
			add:     12,
			change:  0,
			destroy: 0,
		},
		{
			name:    "embedded in larger output",
			output:  "Refreshing...\n\nPlan: 3 to add, 1 to change, 0 to destroy.\n\nDo you want to apply?",
			add:     3,
			change:  1,
			destroy: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := &PlanResult{}
			parsePlanCounts(tc.output, result)
			if result.AddCount != tc.add {
				t.Errorf("AddCount = %d, want %d", result.AddCount, tc.add)
			}
			if result.ChangeCount != tc.change {
				t.Errorf("ChangeCount = %d, want %d", result.ChangeCount, tc.change)
			}
			if result.DestroyCount != tc.destroy {
				t.Errorf("DestroyCount = %d, want %d", result.DestroyCount, tc.destroy)
			}
		})
	}
}

func TestDestroyFailsWhenTerraformMissing(t *testing.T) {
	original := execLookPath
	t.Cleanup(func() { execLookPath = original })

	execLookPath = func(file string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	err := Destroy(t.TempDir(), true)
	if err == nil {
		t.Fatal("expected error when terraform is missing")
	}
}

func TestApplyFailsWhenTerraformMissing(t *testing.T) {
	original := execLookPath
	t.Cleanup(func() { execLookPath = original })

	execLookPath = func(file string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	err := Apply(t.TempDir(), true)
	if err == nil {
		t.Fatal("expected error when terraform is missing")
	}
}

func TestOutputFailsWhenTerraformMissing(t *testing.T) {
	original := execLookPath
	t.Cleanup(func() { execLookPath = original })

	execLookPath = func(file string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	_, err := Output(t.TempDir())
	if err == nil {
		t.Fatal("expected error when terraform is missing")
	}
}

func TestApplyPassesAutoApproveFlag(t *testing.T) {
	originalLookPath := execLookPath
	originalCommand := execCommand
	t.Cleanup(func() {
		execLookPath = originalLookPath
		execCommand = originalCommand
	})

	execLookPath = func(file string) (string, error) {
		return file, nil
	}

	var capturedArgs []string
	execCommand = func(name string, args ...string) *exec.Cmd {
		capturedArgs = append([]string{name}, args...)
		return exec.Command("echo", "ok")
	}

	dir := t.TempDir()
	_ = Apply(dir, true)

	found := false
	for _, arg := range capturedArgs {
		if arg == "-auto-approve" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected -auto-approve in args: %v", capturedArgs)
	}
}

func TestDestroyPassesAutoApproveFlag(t *testing.T) {
	originalLookPath := execLookPath
	originalCommand := execCommand
	t.Cleanup(func() {
		execLookPath = originalLookPath
		execCommand = originalCommand
	})

	execLookPath = func(file string) (string, error) {
		return file, nil
	}

	var capturedArgs []string
	execCommand = func(name string, args ...string) *exec.Cmd {
		capturedArgs = append([]string{name}, args...)
		return exec.Command("echo", "ok")
	}

	dir := t.TempDir()
	_ = Destroy(dir, true)

	found := false
	for _, arg := range capturedArgs {
		if arg == "-auto-approve" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected -auto-approve in args: %v", capturedArgs)
	}
}

func TestParseDriftChangesUpdate(t *testing.T) {
	output := `  # aws_instance.web will be updated in-place
  ~ resource "aws_instance" "web" {
        ami = "ami-123" -> "ami-456"
    }

  # aws_s3_bucket.data will be updated in-place
  ~ resource "aws_s3_bucket" "data" {
        tags = {}
    }
`
	changes := parseDriftChanges(output)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	if changes[0].Address != "aws_instance.web" {
		t.Fatalf("expected aws_instance.web, got %q", changes[0].Address)
	}
	if changes[0].ChangeType != "update" {
		t.Fatalf("expected 'update', got %q", changes[0].ChangeType)
	}
}

func TestParseDriftChangesCreate(t *testing.T) {
	output := `  # aws_security_group.new will be created
  + resource "aws_security_group" "new" {
    }
`
	changes := parseDriftChanges(output)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].ChangeType != "create" {
		t.Fatalf("expected 'create', got %q", changes[0].ChangeType)
	}
}

func TestParseDriftChangesDestroy(t *testing.T) {
	output := `  # aws_instance.old will be destroyed
  - resource "aws_instance" "old" {
    }
`
	changes := parseDriftChanges(output)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].ChangeType != "destroy" {
		t.Fatalf("expected 'destroy', got %q", changes[0].ChangeType)
	}
}

func TestParseDriftChangesReplace(t *testing.T) {
	output := `  # aws_instance.app must be replaced
  -/+ resource "aws_instance" "app" {
    }
`
	changes := parseDriftChanges(output)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].ChangeType != "replace" {
		t.Fatalf("expected 'replace', got %q", changes[0].ChangeType)
	}
}

func TestParseDriftChangesNoDuplicates(t *testing.T) {
	output := `  # aws_instance.web will be updated in-place
  # aws_instance.web will be updated in-place
`
	changes := parseDriftChanges(output)
	if len(changes) != 1 {
		t.Fatalf("expected 1 unique change, got %d", len(changes))
	}
}

func TestParseDriftChangesEmpty(t *testing.T) {
	changes := parseDriftChanges("No changes. Your infrastructure matches the configuration.")
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestNormalizeChangeType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"created", "create"},
		{"destroyed", "destroy"},
		{"updated", "update"},
		{"updated-in-place", "update"},
		{"replaced", "replace"},
		{"changed", "update"},
		{"something-else", "something-else"},
	}
	for _, tc := range tests {
		got := normalizeChangeType(tc.input)
		if got != tc.want {
			t.Errorf("normalizeChangeType(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
