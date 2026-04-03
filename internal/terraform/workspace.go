package terraform

import (
	"fmt"
	"strings"
)

func WorkspaceNew(workDir, name string) error {
	if err := EnsureBinary(); err != nil {
		return err
	}
	cmd := execCommand("terraform", "workspace", "new", name)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform workspace new %q failed:\n%s\n%w", name, string(output), err)
	}
	return nil
}

func WorkspaceSelect(workDir, name string) error {
	if err := EnsureBinary(); err != nil {
		return err
	}
	cmd := execCommand("terraform", "workspace", "select", name)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform workspace select %q failed:\n%s\n%w", name, string(output), err)
	}
	return nil
}

func WorkspaceList(workDir string) ([]string, string, error) {
	if err := EnsureBinary(); err != nil {
		return nil, "", err
	}
	cmd := execCommand("terraform", "workspace", "list")
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", fmt.Errorf("terraform workspace list failed:\n%s\n%w", string(output), err)
	}

	var workspaces []string
	var current string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "* ") {
			name := strings.TrimPrefix(line, "* ")
			current = name
			workspaces = append(workspaces, name)
		} else {
			workspaces = append(workspaces, line)
		}
	}
	return workspaces, current, nil
}
