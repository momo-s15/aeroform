package security

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var PolicyDir = ".aeroform/policies"

func AddPolicy(src string) (string, error) {
	ext := strings.ToLower(filepath.Ext(src))
	if ext != ".py" && ext != ".yaml" && ext != ".yml" {
		return "", fmt.Errorf("unsupported policy file type %q: expected .py, .yaml, or .yml", ext)
	}

	info, err := os.Stat(src)
	if err != nil {
		return "", fmt.Errorf("cannot read policy file: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("%q is a directory, not a policy file", src)
	}

	if err := os.MkdirAll(PolicyDir, 0o755); err != nil {
		return "", fmt.Errorf("create policy directory: %w", err)
	}

	dst := filepath.Join(PolicyDir, filepath.Base(src))
	if err := copyFile(src, dst); err != nil {
		return "", fmt.Errorf("copy policy: %w", err)
	}
	return dst, nil
}

func ListPolicies() ([]string, error) {
	entries, err := os.ReadDir(PolicyDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var policies []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".py" || ext == ".yaml" || ext == ".yml" {
			policies = append(policies, e.Name())
		}
	}
	return policies, nil
}

func RemovePolicy(name string) error {
	path := filepath.Join(PolicyDir, filepath.Base(name))
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("policy %q not found", name)
		}
		return fmt.Errorf("remove policy: %w", err)
	}
	return nil
}

func HasCustomPolicies() bool {
	policies, err := ListPolicies()
	if err != nil {
		return false
	}
	return len(policies) > 0
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	if err := out.Sync(); err != nil {
		return err
	}
	return nil
}
