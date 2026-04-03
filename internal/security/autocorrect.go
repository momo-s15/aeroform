package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Correction describes a single auto-fix applied to rendered Terraform files.
type Correction struct {
	File        string
	Rule        string
	Description string
}

// AutoCorrect scans all .tf files in tfDir and applies known safe fixes.
// Returns the list of corrections applied. Does not fail on unfixable issues —
// those are left for the security gate / Checkov to catch.
func AutoCorrect(tfDir string) ([]Correction, error) {
	entries, err := os.ReadDir(tfDir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", tfDir, err)
	}

	var corrections []Correction

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".tf" {
			continue
		}

		path := filepath.Join(tfDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}

		original := string(data)
		content := original

		if fixed, ok := fixPubliclyAccessible(content); ok {
			content = fixed
			corrections = append(corrections, Correction{
				File:        entry.Name(),
				Rule:        "no-public-db",
				Description: "Set publicly_accessible = false (databases must not be public)",
			})
		}

		if fixed, ok := fixHTTPtoHTTPS(content); ok {
			content = fixed
			corrections = append(corrections, Correction{
				File:        entry.Name(),
				Rule:        "enforce-https",
				Description: "Changed protocol from HTTP to HTTPS",
			})
		}

		if fixed, ok := fixUnencryptedStorage(content); ok {
			content = fixed
			corrections = append(corrections, Correction{
				File:        entry.Name(),
				Rule:        "enforce-encryption",
				Description: "Ensured storage_encrypted = true on database instances",
			})
		}

		if content != original {
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return nil, fmt.Errorf("write %s: %w", path, err)
			}
		}
	}

	return corrections, nil
}

func fixPubliclyAccessible(content string) (string, bool) {
	target := "publicly_accessible = true"
	if !strings.Contains(content, target) {
		return content, false
	}
	return strings.ReplaceAll(content, target, "publicly_accessible = false"), true
}

func fixHTTPtoHTTPS(content string) (string, bool) {
	target := `protocol = "HTTP"`
	replacement := `protocol = "HTTPS"`
	if !strings.Contains(content, target) {
		return content, false
	}
	return strings.ReplaceAll(content, target, replacement), true
}

func fixUnencryptedStorage(content string) (string, bool) {
	target := "storage_encrypted = false"
	if !strings.Contains(content, target) {
		return content, false
	}
	return strings.ReplaceAll(content, target, "storage_encrypted = true"), true
}
