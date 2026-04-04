package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)

// GCP project IDs: 6–30 chars, start with letter, lowercase letters, digits, hyphens.
var gcpProjectIDPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{4,28}[a-z0-9]$`)

func uiSelect(label string, items []string) (string, error) {
	sel := promptui.Select{
		Label: label,
		Items: items,
		Size:  len(items),
	}
	_, result, err := sel.Run()
	return result, err
}

func uiPrompt(label, defaultVal string, validate promptui.ValidateFunc) (string, error) {
	p := promptui.Prompt{
		Label:    label,
		Default:  defaultVal,
		Validate: validate,
	}
	return p.Run()
}

func uiPromptOptional(label string) (string, error) {
	p := promptui.Prompt{
		Label:   label,
		Default: "",
	}
	return p.Run()
}

func uiConfirm(label string) (bool, error) {
	p := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Default:   "y",
	}
	result, err := p.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}
	return strings.EqualFold(result, "y") || strings.EqualFold(result, "yes") || result == "", nil
}

func validateNotEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("value cannot be empty")
	}
	return nil
}

func validateSlug(input string) error {
	v := strings.TrimSpace(input)
	if v == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if len(v) < 2 || len(v) > 63 {
		return fmt.Errorf("project name must be 2–63 characters")
	}
	if !slugPattern.MatchString(v) {
		return fmt.Errorf("use only lowercase letters, digits, and hyphens (no leading/trailing hyphen)")
	}
	return nil
}

func validateGCPProjectID(input string) error {
	v := strings.TrimSpace(input)
	if v == "" {
		return fmt.Errorf("GCP project ID cannot be empty")
	}
	if len(v) < 6 || len(v) > 30 {
		return fmt.Errorf("GCP project ID must be 6–30 characters")
	}
	if !gcpProjectIDPattern.MatchString(v) {
		return fmt.Errorf("must start with a letter; use only lowercase letters, digits, and hyphens")
	}
	return nil
}
