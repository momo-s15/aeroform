package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
)

// One reader for Windows plain stdin prompts (promptui mishandles some keys in PowerShell).
var winStdin = bufio.NewReader(os.Stdin)

func useWindowsStdinPrompts() bool {
	return runtime.GOOS == "windows"
}

func windowsReadLine(label, defaultVal string, validate promptui.ValidateFunc) (string, error) {
	for {
		if strings.TrimSpace(defaultVal) != "" {
			fmt.Fprintf(os.Stdout, "%s [%s]: ", label, defaultVal)
		} else {
			fmt.Fprintf(os.Stdout, "%s: ", label)
		}
		line, err := winStdin.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line == "" {
			line = strings.TrimSpace(defaultVal)
		}
		if validate != nil {
			if err := validate(line); err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
				continue
			}
		}
		return line, nil
	}
}

func windowsReadLineOptional(label string) (string, error) {
	fmt.Fprintf(os.Stdout, "%s (Enter to skip): ", label)
	line, err := winStdin.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.TrimSuffix(line, "\r")), nil
}

func windowsConfirm(label string) (bool, error) {
	for {
		fmt.Fprintf(os.Stdout, "%s [y/N]: ", label)
		line, err := winStdin.ReadString('\n')
		if err != nil {
			return false, err
		}
		line = strings.TrimSpace(strings.ToLower(strings.TrimSuffix(line, "\r")))
		if line == "" || line == "y" || line == "yes" {
			return true, nil
		}
		if line == "n" || line == "no" {
			return false, nil
		}
		fmt.Fprintln(os.Stdout, "Please enter y, n, or Enter for yes.")
	}
}

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
	if useWindowsStdinPrompts() {
		return windowsReadLine(label, defaultVal, validate)
	}
	p := promptui.Prompt{
		Label:    label,
		Default:  defaultVal,
		Validate: validate,
	}
	return p.Run()
}

func uiPromptOptional(label string) (string, error) {
	if useWindowsStdinPrompts() {
		return windowsReadLineOptional(label)
	}
	p := promptui.Prompt{
		Label:       label,
		Default:     "",
		HideEntered: true, // avoids a second echoed line on Windows after Enter (esp. when skipping blank)
	}
	return p.Run()
}

func uiConfirm(label string) (bool, error) {
	if useWindowsStdinPrompts() {
		return windowsConfirm(label)
	}
	p := promptui.Prompt{
		Label:       label,
		IsConfirm:   true,
		Default:     "y",
		HideEntered: true, // cleaner line after Y/n on Windows terminals
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

// Azure region names: lowercase, no spaces (e.g. eastus, canadacentral).
var azureRegionPattern = regexp.MustCompile(`^[a-z]{2,}[-a-z0-9]*$`)

func validateAzureRegion(input string) error {
	v := strings.TrimSpace(strings.ToLower(input))
	if v == "" {
		return fmt.Errorf("Azure region cannot be empty")
	}
	if len(v) < 5 || len(v) > 40 {
		return fmt.Errorf("Azure region looks invalid (length)")
	}
	if !azureRegionPattern.MatchString(v) {
		return fmt.Errorf("use a lowercase Azure region name (e.g. canadacentral, eastus)")
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
