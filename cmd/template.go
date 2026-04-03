package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Inspect available templates",
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		templates, err := listTemplateNames("templates")
		if err != nil {
			return err
		}
		out := cmd.OutOrStdout()
		for _, name := range templates {
			fmt.Fprintln(out, name)
		}
		return nil
	},
}

var templateShowCmd = &cobra.Command{
	Use:   "show [template]",
	Short: "Show a template's metadata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := findTemplateMeta("templates", args[0])
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
}

func listTemplateNames(root string) ([]string, error) {
	var names []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Name() != "meta.yaml" {
			return nil
		}
		names = append(names, filepath.ToSlash(filepath.Dir(path)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func findTemplateMeta(root, query string) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("template name is required")
	}
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() != "meta.yaml" {
			return nil
		}
		if strings.Contains(strings.ToLower(filepath.ToSlash(path)), strings.ToLower(query)) {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("template %q not found", query)
	}
	return matches[0], nil
}
