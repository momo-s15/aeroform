package terraform

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	templatedata "github.com/momo-s15/aeroform/templates"
)

// RenderTemplate copies all files from templateDir into outputDir,
// replacing Terraform variable references (var.KEY) with concrete values
// from the vars map. Files that are not .tf are copied verbatim.
//
// For relative paths under templates/ (e.g. templates/simple/aws/static-site),
// the embedded copy shipped in the binary is used first so a stray templates/
// tree in the current working directory cannot override it. To force disk
// (e.g. when developing templates), set AEROFORM_DISK_TEMPLATES=1.
//
// Absolute templateDir always reads from disk (used by tests and tooling).
func RenderTemplate(templateDir string, vars map[string]string, outputDir string) error {
	rootFS, err := resolveTemplateFS(templateDir)
	if err != nil {
		return fmt.Errorf("read template dir %s: %w", templateDir, err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir %s: %w", outputDir, err)
	}

	entries, err := fs.ReadDir(rootFS, ".")
	if err != nil {
		return fmt.Errorf("read template dir %s: %w", templateDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if err := renderFileFromFS(rootFS, entry, vars, outputDir); err != nil {
			return err
		}
	}
	return nil
}

func resolveTemplateFS(templateDir string) (fs.FS, error) {
	clean := filepath.Clean(templateDir)

	if filepath.IsAbs(clean) {
		if fi, err := os.Stat(clean); err == nil && fi.IsDir() {
			return os.DirFS(clean), nil
		}
		return nil, fmt.Errorf("template dir not found: %s", clean)
	}

	rel := embeddedRel(filepath.ToSlash(clean))
	if rel == "" {
		return nil, fmt.Errorf("empty template path after normalizing %q", templateDir)
	}

	// Default: embedded only. A ./templates tree in CWD (e.g. Downloads) must not
	// override the binary after a failed embed read — that produced stale regexreplace HCL.
	if !envUseDiskTemplates() {
		sub, err := openEmbedded(rel)
		if err != nil {
			return nil, fmt.Errorf("embedded template %q: %w (reinstall: go install github.com/momo-s15/aeroform@main)", rel, err)
		}
		return sub, nil
	}

	if fi, err := os.Stat(clean); err == nil && fi.IsDir() {
		return os.DirFS(clean), nil
	}
	sub, err := openEmbedded(rel)
	if err != nil {
		return nil, fmt.Errorf("template not on disk %s and embedded %q: %w", clean, rel, err)
	}
	return sub, nil
}

func envUseDiskTemplates() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("AEROFORM_DISK_TEMPLATES")))
	return v == "1" || v == "true" || v == "yes"
}

func embeddedRel(slashPath string) string {
	s := strings.TrimPrefix(slashPath, "./")
	if strings.HasPrefix(s, "templates/") {
		return strings.TrimPrefix(s, "templates/")
	}
	return s
}

func openEmbedded(rel string) (fs.FS, error) {
	sub, err := fs.Sub(templatedata.Files, rel)
	if err != nil {
		return nil, err
	}
	entries, err := fs.ReadDir(sub, ".")
	if err != nil || len(entries) == 0 {
		return nil, fmt.Errorf("empty or unreadable embedded dir %q", rel)
	}
	return sub, nil
}

func renderFileFromFS(srcFS fs.FS, entry fs.DirEntry, vars map[string]string, outputDir string) error {
	data, err := fs.ReadFile(srcFS, entry.Name())
	if err != nil {
		return fmt.Errorf("read %s: %w", entry.Name(), err)
	}

	content := string(data)
	if filepath.Ext(entry.Name()) == ".tf" {
		content = substituteVars(content, vars)
	}

	dstPath := filepath.Join(outputDir, entry.Name())
	if err := os.WriteFile(dstPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", dstPath, err)
	}
	return nil
}

// substituteVars replaces var.KEY references in Terraform HCL with the
// concrete value from the vars map. It handles both var.key (bare) and
// "${var.key}" (interpolation) forms.
func substituteVars(content string, vars map[string]string) string {
	for key, value := range vars {
		// "${var.key}" → "value"
		content = strings.ReplaceAll(content, fmt.Sprintf("${var.%s}", key), value)
		// var.key (bare reference) → "value"
		content = strings.ReplaceAll(content, fmt.Sprintf("var.%s", key), fmt.Sprintf("%q", value))
	}
	return content
}

// ProTemplateSource identifies a single template to include in a Pro Mode
// multi-template project. Name is the module name (e.g. "vpc"), Dir is the
// filesystem path to the template directory.
type ProTemplateSource struct {
	Name string
	Dir  string
}

// RenderProProject composes multiple Pro Mode templates into a single
// Terraform root module. Each template is rendered into its own subdirectory
// and a root main.tf is generated that wires them together as modules.
func RenderProProject(templates []ProTemplateSource, vars map[string]string, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir %s: %w", outputDir, err)
	}

	for _, tmpl := range templates {
		moduleDir := filepath.Join(outputDir, tmpl.Name)
		if err := RenderTemplate(tmpl.Dir, nil, moduleDir); err != nil {
			return fmt.Errorf("render module %s: %w", tmpl.Name, err)
		}
	}

	var b strings.Builder
	b.WriteString("terraform {\n  required_version = \">= 1.0\"\n}\n\n")
	for _, tmpl := range templates {
		fmt.Fprintf(&b, "module %q {\n  source = \"./%s\"\n}\n\n", tmpl.Name, tmpl.Name)
	}

	rootMain := filepath.Join(outputDir, "main.tf")
	if err := os.WriteFile(rootMain, []byte(b.String()), 0o644); err != nil {
		return fmt.Errorf("write root main.tf: %w", err)
	}

	if len(vars) > 0 {
		if err := GenerateTfvars(vars, outputDir); err != nil {
			return err
		}
	}

	return nil
}

// GenerateTfvars writes a terraform.tfvars file into outputDir with all
// vars as key = "value" pairs. This is the primary mechanism for passing
// variables to terraform plan/apply.
func GenerateTfvars(vars map[string]string, outputDir string) error {
	if len(vars) == 0 {
		return nil
	}

	keys := sortedKeys(vars)
	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "%s = %q\n", k, vars[k])
	}

	path := filepath.Join(outputDir, "terraform.tfvars")
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return fmt.Errorf("write terraform.tfvars: %w", err)
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Simple insertion sort — var maps are small.
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
